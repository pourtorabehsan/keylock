package keylock_test

import (
	"fmt"
	"sync"
	"time"

	"github.com/pourtorabehsan/keylock"
)

func ExampleKeyLock_Lock() {
	kl := keylock.New()
	key := "resource"

	// Locking a resource
	unlock := kl.Lock(key)
	fmt.Println("Resource locked")

	// Simulating concurrent access
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		unlock2 := kl.Lock(key)
		fmt.Println("Gained access to resource in goroutine")
		unlock2()
	}()

	// Wait a bit and unlock the resource
	time.Sleep(100 * time.Millisecond)
	unlock()
	fmt.Println("Resource unlocked")
	wg.Wait()

	// Output:
	// Resource locked
	// Resource unlocked
	// Gained access to resource in goroutine
}

func ExampleKeyLock_LockWithTimeout() {
	kl := keylock.New()
	key := "resource"

	// Locking a resource
	unlock, err := kl.LockWithTimeout(key, keylock.Immediate)
	if err != nil {
		panic(err)
	}
	fmt.Println("Resource locked")

	// Trying to lock with a timeout
	unlock2, err := kl.LockWithTimeout(key, 100*time.Millisecond)
	if err == keylock.ErrTimeout {
		fmt.Println("Failed to lock resource: timeout")
	}

	// Unlock the resource and try again
	unlock()
	fmt.Println("Resource unlocked")

	unlock2, err = kl.LockWithTimeout(key, 100*time.Millisecond)
	if err == nil {
		fmt.Println("Successfully locked resource with timeout")
		unlock2()
	}

	// Output:
	// Resource locked
	// Failed to lock resource: timeout
	// Resource unlocked
	// Successfully locked resource with timeout
}
