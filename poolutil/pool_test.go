package poolutil

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type StringObject struct {
	d string
}

func (*StringObject) IsActive() bool {
	return true
}
func (s *StringObject) GetObject() string {
	return s.d
}

func (*StringObject) Close() {}

type StringObjectFactory struct {
}

func (*StringObjectFactory) CreateObject() (Object[string], error) {
	return &StringObject{
		d: strconv.Itoa(time.Now().Nanosecond()),
	}, nil
}

func TestNewGenericPool(t *testing.T) {
	pool, err := NewGenericPool(PoolConfig[string]{
		MinIdle:   0,
		MaxIdle:   1,
		MaxActive: 5,
		Factory:   &StringObjectFactory{},
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < 20; i++ {
		fi := i
		go func() {
			until, err2 := pool.BorrowObjectUntil(-1)
			if err2 == nil {
				pool.ReturnObject(until)
			}
			fmt.Println(time.Now(), until, err2, fi)
		}()
	}
	time.Sleep(500 * time.Millisecond)
	pool.Close()
	time.Sleep(time.Minute)
	fmt.Println(pool.GetIdleNum())
}
