package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
)

var cdc = wire.NewCodec()

const MsgType = "checkpoint"

var _ sdk.Msg = &MsgCheckpoint{}

type MsgCheckpoint struct {
	// TODO variable as we dont know who will call this
	Proposer   common.Address `json:"address"` // address of the validator owner
	StartBlock uint64         `json:"startBlock"`
	EndBlock   uint64         `json:"endBlock"`
	RootHash   common.Hash    `json:"rootHash"`
}

//

func NewMsgCheckpointBlock(startBlock uint64, endBlock uint64, roothash common.Hash, proposer string) MsgCheckpoint {
	return MsgCheckpoint{
		// TODO remove after testing
		Proposer:   common.HexToAddress(proposer),
		StartBlock: startBlock,
		EndBlock:   endBlock,
		RootHash:   roothash,
	}
}

//nolint
func (msg MsgCheckpoint) Type() string { return MsgType }
func (msg MsgCheckpoint) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Proposer.Bytes())
	return addrs

}

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
	// TODO add checks
	return nil
}

var _ sdk.Tx = BaseTx{}

// Basetx
type BaseTx struct {
	Msg MsgCheckpoint
}

func NewBaseTx(msg MsgCheckpoint) BaseTx {
	return BaseTx{
		Msg: msg,
	}
}

func (tx BaseTx) GetMsgs() []sdk.Msg { return []sdk.Msg{tx.Msg} }

//
//func (app *HeimdallApp) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
//	var tx = checkpoint.BaseTx{}
//
//	err := rlp.DecodeBytes(txBytes, &tx)
//	if err != nil {
//		return nil, sdk.ErrTxDecode(err.Error())
//	}
//	return tx, nil
//}
