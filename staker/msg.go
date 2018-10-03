package staker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"

	"github.com/tendermint/tendermint/crypto"
)

var cdc = wire.NewCodec()

// name to identify transaction types
const MsgType = "MsgCreateMaticValidator"

// verify interface at compile time
var _ sdk.Msg = &MsgCreateMaticValidator{}


// MsgUnrevoke - struct for unrevoking revoked validator
type MsgCreateMaticValidator struct {
	trialAddress 	sdk.AccAddress `json:"trialAddress"`
	// TODO variable as we dont know who will call this
	ValidatorAddress crypto.Address `json:"address"` // address of the validator owner
	//  TODO we can add multiple block details here , starting with string here
	Pubkey crypto.PubKey `json:"pubkey"`
	Power int64 `json:"power"`
}

func NewCreateMaticValidator(trial sdk.AccAddress,valAddress crypto.Address,pubkey crypto.PubKey,power int64) MsgCreateMaticValidator {
	return MsgCreateMaticValidator{
		trialAddress:trial,
		ValidatorAddress:valAddress,
		Pubkey:pubkey,
		Power:power,
	}
}

//nolint
func (msg MsgCreateMaticValidator) Type() string              { return MsgType }
func (msg MsgCreateMaticValidator) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.trialAddress} }

// get the bytes for the message signer to sign on
func (msg MsgCreateMaticValidator) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgCreateMaticValidator) ValidateBasic() sdk.Error {

	return nil
}

