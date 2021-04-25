// Package cappedbuffer provides a simple bytes.Buffer
// wrapper with a total grow cap.
package cappedbuffer

import (
	"bytes"
	"errors"
)

// ErrBufferOverflow is returned if the write operation
// to the buffer would exceed the specified cap.
var ErrBufferOverflow = errors.New("buffer overflow")

// CappedBuffer wraps bytes.Buffer but with a fixed
// size the internal buffer can grow.
type CappedBuffer struct {
	*bytes.Buffer
	cap int
}

// New returns a new CappedBuffer consuming the passed
// buf array and with the specified grow cap.
//
// If a write operation would grow the size of the
// internal buffer beyond the specified cap, an
// ErrBufferOverflow is returned.
//
// If cap <= 0, the buffer will have no set cap.
func New(buf []byte, cap int) *CappedBuffer {
	return &CappedBuffer{
		Buffer: bytes.NewBuffer(buf),
		cap:    cap,
	}
}

func (cb *CappedBuffer) Write(p []byte) (n int, err error) {
	if cb.cap > 0 && cb.Len()+len(p) > cb.cap {
		return 0, ErrBufferOverflow
	}
	return cb.Buffer.Write(p)
}
