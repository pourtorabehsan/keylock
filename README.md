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
package main

import (
	"time"

	"github.com/pourtorabehsan/keylock"
)

func main() {
    kl := keylock.New()

    // Acquire a lock for a key
    unlock := kl.Lock("my-key")

    // Release the lock
    defer unlock()

    // Do the work
    time.Sleep(1 * time.Second)
}
```

## Testing

Run the tests with:

```bash
go test ./...
```

## License

This project is licensed under the MIT License.
