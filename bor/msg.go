package bor

import (
	"bytes"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
)

var cdc = codec.New()

// CheckpointRoute represents rount in app
const BorProposeSpanRoute = "proposeSpan"

//
// Propose Span Msg
//

var _ sdk.Msg = &MsgProposeSpan{}

type MsgProposeSpan struct {
	Proposer   common.Address `json:"proposer"`
	StartBlock uint64         `json:"startBlock"`
	EndBlock   uint64         `json:"endBlock"`
	RootHash   common.Hash    `json:"rootHash"`
	TimeStamp  uint64         `json:"timestamp"`
}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpan(proposer common.Address, startBlock uint64, endBlock uint64, roothash common.Hash, timestamp uint64) MsgProposeSpan {
	return MsgProposeSpan{
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   roothash,
		TimeStamp:  timestamp,
	}
}

// Type returns message type
func (msg MsgProposeSpan) Type() string {
	return "checkpoint"
}

func (msg MsgProposeSpan) Route() string {
	return CheckpointRoute
}

// GetSigners returns address of the signer
func (msg MsgProposeSpan) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Proposer.Bytes())
	return addrs
}

func (msg MsgProposeSpan) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgProposeSpan) ValidateBasic() sdk.Error {
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
