package settings

import (
	t "github.com/faelmori/kubex-interfaces/types"
	"sync"
)

// KubexThreading is a struct that holds the mutexes for the spider
type KubexThreading struct {
	// IThreading interface for threading
	t.IThreading
	// Mutex for thread safety
	mu sync.RWMutex
	// SyncGroup for synchronization
	wg sync.WaitGroup
}

// NewThreading creates a new Threading instance
func NewThreading() *KubexThreading {
	return &KubexThreading{
		mu: sync.RWMutex{},
		wg: sync.WaitGroup{},
	}
}

// TryLock tries to lock the mutex
func (k *KubexThreading) TryLock() bool { return k.mu.TryLock() }

// Lock locks the mutex
func (k *KubexThreading) Lock() { k.mu.Lock() }

// Unlock unlocks the mutex
func (k *KubexThreading) Unlock() { k.mu.Unlock() }

// RLock locks the mutex for reading
func (k *KubexThreading) RLock() { k.mu.RLock() }

// RUnlock unlocks the mutex for reading
func (k *KubexThreading) RUnlock() { k.mu.RUnlock() }

// LockFunc locks the mutex and executes the function
func (k *KubexThreading) LockFunc(f func()) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if f == nil {
		return
	}
	f()
}

// UnlockFunc unlocks the mutex and executes the function
func (k *KubexThreading) UnlockFunc(f func()) {
	k.mu.Unlock()
	defer k.mu.Lock()
	if f == nil {
		return
	}
	f()
}

// LockFuncWithArgs locks the mutex and executes the function with args
func (k *KubexThreading) LockFuncWithArgs(f func(interface{}), args interface{}) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if f == nil {
		return
	}
	f(args)
}

// UnlockFuncWithArgs unlocks the mutex and executes the function with args
func (k *KubexThreading) UnlockFuncWithArgs(f func(interface{}), args interface{}) {
	k.mu.Unlock()
	defer k.mu.Lock()
	if f == nil {
		return
	}
	f(args)
}

// Defer returns a function that will be executed when the mutex is unlocked
func (k *KubexThreading) Defer() t.IDeferFunc {
	return func(f func(), args interface{}) func() error {
		return func() error {
			k.mu.Unlock()
			if f != nil {
				f()
			}
			return nil
		}
	}
}

// Add adds delta to the WaitGroup counter
func (k *KubexThreading) Add(delta int) { k.wg.Add(delta) }

// Wait waits for the WaitGroup counter to reach zero
func (k *KubexThreading) Wait() { k.wg.Wait() }

// Done decrements the WaitGroup counter
func (k *KubexThreading) Done() { k.wg.Done() }

// WaitGroup returns the WaitGroup
func (k *KubexThreading) WaitGroup() *sync.WaitGroup { return &k.wg }

// WaitGroupAdd adds delta to the WaitGroup counter
func (k *KubexThreading) WaitGroupAdd(delta int) { k.wg.Add(delta) }

// WaitGroupDone decrements the WaitGroup counter
func (k *KubexThreading) WaitGroupDone() { k.wg.Done() }
