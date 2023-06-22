package keylock

import (
	"errors"
	"sync"
	"time"
)

var ErrTimeout = errors.New("timeout")

// KeyLock is a key-based lock structure for managing concurrency.
type KeyLock struct {
	m     sync.Mutex
	locks map[string]chan struct{}
}

// New creates a new KeyLock instance.
func New() *KeyLock {
	return &KeyLock{locks: make(map[string]chan struct{})}
}

// Lock acquires a lock for the given key.
// It returns a function that releases the lock when called.
func (kl *KeyLock) Lock(key string) (unlock func()) {
	ch := kl.provide(key)
	ch <- struct{}{}
	return func() { <-ch }
}

// LockWithTimeout tries to acquire a lock for the given key within the
// specified timeout duration. If the lock is successfully acquired,
// it returns a function that releases the lock when called. If not,
// it returns an error.
func (kl *KeyLock) LockWithTimeout(key string, timeout time.Duration) (unlock func(), err error) {
	ch := kl.provide(key)
	select {
	case ch <- struct{}{}:
		return func() { <-ch }, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

// provide makes sure that a lock for the given key exists in the
// internal map of locks. If it doesn't exist, it creates a new lock.
func (kl *KeyLock) provide(key string) chan struct{} {
	kl.m.Lock()
	defer kl.m.Unlock()
	ch, ok := kl.locks[key]
	if !ok {
		ch = make(chan struct{}, 1)
		kl.locks[key] = ch
	}
	return ch
}
