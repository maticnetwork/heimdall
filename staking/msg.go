package staking

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
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
	TxHash       common.Hash       `json:"tx_hash"`
	TimeStamp    uint64            `json:"timestamp"`
}

func NewMsgValidatorJoin(_id uint64, _pubkey types.PubKey, txhash common.Hash, timestamp uint64) MsgValidatorJoin {
	return MsgValidatorJoin{
		ID:           types.NewValidatorID(_id),
		SignerPubKey: _pubkey,
		TxHash:       txhash,
		TimeStamp:    timestamp,
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

	return nil
}

//
// validator update
//
var _ sdk.Msg = &MsgSignerUpdate{}

// MsgSignerUpdate signer update struct
type MsgSignerUpdate struct {
	ID              types.ValidatorID `json:"ID"`
	NewSignerPubKey types.PubKey      `json:"pubKey"`
	TxHash          common.Hash       `json:"tx_hash"`
	TimeStamp       uint64            `json:"timestamp"`
}

func NewMsgValidatorUpdate(_id uint64, pubKey types.PubKey, txhash common.Hash, timestamp uint64) MsgSignerUpdate {
	return MsgSignerUpdate{
		ID:              types.NewValidatorID(_id),
		NewSignerPubKey: pubKey,
		TxHash:          txhash,
		TimeStamp:       timestamp,
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

	return nil
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

type MsgValidatorExit struct {
	ID        types.ValidatorID `json:"ID"`
	TxHash    common.Hash       `json:"tx_hash"`
	TimeStamp uint64            `json:"timestamp"`
}

func NewMsgValidatorExit(_id uint64, txhash common.Hash, timestamp uint64) MsgValidatorExit {
	return MsgValidatorExit{
		ID:        types.NewValidatorID(_id),
		TxHash:    txhash,
		TimeStamp: timestamp,
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

//
// Update power
//
var _ sdk.Msg = &MsgPowerUpdate{}

type MsgPowerUpdate struct {
	ID        types.ValidatorID `json:"ID"`
	TxHash    common.Hash       `json:"tx_hash"`
	TimeStamp uint64            `json:"timestamp"`
}

func NewMsgPowerUpdate(_id uint64, txhash common.Hash, timestamp uint64) MsgPowerUpdate {
	return MsgPowerUpdate{
		ID:        types.NewValidatorID(_id),
		TxHash:    txhash,
		TimeStamp: timestamp,
	}
}

func (msg MsgPowerUpdate) Type() string {
	return "power-update"
}

func (msg MsgPowerUpdate) Route() string {
	return StakingRoute
}

func (msg MsgPowerUpdate) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgPowerUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgPowerUpdate) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	return nil
}
