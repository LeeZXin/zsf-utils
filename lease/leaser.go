package lease

import (
	"context"
	"errors"
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
	Id      int64     `json:"id" xorm:"pk autoincr"`
	Key     string    `json:"key"`
	Owner   string    `json:"owner"`
	Renewed time.Time `json:"renewed"`
	Created time.Time `json:"created" xorm:"created"`
}

type dbLease struct {
	Key             string
	Owner           string
	TableName       string
	Engine          *xorm.Engine
	ExpiredDuration time.Duration
}

type dbReleaser struct {
	Lease           *DbModel
	TableName       string
	Owner           string
	Engine          *xorm.Engine
	ExpiredDuration time.Duration
}

func (r *dbReleaser) Release() error {
	session := r.Engine.NewSession()
	defer session.Close()
	_, err := session.
		Where("id = ?", r.Lease.Id).
		And("owner = ?", r.Owner).
		And("renewed > ?", time.Now().Add(-r.ExpiredDuration).Format(time.DateTime)).
		Table(r.TableName).
		Delete(new(DbModel))
	return err
}

type dbRenewer struct {
	Lease           *DbModel
	TableName       string
	Owner           string
	Engine          *xorm.Engine
	ExpiredDuration time.Duration
}

func (r *dbRenewer) Renew(context.Context) (bool, error) {
	session := r.Engine.NewSession()
	defer session.Close()
	rows, err := session.
		Where("id = ?", r.Lease.Id).
		And("owner = ?", r.Owner).
		And("renewed > ?", time.Now().Add(-r.ExpiredDuration).Format(time.DateTime)).
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
	err := session.Begin()
	if err != nil {
		return nil, nil, false, err
	}
	md, success, err := l.tryGrantInTx(session)
	if err != nil {
		return nil, nil, false, err
	}
	err = session.Commit()
	if err != nil {
		return nil, nil, false, err
	}
	if success {
		return &dbReleaser{
				Lease:           &md,
				TableName:       l.TableName,
				Engine:          l.Engine,
				Owner:           l.Owner,
				ExpiredDuration: l.ExpiredDuration,
			}, &dbRenewer{
				Lease:           &md,
				TableName:       l.TableName,
				Engine:          l.Engine,
				Owner:           l.Owner,
				ExpiredDuration: l.ExpiredDuration,
			}, true, nil
	}
	return nil, nil, false, nil
}

func (l *dbLease) tryGrantInTx(session *xorm.Session) (DbModel, bool, error) {
	var md DbModel
	// select for update
	b, err := session.
		Where("key = ?", l.Key).
		Table(l.TableName).
		ForUpdate().
		Get(&md)
	if err != nil {
		return md, false, err
	}
	now := time.Now()
	if !b {
		md = DbModel{
			Key:     l.Key,
			Owner:   l.Owner,
			Renewed: now,
		}
		// 不存在则插入
		_, err = session.Table(l.TableName).Insert(&md)
		if err != nil {
			return md, false, err
		}
		// 加锁成功
		return md, true, nil
	}
	// 锁过期
	if md.Renewed.Before(now.Add(-l.ExpiredDuration)) {
		md.Owner = l.Owner
		md.Renewed = now
		_, err = session.Where("id = ?", md.Id).
			Cols("owner", "renewed").
			Table(l.TableName).
			Update(&md)
		if err != nil {
			return md, false, nil
		}
		// 加锁成功
		return md, true, nil
	}
	// 没过期则判断有效的锁的owner是不是自己
	return md, md.Owner == l.Owner, nil
}
