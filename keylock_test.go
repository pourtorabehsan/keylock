package keylock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/pourtorabehsan/keylock"
)

func TestLock(t *testing.T) {
	kl := keylock.New()
	key := "test-key"

	unlock := kl.Lock(key)

	var wg sync.WaitGroup
	wg.Add(1)

	started := make(chan struct{})
	go func() {
		defer wg.Done()
		started <- struct{}{}
		unlock2 := kl.Lock(key)
		unlock2()
	}()

	<-started

	unlock()
	wg.Wait()
}

func TestLockAndUnlock(t *testing.T) {
	kl := keylock.New()
	key := "test-key"

	unlock := kl.Lock(key)

	acquired := false
	started := make(chan struct{})

	go func() {
		started <- struct{}{}
		unlock2 := kl.Lock(key)
		defer unlock2()
		acquired = true
	}()

	<-started
	time.Sleep(100 * time.Millisecond) // Give the goroutine time to lock

	if acquired {
		t.Errorf("Expected lock to be held, but it was released")
	}

	unlock()

	time.Sleep(100 * time.Millisecond) // Give the goroutine time to lock

	if !acquired {
		t.Errorf("Expected lock to be released, but it was held")
	}
}

func TestImmediateLock(t *testing.T) {
	kl := keylock.New()
	key := "test-key"

	unlock, err := kl.LockWithTimeout(key, keylock.Immediate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if unlock == nil {
		t.Errorf("Expected unlock function, got nil")
	}
	defer unlock()

	// the lock is already held, so this should fail
	_, err = kl.LockWithTimeout(key, keylock.Immediate)
	if err != keylock.ErrTimeout {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

func TestLockWithTimeout(t *testing.T) {
	kl := keylock.New()
	key := "test-key"

	unlock, err := kl.LockWithTimeout(key, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// the lock is already held, so this should fail
	_, err = kl.LockWithTimeout(key, 100*time.Millisecond)
	if err != keylock.ErrTimeout {
		t.Errorf("Expected timeout error, got %v", err)
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		unlock()
	}()

	// The lock should be released after 200 milliseconds
	//  in the goroutine above, so this should succeed.
	unlock2, err := kl.LockWithTimeout(key, 1*time.Second)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	unlock2()
}

func TestLocksAreIndependent(t *testing.T) {
	kl := keylock.New()
	key1 := "test-key1"
	key2 := "test-key2"

	unlock1 := kl.Lock(key1)

	var wg sync.WaitGroup
	wg.Add(1)

	started := make(chan struct{})

	go func() {
		defer wg.Done()
		started <- struct{}{}
		unlock2 := kl.Lock(key2)
		unlock2()
	}()

	<-started

	// The goroutine should be able to lock key and unlock even though key1 is locked
	wg.Wait()
	unlock1()
}
