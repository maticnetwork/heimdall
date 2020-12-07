package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/provider module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/provider and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// RegisterLegacyAminoCodec register concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// no legacy amino needed
}

// RegisterInterfaces registers the x/sidechannel interfaces types with the interface registry
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// no interface needed for sidechannel
}
