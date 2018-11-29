package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var cdc = codec.New()

const StakingRoute = "staking"

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

type MsgValidatorJoin struct {
	ValidatorAddress common.Address `json:"address"`
	ValidatorPubKey  []byte         `json:"pubKey"`
}

func NewMsgValidatorJoin(Address common.Address, pubKey []byte) MsgValidatorJoin {
	return MsgValidatorJoin{
		ValidatorAddress: Address,
		ValidatorPubKey:  pubKey,
	}
}

func (msg MsgValidatorJoin) Type() string {
	return "validator-join"
}

func (msg MsgValidatorJoin) Route() string {
	return StakingRoute
}

func (msg MsgValidatorJoin) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgValidatorJoin) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorJoin) ValidateBasic() sdk.Error {
	// add length checks
	return nil
}

//
// validator update
//

var _ sdk.Msg = &MsgSignerUpdate{}

// MsgSignerUpdate signer update struct
type MsgSignerUpdate struct {
	ValidatorAddress   common.Address `json:"address"`
	NewValidatorPubKey []byte         `json:"pubKey"`
}

func NewMsgValidatorUpdate(address common.Address, pubKey []byte) MsgSignerUpdate {
	return MsgSignerUpdate{
		ValidatorAddress:   address,
		NewValidatorPubKey: pubKey,
	}
}

func (msg MsgSignerUpdate) Type() string {
	return "validator-update"
}

func (msg MsgSignerUpdate) Route() string {
	return StakingRoute
}

func (msg MsgSignerUpdate) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgSignerUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSignerUpdate) ValidateBasic() sdk.Error {
	// add length checks
	return nil
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

type MsgValidatorExit struct {
	ValidatorAddress common.Address
}

func NewMsgValidatorExit(address common.Address) MsgValidatorExit {
	return MsgValidatorExit{
		ValidatorAddress: address,
	}
}

func (msg MsgValidatorExit) Type() string {
	return "validator-exit"
}

func (msg MsgValidatorExit) Route() string {
	return StakingRoute
}

func (msg MsgValidatorExit) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgValidatorExit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorExit) ValidateBasic() sdk.Error {
	return nil
}
