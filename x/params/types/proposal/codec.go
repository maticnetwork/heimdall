package proposal

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/maticnetwork/heimdall/x/gov/types"
)

// RegisterLegacyAminoCodec registers all necessary param module types with a given LegacyAmino codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&ParameterChangeProposal{}, "heimdall/ParameterChangeProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ParameterChangeProposal{},
	)
}
