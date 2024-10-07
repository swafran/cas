package cas

import (
	"sync/atomic"
)

const free = int32(0)

type SpinLock struct {
	state *int32
}

// Lock provides mutually exclusive access
// 42 est une valeur arbitraire, ça peut être n'importe quoi sauf 0 (la valeur de free)
// Avant go(1.4?) un appel à Gosched() était néc. sur un loop aussi maigre pour éviter lockups si sur un seul core
func (l *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(l.state, free, 42) {
		//runtime.Gosched()
	}
}

func (l *SpinLock) Unlock() {
	atomic.StoreInt32(l.state, free)
}

func NewSpinLock() *SpinLock {
	startState := free
	return &SpinLock{state: &startState}
}

type NullLock struct{}

func (n *NullLock) Lock()   {}
func (n *NullLock) Unlock() {}
