package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	delPk1   = secp256k1.GenPrivKey().PubKey()
	delAddr1 = hmTypes.BytesToHeimdallAddress(delPk1.Address().Bytes())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()
	acc := types.NewBaseAccountWithAddress(delAddr1)
	globalAccNumber := uint64(10)

	kvPairs := []sdk.KVPair{
		{Key: types.AddressStoreKey(delAddr1), Value: cdc.MustMarshalBinaryBare(acc)},
		{Key: types.GlobalAccountNumberKey, Value: cdc.MustMarshalBinaryBare(globalAccNumber)},
		{Key: []byte{0x99}, Value: []byte{0x99}},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Account", fmt.Sprintf("%v\n%v", acc, acc)},
		{"GlobalAccNumber", fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumber, globalAccNumber)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
