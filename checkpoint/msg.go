package checkpoint

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

var cdc = codec.New()

//
// Checkpoint Msg
//

var _ sdk.Msg = &MsgCheckpoint{}

type MsgCheckpoint struct {
	Proposer   common.Address `json:"proposer"`
	StartBlock uint64         `json:"startBlock"`
	EndBlock   uint64         `json:"endBlock"`
	RootHash   common.Hash    `json:"rootHash"`
	TimeStamp  uint64         `json:"timestamp"`
}

// NewMsgCheckpointBlock creates new checkpoint message using mentioned arguments
func NewMsgCheckpointBlock(proposer common.Address, startBlock uint64, endBlock uint64, roothash common.Hash, timestamp uint64) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   roothash,
		TimeStamp:  timestamp,
	}
}

// Type returns message type
func (msg MsgCheckpoint) Type() string {
	return "checkpoint"
}

func (msg MsgCheckpoint) Route() string {
	return checkpointTypes.RouterKey
}

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
	if bytes.Equal(msg.RootHash.Bytes(), helper.ZeroHash.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid rootHash %v", msg.RootHash.String())
	}

	if bytes.Equal(msg.Proposer.Bytes(), helper.ZeroAddress.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
	}

	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid timestamp %d", msg.TimeStamp)
	}

	if msg.StartBlock >= msg.EndBlock || msg.EndBlock == 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid startBlock %v or/and endBlock %v", msg.StartBlock, msg.EndBlock)
	}

	return nil
}

//
// Msg Checkpoint Ack
//

var _ sdk.Msg = &MsgCheckpointAck{}

type MsgCheckpointAck struct {
	HeaderBlock uint64 `json:"headerBlock"`
	Timestamp   uint64 `json:"timestamp"`
}

func NewMsgCheckpointAck(headerBlock uint64, timestamp uint64) MsgCheckpointAck {
	return MsgCheckpointAck{
		HeaderBlock: headerBlock,
		Timestamp:   timestamp,
	}
}

func (msg MsgCheckpointAck) Type() string {
	return "checkpoint-ack"
}

func (msg MsgCheckpointAck) Route() string {
	return checkpointTypes.RouterKey
}

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
	childBlockInterval := helper.GetConfig().ChildBlockInterval
	if msg.HeaderBlock > 0 && msg.HeaderBlock%childBlockInterval != 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid header block %d", msg.HeaderBlock)
	}

	return nil
}

//
// Msg Checkpoint No Ack
//

var _ sdk.Msg = &MsgCheckpointNoAck{}

type MsgCheckpointNoAck struct {
	TimeStamp uint64 `json:"timestamp"`
}

func NewMsgCheckpointNoAck(timestamp uint64) MsgCheckpointNoAck {
	return MsgCheckpointNoAck{
		TimeStamp: timestamp,
	}
}

func (msg MsgCheckpointNoAck) Type() string {
	return "checkpoint-no-ack"
}

func (msg MsgCheckpointNoAck) Route() string {
	return checkpointTypes.RouterKey
}

func (msg MsgCheckpointNoAck) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0)
	return addrs
}

func (msg MsgCheckpointNoAck) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpointNoAck) ValidateBasic() sdk.Error {
	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid timestamp %d", msg.TimeStamp)
	}

	return nil
}
