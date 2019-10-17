package base

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

func (s *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		runtime.Gosched()
	}
}

func (s *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(s), 0)
}

func SpinLock() sync.Locker {
	return new(spinLock)
}
