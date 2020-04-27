package app

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
)

func TestGetSimulationLog(t *testing.T) {
	cdc := MakeCodec()

	decoders := make(hmModule.StoreDecoderRegistry)
	decoders[authTypes.StoreKey] = func(cdc *codec.Codec, kvAs, kvBs sdk.KVPair) string { return "10" }

	tests := []struct {
		store       string
		kvPairs     []sdk.KVPair
		expectedLog string
	}{
		{
			"Empty",
			[]sdk.KVPair{{}},
			"",
		},
		{
			authTypes.StoreKey,
			[]sdk.KVPair{{Key: authTypes.GlobalAccountNumberKey, Value: cdc.MustMarshalBinaryBare(uint64(10))}},
			"10",
		},
		{
			"OtherStore",
			[]sdk.KVPair{{Key: []byte("key"), Value: []byte("value")}},
			fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", []byte("key"), []byte("value"), []byte("key"), []byte("value")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.store, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, GetSimulationLog(tt.store, decoders, cdc, tt.kvPairs, tt.kvPairs), tt.store)
		})
	}
}
