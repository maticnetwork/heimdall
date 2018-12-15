package checkpoint

import (
	"github.com/cosmos/cosmos-sdk/codec"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func RegisterWire(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCheckpoint{}, "checkpoint/MsgCheckpoint", nil)
	cdc.RegisterConcrete(MsgCheckpointAck{}, "checkpoint/MsgCheckpointACK", nil)
	cdc.RegisterConcrete(MsgCheckpointNoAck{}, "checkpoint/MsgCheckpointNoACK", nil)
}

func RegisterPulp(pulp *hmTypes.Pulp) {
	pulp.RegisterConcrete(MsgCheckpoint{})
	pulp.RegisterConcrete(MsgCheckpointAck{})
	pulp.RegisterConcrete(MsgCheckpointNoAck{})
}

var cdcEmpty = codec.New()

func init() {
	RegisterWire(cdcEmpty)
	codec.RegisterCrypto(cdcEmpty)
}
