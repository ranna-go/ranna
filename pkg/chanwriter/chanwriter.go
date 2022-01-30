package chanwriter

import (
	"io"
)

type Chanwriter struct {
	c chan []byte
}

var _ io.Writer = (*Chanwriter)(nil)

func New(c chan []byte) Chanwriter {
	return Chanwriter{c}
}

func (w Chanwriter) Write(p []byte) (n int, err error) {
	n = len(p)
	cp := make([]byte, n)
	copy(cp, p)
	w.c <- cp
	return
}
