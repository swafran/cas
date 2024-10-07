package cas

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

type Locker interface {
	Lock()
	Unlock()
}

type Exclusive struct {
	locker Locker
	vals   map[string]int
}

func (e *Exclusive) incr(key string) {
	e.locker.Lock()
	defer e.locker.Unlock()
	e.vals[key]++
}

func (e *Exclusive) IncrMultiple(wg *sync.WaitGroup, key string, n int) {
	for i := 0; i < n; i++ {
		e.incr(key)
	}
	wg.Done()
}

func runTest(e *Exclusive) {
	e.vals = map[string]int{"a": 0, "b": 0}

	var wg sync.WaitGroup
	wg.Add(3)
	go e.IncrMultiple(&wg, "a", 10000)
	go e.IncrMultiple(&wg, "a", 10000)
	go e.IncrMultiple(&wg, "b", 10000)

	wg.Wait()
	fmt.Println(e.vals)
}

func TestRaceCondition(t *testing.T) {
	exclusive := &Exclusive{
		locker: &NullLock{},
		vals:   map[string]int{},
	}
	runTest(exclusive)
}

func TestCAS(t *testing.T) {
	exclusive := &Exclusive{
		locker: NewSpinLock(),
		vals:   map[string]int{},
	}
	runTest(exclusive)
}

func TestSyncMutex(t *testing.T) {
	exclusive := &Exclusive{
		locker: &sync.Mutex{},
		vals:   map[string]int{},
	}
	runTest(exclusive)
}

func TestGoSched(t *testing.T) {
	runtime.GOMAXPROCS(1)
	exclusive := &Exclusive{
		locker: NewSpinLock(),
		vals:   map[string]int{},
	}
	runTest(exclusive)
}
