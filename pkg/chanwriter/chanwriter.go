// Package chanwriter provides an io.Writer
// implementation which writes into a channel.
package chanwriter

import (
	"io"
)

// Chanwriter is an io.Writer implementation
// which writes the passed byte slice into the
// specified channel.
type Chanwriter struct {
	c chan<- []byte
}

var _ io.Writer = (*Chanwriter)(nil)

// New creates a new Chanwriter wrapping the
// passed channel.
func New(c chan<- []byte) Chanwriter {
	return Chanwriter{c}
}

func (w Chanwriter) Write(p []byte) (n int, err error) {
	n = len(p)
	cp := make([]byte, n)
	copy(cp, p)
	w.c <- cp
	return
}
