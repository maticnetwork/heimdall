package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterInterface((*auth.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)
	cdc.RegisterConcrete(&auth.BaseVestingAccount{}, "auth/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&auth.ContinuousVestingAccount{}, "auth/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&auth.DelayedVestingAccount{}, "auth/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(StdTx{}, "auth/StdTx", nil)
}

// RegisterBaseAccount most users shouldn't use this, but this comes in handy for tests.
func RegisterBaseAccount(cdc *codec.Codec) {
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterInterface((*auth.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	cdc.RegisterConcrete(&auth.BaseVestingAccount{}, "cosmos-sdk/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&auth.ContinuousVestingAccount{}, "cosmos-sdk/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&auth.DelayedVestingAccount{}, "cosmos-sdk/DelayedVestingAccount", nil)
	codec.RegisterCrypto(cdc)
}

// MsgCdc module wide codec
var MsgCdc *codec.Codec

func init() {
	MsgCdc = codec.New()
	RegisterCodec(MsgCdc)
	codec.RegisterCrypto(MsgCdc)
	MsgCdc.Seal()
}
