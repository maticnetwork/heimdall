package sideBlock


import "github.com/cosmos/cosmos-sdk/wire"

func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterConcrete(MsgSideBlock{}, "sideBlock/MsgSideBlock", nil)
}

var cdcEmpty = wire.NewCodec()


func init() {
	RegisterWire(cdcEmpty)
	wire.RegisterCrypto(cdcEmpty)
}
