package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// TODO we most likely dont need to register to amino as we are using RLP to encode

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgValidatorJoin{}, "staking/MsgValidatorJoin", nil)
	cdc.RegisterConcrete(MsgSignerUpdate{}, "staking/MsgSignerUpdate", nil)
	cdc.RegisterConcrete(MsgValidatorExit{}, "staking/MsgValidatorExit", nil)
	cdc.RegisterConcrete(MsgStakeUpdate{}, "staking/MsgStakeUpdate", nil)
}

func RegisterPulp(pulp *authTypes.Pulp) {
	pulp.RegisterConcrete(MsgValidatorJoin{})
	pulp.RegisterConcrete(MsgSignerUpdate{})
	pulp.RegisterConcrete(MsgValidatorExit{})
	pulp.RegisterConcrete(MsgStakeUpdate{})
}

var cdcEmpty = codec.New()

func init() {
	RegisterCodec(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
