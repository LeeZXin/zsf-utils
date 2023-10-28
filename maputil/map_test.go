package maputil

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewConcurrentMap(t *testing.T) {
	m := NewConcurrentMap[string, string](nil)
	fmt.Println(m.LoadOrStore("k", "v"))
	fmt.Println(m.Load("k"))
	fmt.Println(m.LoadOrStoreWithLoader("k", func() (string, error) {
		return "v2", nil
	}))
	fmt.Println(m.Load("k"))
	fmt.Println(m.LoadOrStoreWithLoader("k2", func() (string, error) {
		return "v2", errors.New("fff")
	}))
	fmt.Println(m.Load("k2"))
	fmt.Println(m.AllKeys())
	m.RemoveKey("k")
	fmt.Println(m.AllKeys())
}
