The jitter package simulates an unreliable connection, taking an io.ReadWriter and a configuration describing the frequency of errors, and returning an io.ReadWriter.

# Example
```go
package main

import (
	"net"

	"github.com/ebusto/jitter"
)

func main() {
	a, b := net.Pipe()

	j := jitter.New(a, &jitter.Config{
		Delay:  5 * time.Second, // Read/Write will take between 0s and 5s.
		Random: 0.05,            // 5% of bytes will be replaced with a random byte.
		Skip:   0.01,            // 1% of bytes will be lost.
	})

	n, err := j.Read(buf)  // Reads from j will be corrupted.
	n, err := j.Write(buf) // Writes to j will be corrupted.
	...
}
```
