package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/eth/go-ethereum/common"
)

var cdc = codec.New()

// MsgType represents string for message type
const ValidatorJoin = "validatorJoin"

var _ sdk.Msg = &MsgValidatorJoin{}

type MsgValidatorJoin struct {
	ValidatorAddr common.Address `json:"validatorAddr"`
	StartEpoch    uint64         `json:"startEpoch"`
	EndEpoch      uint64         `json:"endEpoch"`
}

func NewMsgValidatorJoin(validatorAddr common.Address, startEpoch uint64, endEpoch uint64) MsgValidatorJoin {
	return MsgValidatorJoin{
		ValidatorAddr: validatorAddr,
		StartEpoch:    startEpoch,
		EndEpoch:      endEpoch,
	}
}

func (msg MsgValidatorJoin) Type() string {
	return ValidatorJoin
}

func (msg MsgValidatorJoin) Route() string { return ValidatorJoin }

// GetSigners returns address of the signer
func (msg MsgValidatorJoin) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

// GetSignBytes returns the bytes for the message signer to sign on
func (msg MsgValidatorJoin) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic checks quick validation
func (msg MsgValidatorJoin) ValidateBasic() sdk.Error {

	return nil
}

//
// validator update
//

// MsgType represents string for message type
const ValidatorUpdateSigner = "validatorUpdateSigner"

var _ sdk.Msg = &MsgValidatorUpdate{}

type MsgValidatorUpdate struct {
	CurrentValPubkey string
	NewValPubkey     string
}

func NewMsgValidatorUpdate(currentValPubKey string, newValPubkey string) MsgValidatorUpdate {
	return MsgValidatorUpdate{
		CurrentValPubkey: currentValPubKey,
		NewValPubkey:     newValPubkey,
	}
}

func (msg MsgValidatorUpdate) Type() string {
	return ValidatorUpdateSigner
}

func (msg MsgValidatorUpdate) Route() string { return ValidatorUpdateSigner }

// GetSigners returns address of the signer
func (msg MsgValidatorUpdate) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

// GetSignBytes returns the bytes for the message signer to sign on
func (msg MsgValidatorUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic checks quick validation
func (msg MsgValidatorUpdate) ValidateBasic() sdk.Error {

	return nil
}
