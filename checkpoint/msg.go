package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
)

var cdc = wire.NewCodec()

// MsgType represents string for message type
const MsgType = "checkpoint"

var _ sdk.Msg = &MsgCheckpoint{}

// MsgCheckpoint represents incoming checkpoint format
type MsgCheckpoint struct {
	Proposer   common.Address `json:"address"` // address of the validator owner
	StartBlock uint64         `json:"start_block"`
	EndBlock   uint64         `json:"end_block"`
	RootHash   common.Hash    `json:"root_hash"`
}

// NewMsgCheckpointBlock creates new checkpoint message using mentioned arguments
func NewMsgCheckpointBlock(startBlock uint64, endBlock uint64, roothash common.Hash, proposer string) MsgCheckpoint {
	return MsgCheckpoint{
		Proposer:   common.HexToAddress(proposer),
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   roothash,
	}
}

// Type returns message type
func (msg MsgCheckpoint) Type() string {
	return MsgType
}

// GetSigners returns address of the signer
func (msg MsgCheckpoint) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Proposer.Bytes())
	return addrs

}

// GetSignBytes returns the bytes for the message signer to sign on
func (msg MsgCheckpoint) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic checks quick validation
func (msg MsgCheckpoint) ValidateBasic() sdk.Error {
	// TODO add checks
	return nil
}

// assertion
var _ sdk.Tx = BaseTx{}

// BaseTx represents base tx tendermint needs
type BaseTx struct {
	Msg MsgCheckpoint
}

// NewBaseTx drafts BaseTx with messages
func NewBaseTx(msg MsgCheckpoint) BaseTx {
	return BaseTx{
		Msg: msg,
	}
}

// GetMsgs returns array of messages
func (tx BaseTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.Msg}
}
