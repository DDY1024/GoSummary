package rand

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type FatRand struct {
	idx       uint32
	length    uint32
	locks     []*sync.Mutex
	generator []*rand.Rand
}

func NewFatRand(n uint32) *FatRand {
	r := &FatRand{
		length:    n,
		locks:     make([]*sync.Mutex, n),
		generator: make([]*rand.Rand, n),
	}
	for i := uint32(0); i < n; i++ {
		r.locks[i] = new(sync.Mutex)
		// rand.NewSource 是非线程安全的，需要加锁进行保护
		r.generator[i] = rand.New(rand.NewSource(time.Now().Unix()))
	}
	return r
}

func (r *FatRand) Uint32() uint32 {
	i := atomic.AddUint32(&r.idx, 1)
	i %= r.length
	r.locks[i].Lock()
	val := r.generator[i].Uint32()
	r.locks[i].Unlock()
	return val
}
