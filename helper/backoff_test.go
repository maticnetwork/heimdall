package helper

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		i := 0
		outcomes := []bool{false, false, true}
		t0 := time.Now()
		err := ExponentialBackoff(func() error {
			outcome := outcomes[i]
			i++
			if outcome {
				return nil
			}
			return errors.New("bad")
		}, 3, 150*time.Millisecond)

		elapsed := time.Since(t0)

		require.NoError(t, err)
		require.Equal(t, i, 3)
		require.True(t, elapsed >= 450*time.Millisecond)
	})

	t.Run("failed", func(t *testing.T) {
		t.Parallel()

		i := 0
		t0 := time.Now()
		err := ExponentialBackoff(func() error {
			i++
			return errors.New("bad")
		}, 3, 100*time.Millisecond)

		elapsed := time.Since(t0)

		require.Error(t, err)
		require.Equal(t, i, 3)
		require.True(t, elapsed >= 600*time.Millisecond)
	})
}
