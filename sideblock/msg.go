package sideBlock


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var cdc = wire.NewCodec()

// name to identify transaction types
const MsgType = "sideBlock"

// verify interface at compile time
var _ sdk.Msg = &MsgSideBlock{}

// MsgUnrevoke - struct for unrevoking revoked validator
type MsgSideBlock struct {
	// TODO variable as we dont know who will call this
	VariableAddress sdk.AccAddress `json:"address"` // address of the validator owner
	//  TODO we can add multiple block details here , starting with string here
	BlockHash string `json:"blockhash"`
	TxRoot string `json:"tx_root"`
	ReceiptRoot string `json:"receipt_root"`
	//BlockNumber big.Int `json:"block_number"`
}

func NewMsgSideBlock(variableAddr sdk.AccAddress,blockhash string,txroot string,rRoot string) MsgSideBlock {
	return MsgSideBlock{
		VariableAddress: variableAddr,
		BlockHash:      blockhash,
		TxRoot:			txroot,
		ReceiptRoot:    rRoot,
		//BlockNumber: 	number,
	}
}

//nolint
func (msg MsgSideBlock) Type() string              { return MsgType }
func (msg MsgSideBlock) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.VariableAddress} }

// get the bytes for the message signer to sign on
func (msg MsgSideBlock) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgSideBlock) ValidateBasic() sdk.Error {
	if msg.VariableAddress == nil {
		//TODO create error and return respective error here, right now it will allow nil
		//return ErrBadValidatorAddr(DefaultCodespace)
		return nil
	}
	return nil
}

