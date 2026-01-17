// Package timeout provides simple functionalities
// to timeout executions.
package timeout

import "time"

// RunBlockingWithTimeout runs the passed function
// in a new go routine and blocks the current go
// routine until it finishes or until the passed
// timeut duration exceeded.
//
// Returns true if the timeout duration exceeded.
// The function execution is not canceled after
// timeout.
func RunBlockingWithTimeout(f func(), timeout time.Duration) (v bool) {
	if f == nil {
		return false
	}

	cFinished := make(chan struct{}, 1)
	timer := time.NewTimer(timeout)

	go func() {
		f()
		cFinished <- struct{}{}
	}()

	select {
	case <-timer.C:
		return true
	case <-cFinished:
		return false
	}
}
