package hashset

import (
	"fmt"
	"testing"
)

func TestNewConcurrentHashSet(t *testing.T) {
	set := NewConcurrentHashSet([]string{"x"})
	fmt.Println(set.Intersect(set))
}
