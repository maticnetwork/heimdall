package checkpoint

import "github.com/cosmos/cosmos-sdk/wire"

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCheckpoint{}, "checkpoint/MsgCheckpoint", nil)
}

var cdcEmpty = wire.NewCodec()

func init() {
	RegisterWire(cdcEmpty)
	wire.RegisterCrypto(cdcEmpty)
}
