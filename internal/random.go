package internal

import (
	"sync"
	"math/rand"
)

var mu sync.Mutex

func RandomUint64() uint64 {
	mu.Lock()
	defer mu.Unlock()
	return rand.Uint64()
}
