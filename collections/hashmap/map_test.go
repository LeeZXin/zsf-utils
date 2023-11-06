package hashmap

import (
	"fmt"
	"testing"
)

type TM[K comparable, V any] struct {
	*LinkedHashMap[K, V]
}

func TestNewLinkedHashMap(t *testing.T) {
	m := &TM[string, int]{
		LinkedHashMap: NewLinkedHashMapWithRemoveEldestKeyFn[string, int](true, func(m map[string]int) bool {
			fmt.Println(m)
			return len(m) > 2
		}),
	}
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("c", 3)
	m.Put("d", 4)
	m.Put("e", 5)
	fmt.Println(m.AllKeys())
	m.Get("a")
	fmt.Println(m.AllKeys())
	m.Get("a")
	fmt.Println(m.AllKeys())
	m.Get("b")
	fmt.Println(m.AllKeys())
	m.Get("d")
	fmt.Println(m.AllKeys())
	m.Get("f")
	fmt.Println(m.AllKeys())
	m.Remove("a")
	fmt.Println(m.AllKeys())
}
