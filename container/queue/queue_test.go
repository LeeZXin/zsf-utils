package queue

import (
	"fmt"
	"testing"
	"time"
)

type Item struct {
	v string
	d time.Duration
}

func (i *Item) GetObject() string {
	return i.v
}

func (i *Item) GetDelayedDuration() time.Duration {
	return i.d
}

func TestNewDelayedQueue(t *testing.T) {
	queue := NewDelayedQueue[string]()
	go func() {
		for {
			fmt.Println("take", time.Now(), 1)
			d := queue.Take()
			o := d.GetObject()
			fmt.Println(o, time.Now(), 1)
		}
	}()
	go func() {
		for {
			fmt.Println("take", time.Now(), 2)
			d := queue.Take()
			o := d.GetObject()
			fmt.Println(o, time.Now(), 2)
			time.Sleep(time.Second)
		}
	}()
	queue.Push(&Item{
		v: "a",
		d: 3 * time.Second,
	})
	queue.Push(&Item{
		v: "b",
		d: 5 * time.Second,
	})
	queue.Push(&Item{
		v: "c",
		d: 1 * time.Second,
	})
	queue.Push(&Item{
		v: "d",
		d: 10 * time.Second,
	})
	queue.Push(&Item{
		v: "f",
		d: 100 * time.Millisecond,
	})
	time.Sleep(time.Hour)
}
