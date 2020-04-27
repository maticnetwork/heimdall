package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	exported "github.com/maticnetwork/heimdall/auth/exported"
	"github.com/maticnetwork/heimdall/auth/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding auth type
func DecodeStore(cdc *codec.Codec, kvA, kvB sdk.KVPair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.AddressStoreKeyPrefix):
		var accA, accB exported.Account
		cdc.MustUnmarshalBinaryBare(kvA.Value, &accA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &accB)
		return fmt.Sprintf("%v\n%v", accA, accB)

	case bytes.Equal(kvA.Key, types.GlobalAccountNumberKey):
		var globalAccNumberA, globalAccNumberB uint64
		cdc.MustUnmarshalBinaryBare(kvA.Value, &globalAccNumberA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &globalAccNumberB)
		return fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumberA, globalAccNumberB)

	default:
		panic(fmt.Sprintf("invalid account key %X", kvA.Key))
	}
}
