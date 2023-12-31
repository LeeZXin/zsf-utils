package randutil

import (
	"math/rand"
	"sync"
	"time"
)

var (
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu sync.Mutex
)

func Int() int {
	mu.Lock()
	defer mu.Unlock()
	return r.Int()
}

func Int63n(n int64) int64 {
	mu.Lock()
	defer mu.Unlock()
	return r.Int63n(n)
}

func Intn(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return r.Intn(n)
}

func Int31n(n int32) int32 {
	mu.Lock()
	defer mu.Unlock()
	return r.Int31n(n)
}

func Float64() float64 {
	mu.Lock()
	defer mu.Unlock()
	return r.Float64()
}

func Uint64() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return r.Uint64()
}
