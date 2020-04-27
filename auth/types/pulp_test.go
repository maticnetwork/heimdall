package types

import (
	"testing"

	assert "github.com/attic-labs/testify/require"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetPulpHash(t *testing.T) {
	tc := struct {
		in  sdk.Msg
		out []byte
	}{
		in:  sdk.NewTestMsg(nil),
		out: []byte{142, 88, 179, 79},
	}
	out := GetPulpHash(tc.in)
	assert.Equal(t, string(tc.out), string(out))
}
