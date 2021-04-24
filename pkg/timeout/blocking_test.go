package timeout

import (
	"testing"
	"time"
)

func TestRunBlockingWithTimeout(t *testing.T) {
	wasExecuted := false
	timedOut := RunBlockingWithTimeout(func() {
		wasExecuted = true
	}, time.Second)
	if !wasExecuted {
		t.Error("function was not executed")
	}
	if timedOut {
		t.Error("timed out even it should not have been")
	}

	timedOut = RunBlockingWithTimeout(func() {
		time.Sleep(100 * time.Millisecond)
	}, time.Millisecond)
	if !timedOut {
		t.Error("did not time out")
	}

	timedOut = RunBlockingWithTimeout(nil, time.Second)
	if timedOut {
		t.Error("timed out even if function was not executed")
	}
}
