# KeyLock

KeyLock is a lightweight Go package for managing concurrency with a key-based locking mechanism. It allows you to easily manage access to shared resources by multiple goroutines in a safe and efficient manner.

## Features

- Key-based locks
- Timeout feature for lock acquisition
- Simple and intuitive API
- High-performance locking mechanism
- Thoroughly tested with edge cases in mind

## Installation

Install KeyLock with:

```bash
go get github.com/pourtorabehsan/keylock
```

## Usage

```go
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
```

## Testing

Run the tests with:

```bash
go test ./...
```

## License

This project is licensed under the MIT License.
