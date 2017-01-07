package jitter

import (
	"io"
	"math/rand"
	"time"
)

type Config struct {
	Delay  time.Duration
	Random float64
	Skip   float64
}

type Jitter struct {
	cn io.ReadWriter
	cf *Config
}

var defaultConfig = &Config{
	Delay:  5 * time.Millisecond,
	Random: 0.001,
	Skip:   0.001,
}

func New(cn io.ReadWriter, cf *Config) *Jitter {
	if cf == nil {
		cf = defaultConfig
	}

	rand.Seed(time.Now().Unix())

	return &Jitter{cn, cf}
}

func (j *Jitter) Copy(dst []byte, src []byte) int {
	c := j.cf
	n := 0

	for _, b := range src {
		if rand.Float64() < c.Skip {
			continue
		}

		if rand.Float64() < c.Random {
			if _, err := rand.Read(dst[n : n+1]); err != nil {
				panic(err)
			}
		} else {
			dst[n] = b
		}

		n++
	}

	if c.Delay > 0 {
		time.Sleep(time.Duration(rand.Int63n(int64(c.Delay))))
	}

	return n
}

func (j *Jitter) Read(p []byte) (int, error) {
	b := make([]byte, len(p))

	n, err := j.cn.Read(b)

	return j.Copy(p, b[:n]), err
}

func (j *Jitter) Write(p []byte) (int, error) {
	b := make([]byte, len(p))

	n := j.Copy(b, p)
	_, err := j.cn.Write(b[:n])

	return len(p), err
}
