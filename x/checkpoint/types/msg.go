package types

import (
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
)

var cdc = codec.NewLegacyAmino()

//
// Checkpoint Msg
//

var _ sdk.Msg = &MsgCheckpoint{}

// NewMsgCheckpointBlock creates new checkpoint message using mentioned arguments
func NewMsgCheckpointBlock(
	proposer sdk.AccAddress,
	startBlock uint64,
	endBlock uint64,
	roothash hmCommonTypes.HeimdallHash,
	accountRootHash hmCommonTypes.HeimdallHash,
	borChainID string,
) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer:        proposer.String(),
		StartBlock:      startBlock,
		EndBlock:        endBlock,
		RootHash:        roothash.Bytes(),
		AccountRootHash: accountRootHash.Bytes(),
		BorChainID:      borChainID,
	}
}

// Type returns message type
func (msg MsgCheckpoint) Type() string {
	return "checkpoint"
}

func (msg MsgCheckpoint) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgCheckpoint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress([]byte(msg.Proposer))}
}

func (msg MsgCheckpoint) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpoint) ValidateBasic() error {
	// if bytes.Equal(msg.RootHash.Bytes(), helper.ZeroHash.Bytes()) {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid rootHash %v", msg.RootHash.String())
	// }

	// if msg.Proposer.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.Proposer.String())
	// }

	if msg.StartBlock >= msg.EndBlock || msg.EndBlock == 0 {
		return hmCommon.ErrInvalidMsg
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgCheckpoint) GetSideSignBytes() []byte {
	// keccak256(abi.encoded(proposer, startBlock, endBlock, rootHash, accountRootHash, bor chain id))
	borChainID, _ := strconv.ParseUint(msg.BorChainID, 10, 64)
	return appendBytes32(
		[]byte(msg.Proposer),
		new(big.Int).SetUint64(msg.StartBlock).Bytes(),
		new(big.Int).SetUint64(msg.EndBlock).Bytes(),
		msg.RootHash,
		msg.AccountRootHash,
		new(big.Int).SetUint64(borChainID).Bytes(),
	)
}

//
// Msg Checkpoint Ack
//

var _ sdk.Msg = &MsgCheckpointAck{}

// NewMsgCheckpointAck Add mainchain commit transaction hash to MsgCheckpointAck
func NewMsgCheckpointAck(
	from sdk.AccAddress,
	number uint64,
	proposer sdk.AccAddress,
	startBlock uint64,
	endBlock uint64,
	rootHash hmCommonTypes.HeimdallHash,
	txHash hmCommonTypes.HeimdallHash,
	logIndex uint64,
) MsgCheckpointAck {
	return MsgCheckpointAck{
		From:       from.String(),
		Number:     number,
		Proposer:   proposer.String(),
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   rootHash.Bytes(),
		TxHash:     txHash.Bytes(),
		LogIndex:   logIndex,
	}
}

func (msg MsgCheckpointAck) Type() string {
	return "checkpoint-ack"
}

func (msg MsgCheckpointAck) Route() string {
	return RouterKey
}

// GetSigners returns signers
func (msg MsgCheckpointAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress([]byte(msg.From))}

}

// GetSignBytes returns sign bytes
func (msg MsgCheckpointAck) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic validate basic
func (msg MsgCheckpointAck) ValidateBasic() error {
	// if msg.From.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid from %v", msg.From.String())
	// }

	// if msg.Proposer.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid empty proposer")
	// }

	// if msg.RootHash.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid empty root hash")
	// }

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgCheckpointAck) GetTxHash() hmCommonTypes.HeimdallHash {
	return hmCommonTypes.BytesToHeimdallHash(msg.TxHash)
}

// GetLogIndex Returns log index
func (msg MsgCheckpointAck) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgCheckpointAck) GetSideSignBytes() []byte {
	return nil
}

//
// Msg Checkpoint No Ack
//

var _ sdk.Msg = &MsgCheckpointNoAck{}

func NewMsgCheckpointNoAck(
	from sdk.AccAddress,
) MsgCheckpointNoAck {
	return MsgCheckpointNoAck{
		From: from.String(),
	}
}

func (msg MsgCheckpointNoAck) Type() string {
	return "checkpoint-no-ack"
}

func (msg MsgCheckpointNoAck) Route() string {
	return RouterKey
}

func (msg MsgCheckpointNoAck) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress([]byte(msg.From))}
}

func (msg MsgCheckpointNoAck) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCheckpointNoAck) ValidateBasic() error {
	if msg.From == "" {
		return hmCommon.ErrInvalidMsg
	}

	return nil
}
