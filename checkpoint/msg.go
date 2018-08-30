package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

var cdc = wire.NewCodec()
const MsgType = "checkpoint"

var _ sdk.Msg = &MsgCheckpoint{}

type MsgCheckpoint struct {
	// TODO variable as we dont know who will call this
	Proposer 		sdk.AccAddress 	`json:"address"` // address of the validator owner
	CheckpointData []BlockHeader 	`json:"checkpointData"`
}
type BlockHeader struct {
	BlockHash string `json:"blockhash"`
	TxRoot string `json:"tx_root"`
	ReceiptRoot string `json:"receipt_root"`
}
//

func NewMsgSideBlock(proposer sdk.AccAddress,checkpointdata []BlockHeader) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer: proposer,
		CheckpointData:checkpointdata,
	}
}

//nolint
func (msg MsgCheckpoint) Type() string              { return MsgType }
func (msg MsgCheckpoint) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proposer} }

// get the bytes for the message signer to sign on
func (msg MsgCheckpoint) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// quick validity check
func (msg MsgCheckpoint) ValidateBasic() sdk.Error {
	if msg.Proposer == nil {
		//TODO create error and return respective error here, right now it will allow nil
		//return ErrBadValidatorAddr(DefaultCodespace)
		return nil
	}
	return nil
}

