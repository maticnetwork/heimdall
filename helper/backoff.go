package helper

import (
	"time"
)

// ExponentialBackoff performs exponential backoff attempts on a given action
func ExponentialBackoff(action func() error, maxV uint, wait time.Duration) error {
	var err error
	for i := uint(0); i < maxV; i++ {
		if err = action(); err == nil {
			break
		}

		time.Sleep(wait)
		wait *= 2
	}

	return err
}
