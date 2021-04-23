package util

import "time"

func RunBlockingWithTimeout(f func(), timeout time.Duration) (v bool) {
	cFinished := make(chan struct{}, 1)
	timer := time.NewTimer(timeout)

	go func() {
		f()
		cFinished <- struct{}{}
	}()

	select {
	case <-timer.C:
		v = true
	case <-cFinished:
	}

	return
}
