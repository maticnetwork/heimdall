package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterInterface((*GenesisAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "auth/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)
}

// RegisterBaseAccount most users shouldn't use this, but this comes in handy for tests.
func RegisterBaseAccount(cdc *codec.Codec) {
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	codec.RegisterCrypto(cdc)
}

// ModuleCdc module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
