package checkpoint

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

var cdc = codec.New()

//
// Checkpoint Msg
//

var _ sdk.Msg = &MsgCheckpoint{}

// MsgCheckpoint represents checkpoint
type MsgCheckpoint struct {
	Proposer       types.HeimdallAddress `json:"proposer"`
	StartBlock     uint64                `json:"startBlock"`
	EndBlock       uint64                `json:"endBlock"`
	RootHash       types.HeimdallHash    `json:"rootHash"`
	RewardRootHash types.HeimdallHash    `json:"rewardRootHash"`
	TimeStamp      uint64                `json:"timestamp"`
}

// NewMsgCheckpointBlock creates new checkpoint message using mentioned arguments
func NewMsgCheckpointBlock(
	proposer types.HeimdallAddress,
	startBlock uint64,
	endBlock uint64,
	roothash types.HeimdallHash,
	rewardRootHash types.HeimdallHash,
	timestamp uint64,
) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer:       proposer,
		StartBlock:     startBlock,
		EndBlock:       endBlock,
		RootHash:       roothash,
		RewardRootHash: rewardRootHash,
		TimeStamp:      timestamp,
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
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.Proposer)}
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

	if msg.Proposer.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
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

// MsgCheckpointAck Add mainchain commit transaction hash to MsgCheckpointAck
type MsgCheckpointAck struct {
	From        types.HeimdallAddress `json:"from"`
	HeaderBlock uint64                `json:"headerBlock"`
	TxHash      types.HeimdallHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
}

func NewMsgCheckpointAck(from types.HeimdallAddress, headerBlock uint64, txHash types.HeimdallHash, logIndex uint64) MsgCheckpointAck {
	return MsgCheckpointAck{
		From:        from,
		HeaderBlock: headerBlock,
		TxHash:      txHash,
		LogIndex:    logIndex,
	}
}

func (msg MsgCheckpointAck) Type() string {
	return "checkpoint-ack"
}

func (msg MsgCheckpointAck) Route() string {
	return checkpointTypes.RouterKey
}

func (msg MsgCheckpointAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgCheckpointAck) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpointAck) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	}

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
	From      types.HeimdallAddress `json:"from"`
	TimeStamp uint64                `json:"timestamp"`
}

func NewMsgCheckpointNoAck(from types.HeimdallAddress, timestamp uint64) MsgCheckpointNoAck {
	return MsgCheckpointNoAck{
		From:      from,
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
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
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

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	}

	return nil
}
