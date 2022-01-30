package chanwriter

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	c := make(chan []byte)
	cw := New(c)

	var rec []byte
	go func() {
		rec = <-c
	}()

	data := []byte("Hello world!")
	n, err := cw.Write(data)
	if err != nil {
		t.Error("Write failed: ", err)
	}
	if n != len(data) {
		t.Errorf("value of n does not equal expected value (%d != %d)", n, len(data))
	}

	if !bytes.Equal(rec, data) {
		t.Errorf("received data does not match (%x != %x)", rec, data)
	}
}
