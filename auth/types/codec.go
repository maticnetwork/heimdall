package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	// TODO figure out the way to not to make authTypes.ModuleCdc global codec
	ModuleCdc = cdc

	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "auth/Account", nil)
	cdc.RegisterConcrete(&GenesisAccount{}, "auth/GenesisAccount", nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)
}

// ModuleCdc module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
