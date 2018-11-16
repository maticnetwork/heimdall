package checkpoint

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

var cdc = codec.New()

//
// Checkpoint Msg
//

const Checkpoint = "checkpoint"

var _ sdk.Msg = &MsgCheckpoint{}

type MsgCheckpoint struct {
	Proposer   common.Address `json:"proposer"`
	StartBlock uint64         `json:"startBlock"`
	EndBlock   uint64         `json:"endBlock"`
	RootHash   common.Hash    `json:"rootHash"`
}

// NewMsgCheckpointBlock creates new checkpoint message using mentioned arguments
func NewMsgCheckpointBlock(proposer common.Address, startBlock uint64, endBlock uint64, roothash common.Hash) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   roothash,
	}
}

// Type returns message type
func (msg MsgCheckpoint) Type() string {
	return Checkpoint
}

func (msg MsgCheckpoint) Route() string { return Checkpoint }

// GetSigners returns address of the signer
func (msg MsgCheckpoint) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Proposer.Bytes())
	return addrs

}

func (msg MsgCheckpoint) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpoint) ValidateBasic() sdk.Error {

	return nil
}

//
// Msg Checkpoint Ack
//

const CheckpointACK = "checkpointACK"

var _ sdk.Msg = &MsgCheckpointAck{}

type MsgCheckpointAck struct {
	HeaderBlock uint64 `json:"headerBlock"`
}

func NewMsgCheckpointAck(headerBlock uint64) MsgCheckpointAck {
	return MsgCheckpointAck{
		HeaderBlock: headerBlock,
	}
}

func (msg MsgCheckpointAck) Type() string {
	return CheckpointACK
}

func (msg MsgCheckpointAck) Route() string { return CheckpointACK }

func (msg MsgCheckpointAck) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgCheckpointAck) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpointAck) ValidateBasic() sdk.Error {

	return nil
}
