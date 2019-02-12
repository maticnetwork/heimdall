package staking

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

var cdc = codec.New()

const StakingRoute = "staking"

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

type MsgValidatorJoin struct {
	ID           types.ValidatorID `json:"ID"`
	SignerPubKey types.PubKey      `json:"pubKey"`
	StartEpoch   uint64            `json:"startEpoch"`
	EndEpoch     uint64            `json:"endEpoch"`
	Amount       json.Number       `json:"amount"`
}

func NewMsgValidatorJoin(
	_id uint64,
	_pubkey types.PubKey,
	_startEpoch uint64,
	_endEpoch uint64,
	_amount json.Number,
) MsgValidatorJoin {
	return MsgValidatorJoin{
		ID:           types.NewValidatorID(_id),
		SignerPubKey: _pubkey,
		StartEpoch:   _startEpoch,
		EndEpoch:     _endEpoch,
		Amount:       _amount,
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
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if bytes.Equal(msg.SignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.SignerPubKey.String())
	}

	r, _ := regexp.Compile("[0-9]+")
	if msg.Amount == "" || !r.MatchString(msg.Amount.String()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid new amount %v", msg.Amount.String())
	}

	return nil
}

func (msg MsgValidatorJoin) GetPower() uint64 {
	return types.GetValidatorPower(msg.Amount.String())
}

//
// validator update
//

var _ sdk.Msg = &MsgSignerUpdate{}

// MsgSignerUpdate signer update struct
// TODO add old signer sig check
type MsgSignerUpdate struct {
	ID              types.ValidatorID `json:"ID"`
	NewSignerPubKey types.PubKey      `json:"pubKey"`
	NewAmount       json.Number       `json:"amount"`
}

func NewMsgValidatorUpdate(_id uint64, pubKey types.PubKey, amount json.Number) MsgSignerUpdate {
	return MsgSignerUpdate{
		ID:              types.NewValidatorID(_id),
		NewSignerPubKey: pubKey,
		NewAmount:       amount,
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
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if bytes.Equal(msg.NewSignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.NewSignerPubKey.String())
	}

	r, _ := regexp.Compile("[0-9]+")
	if msg.NewAmount != "" && !r.MatchString(msg.NewAmount.String()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid new amount %v", msg.NewAmount.String())
	}

	return nil
}

func (msg MsgSignerUpdate) GetNewPower() uint64 {
	return types.GetValidatorPower(msg.NewAmount.String())
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

type MsgValidatorExit struct {
	ID types.ValidatorID `json:"ID"`
}

func NewMsgValidatorExit(_id uint64) MsgValidatorExit {
	return MsgValidatorExit{
		ID: types.NewValidatorID(_id),
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
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	return nil
}
