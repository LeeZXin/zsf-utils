package lease

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"
	"xorm.io/xorm"
)

type Leaser interface {
	TryGrant() (Releaser, Renewer, bool, error)
}

type Releaser interface {
	Release() error
}

type Renewer interface {
	Renew(context.Context) (bool, error)
}

type DbModel struct {
	Id       int64     `json:"id" xorm:"pk autoincr"`
	LeaseKey string    `json:"leaseKey"`
	Owner    string    `json:"owner"`
	Renewed  time.Time `json:"renewed"`
	Created  time.Time `json:"created" xorm:"created"`
}

type dbLease struct {
	Key             string
	Owner           string
	TableName       string
	Engine          *xorm.Engine
	ExpiredDuration time.Duration
}

type dbReleaser struct {
	Lease       *DbModel
	TableName   string
	Owner       string
	Engine      *xorm.Engine
	releaseOnce sync.Once
}

func (r *dbReleaser) Release() (err error) {
	r.releaseOnce.Do(func() {
		session := r.Engine.NewSession()
		defer session.Close()
		_, err = session.
			Where("id = ?", r.Lease.Id).
			And("owner = ?", r.Owner).
			Table(r.TableName).
			Delete(new(DbModel))
	})
	return
}

type dbRenewer struct {
	Lease     *DbModel
	TableName string
	Owner     string
	Engine    *xorm.Engine
}

func (r *dbRenewer) Renew(context.Context) (bool, error) {
	session := r.Engine.NewSession()
	defer session.Close()
	rows, err := session.
		Where("id = ?", r.Lease.Id).
		And("owner = ?", r.Owner).
		Table(r.TableName).
		Cols("renewed").
		Update(&DbModel{
			Renewed: time.Now(),
		})
	return rows == 1, err
}

func NewDbLease(key, owner, tableName string, engine *xorm.Engine, expiredDuration time.Duration) (Leaser, error) {
	if key == "" {
		return nil, errors.New("empty key")
	}
	if owner == "" {
		return nil, errors.New("empty owner")
	}
	if tableName == "" {
		return nil, errors.New("empty table name")
	}
	if engine == nil {
		return nil, errors.New("nil Engine")
	}
	if expiredDuration <= 0 {
		return nil, errors.New("wrong duration")
	}
	return &dbLease{
		Key:             key,
		Owner:           owner,
		TableName:       tableName,
		Engine:          engine,
		ExpiredDuration: expiredDuration,
	}, nil
}

func (l *dbLease) TryGrant() (Releaser, Renewer, bool, error) {
	session := l.Engine.NewSession()
	defer session.Close()
	md, success, err := l.tryGrant(session)
	if err != nil {
		return nil, nil, false, err
	}
	if success {
		return &dbReleaser{
				Lease:     &md,
				TableName: l.TableName,
				Engine:    l.Engine,
				Owner:     l.Owner,
			}, &dbRenewer{
				Lease:     &md,
				TableName: l.TableName,
				Engine:    l.Engine,
				Owner:     l.Owner,
			}, true, nil
	}
	return nil, nil, false, nil
}

// tryGrant 使用乐观锁
func (l *dbLease) tryGrant(session *xorm.Session) (DbModel, bool, error) {
	var md DbModel
	// select for update
	b, err := session.
		Where("lease_key = ?", l.Key).
		Table(l.TableName).
		Get(&md)
	if err != nil {
		return md, false, err
	}
	now := time.Now()
	if !b {
		md = DbModel{
			LeaseKey: l.Key,
			Owner:    l.Owner,
			Renewed:  now,
		}
		// 不存在则插入
		_, err = session.Table(l.TableName).Insert(&md)
		if err != nil {
			// 唯一键冲突
			if strings.Contains(err.Error(), "Error 1062") {
				return md, false, nil
			}
			return md, false, err
		}
		// 加锁成功
		return md, true, nil
	}
	// 锁过期
	if md.Renewed.Before(now.Add(-l.ExpiredDuration)) {
		oldOwner := md.Owner
		md.Owner = l.Owner
		md.Renewed = now
		rows, err := session.Where("id = ?", md.Id).
			And("owner = ?", oldOwner).
			Cols("owner", "renewed").
			Table(l.TableName).
			Update(&md)
		if err != nil {
			return md, false, err
		}
		return md, rows == 1, nil
	}
	// 没过期则判断有效的锁的owner是不是自己
	return md, md.Owner == l.Owner, nil
}
