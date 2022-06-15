package helper

import "time"

// ExponentialBackoff performs exponential backoff attempts on a given action
func ExponentialBackoff(action func() error, max int, initial time.Duration) error {
	wait := initial
	var err error
	for i := 0; i < max; i++ {
		err = action()
		if err == nil {
			return nil
		}
		time.Sleep(wait)
		wait *= 2
	}
	return err
}
