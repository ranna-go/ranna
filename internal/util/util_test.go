package util

import (
	"testing"
	"time"
)

func TestParseMemoryStr(t *testing.T) {
	expect := func(v string, toBe int64) {
		res, err := ParseMemoryStr(v)
		if err != nil {
			t.Error(err)
		}
		if res != toBe {
			t.Errorf("value was %d (exppected: %d)", res, toBe)
		}
	}

	expectErr := func(v string, errMsg string) {
		_, err := ParseMemoryStr(v)
		if err == nil {
			t.Error("no error was returned")
		}
		if errMsg != err.Error() {
			t.Errorf("invalid err msg: %s (expected: %s)",
				err.Error(), errMsg)
		}
	}

	expect("", 0)
	expect("50", 50)
	expect("50K", 50*1024)
	expect("50k", 50*1024)
	expect("50M", 50*1024*1024)
	expect("50m", 50*1024*1024)
	expect("50G", 50*1024*1024*1024)
	expect("50g", 50*1024*1024*1024)
	expect("50T", 50*1024*1024*1024*1024)
	expect("50t", 50*1024*1024*1024*1024)

	expect("50Kb", 50*1024)
	expect("50kb", 50*1024)
	expect("50Mb", 50*1024*1024)
	expect("50mb", 50*1024*1024)
	expect("50Gb", 50*1024*1024*1024)
	expect("50gb", 50*1024*1024*1024)
	expect("50Tb", 50*1024*1024*1024*1024)
	expect("50tb", 50*1024*1024*1024*1024)

	expect("50KB", 50*1024)
	expect("50kB", 50*1024)
	expect("50MB", 50*1024*1024)
	expect("50mB", 50*1024*1024)
	expect("50GB", 50*1024*1024*1024)
	expect("50gB", 50*1024*1024*1024)
	expect("50TB", 50*1024*1024*1024*1024)
	expect("50tB", 50*1024*1024*1024*1024)

	expectErr("50E", ErrInvalidSyntax.Error())
	expectErr("50e", ErrInvalidSyntax.Error())
	expectErr("k", ErrInvalidSyntax.Error())
	expectErr("K", ErrInvalidSyntax.Error())
	expectErr("poggers", ErrInvalidSyntax.Error())
	expectErr("1pK", "strconv.ParseInt: parsing \"1p\": invalid syntax")
}

func TestMeasureTime(t *testing.T) {
	d := MeasureTime(func() {})
	if d != 0 {
		t.Errorf("%s delay measured, expected: 0", d)
	}

	d = MeasureTime(func() {
		time.Sleep(1 * time.Second)
	})
	if d > 1*time.Second+10*time.Millisecond ||
		d < 1*time.Second-10*time.Millisecond {
		t.Errorf("%s delay measured, expected: ~1s", d)
	}
}
