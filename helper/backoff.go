package helper

import (
	"fmt"
	"time"
)

// ExponentialBackoff performs exponential backoff attempts on a given action
func ExponentialBackoff(action func() error, max int, initial time.Duration) error {
	if max < 0 {
		return fmt.Errorf("max number should be more or equal to zero, but %d given", max)
	}

	wait := initial
	var err error
	for i := 0; i < max; i++ {
		if err = action(); err == nil {
			return nil
		}
		time.Sleep(wait)
		wait *= 2
	}
	return err
}
