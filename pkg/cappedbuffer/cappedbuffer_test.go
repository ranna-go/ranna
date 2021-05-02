package cappedbuffer

import "testing"

func TestWrite(t *testing.T) {
	buf := New([]byte{}, 0)
	n, err := buf.Write(make([]byte, 1000))
	if n != 1000 {
		t.Errorf("%d bytes written", n)
	}
	if err != nil {
		t.Error(err)
	}

	buf = New([]byte{}, -1)
	n, err = buf.Write(make([]byte, 1000))
	if n != 1000 {
		t.Errorf("%d bytes written", n)
	}
	if err != nil {
		t.Error(err)
	}

	buf = New([]byte{}, 1000)
	n, err = buf.Write(make([]byte, 950))
	if n != 950 {
		t.Errorf("%d bytes written", n)
	}
	if err != nil {
		t.Error(err)
	}

	n, err = buf.Write(make([]byte, 40))
	if n != 40 {
		t.Errorf("%d bytes written", n)
	}
	if err != nil {
		t.Error(err)
	}

	n, err = buf.Write(make([]byte, 20))
	if n != 0 {
		t.Errorf("%d bytes written", n)
	}
	if err == nil {
		t.Error("did not return error")
	}
}
