package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maticnetwork/heimdall/sidechannel/types"
)

func TestCodec(t *testing.T) {
	require.NotNil(t, types.ModuleCdc, "ModuleCdc shouldn't be nil")
}
