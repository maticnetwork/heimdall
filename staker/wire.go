package staker

import "github.com/cosmos/cosmos-sdk/wire"

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgCreateMaticValidator{}, "staker/MsgCreateMaticValidator", nil)
}

var cdcEmpty = wire.NewCodec()

func init() {
	RegisterWire(cdcEmpty)
	wire.RegisterCrypto(cdcEmpty)
}
