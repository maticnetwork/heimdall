package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfigFromFlags(t *testing.T) {
	t.Parallel()

	config := NewConfigFromFlags()
	require.Equal(t, config.ChainID, "", "Default chain id should be empty")
}
