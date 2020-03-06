package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// ModuleCdc module codec
var ModuleCdc = codec.New()

// RegisterCodec registers all the necessary types and interfaces for
// governance.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Content)(nil), nil)

	cdc.RegisterConcrete(MsgSubmitProposal{}, "heimdall/MsgSubmitProposal", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "heimdall/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "heimdall/MsgVote", nil)

	cdc.RegisterConcrete(TextProposal{}, "heimdall/TextProposal", nil)
	cdc.RegisterConcrete(SoftwareUpgradeProposal{}, "heimdall/SoftwareUpgradeProposal", nil)
}

// RegisterProposalTypeCodec registers an external proposal content type defined
// in another module for the internal ModuleCdc. This allows the MsgSubmitProposal
// to be correctly Amino encoded and decoded.
func RegisterProposalTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

// RegisterPulp register pulp
func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgSubmitProposal{})
	pulp.RegisterConcrete(MsgDeposit{})
	pulp.RegisterConcrete(MsgVote{})
}

// TODO determine a good place to seal this codec
func init() {
	RegisterCodec(ModuleCdc)
}
