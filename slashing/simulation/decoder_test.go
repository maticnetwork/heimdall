package simulation

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// nolint:deadcode,unused,varcheck
var (
	delPk1    = secp256k1.GenPrivKey().PubKey()
	consAddr1 = hmTypes.BytesToHeimdallAddress(delPk1.Address().Bytes())
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

	info := hmTypes.NewValidatorSigningInfo(1, 0, 1, 0)
	missed := gogotypes.BoolValue{Value: true}

	kvPairs := []sdk.KVPair{
		{Key: types.GetValidatorSigningInfoKey(consAddr1.Bytes()), Value: cdc.MustMarshalBinaryBare(info)},
		{Key: types.GetValidatorMissedBlockBitArrayKey(consAddr1.Bytes(), 6), Value: cdc.MustMarshalBinaryBare(&missed)},
		{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"ValidatorSigningInfo", fmt.Sprintf("%v\n%v", info, info)},
		{"ValidatorMissedBlockBitArray", fmt.Sprintf("missedA: %v\nmissedB: %v", missed.Value, missed.Value)},
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
