package jitter

import (
	"bytes"
	"io"
	"net"
	"testing"
)

func TestJitter(t *testing.T) {
	a, b := net.Pipe()

	j := New(a, nil)
	s := "The quick brown fox jumps over the lazy dog."

	t.Logf("Jitter read")
	testJitter(t, []byte(s), j, b)

	t.Logf("Jitter write")
	testJitter(t, []byte(s), b, j)
}

func testJitter(t *testing.T, v []byte, r io.Reader, w io.Writer) {
	done := make(chan bool)
	same := 0
	size := 100

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				if _, err := w.Write(v); err != nil {
					t.Fatal(err)
				}
			}
		}
	}()

	b := make([]byte, len(v))

	for i := 0; i < size; i++ {
		n, err := r.Read(b)

		if err != nil {
			t.Fatal(err)
		}

		t.Logf("[%d] [%t] %s\n", n, bytes.Equal(b, v), string(b))

		if bytes.Equal(b, v) {
			same++
		}
	}

	if same == size {
		t.Fatal("All bytes were modified! :-(")
	}

	if same == 0 {
		t.Fatal("No bytes were modified! :-(")
	}
}
