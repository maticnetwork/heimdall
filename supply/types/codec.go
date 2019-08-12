package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*ModuleAccountInterface)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "cosmos-sdk/ModuleAccount", nil)
}

// RegisterPulp register pulp
func RegisterPulp(pulp *authTypes.Pulp) {
	// pulp.RegisterConcrete(ModuleAccount{})
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	ModuleCdc = cdc.Seal()
}
