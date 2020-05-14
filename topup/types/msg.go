package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
)

//
// Fee token
//

// MsgTopup - high level transaction of the fee coin module
type MsgTopup struct {
	FromAddress types.HeimdallAddress `json:"from_address"`
	User        types.HeimdallAddress `json:"user"`
	Fee         sdk.Int               `json:"fee"`
	TxHash      types.HeimdallHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
	BlockNumber uint64                `json:"block_number"`
}

var _ sdk.Msg = MsgTopup{}

// NewMsgTopup - construct arbitrary multi-in, multi-out send msg.
func NewMsgTopup(
	fromAddr types.HeimdallAddress,
	user types.HeimdallAddress,
	fee sdk.Int,
	txhash types.HeimdallHash,
	logIndex uint64,
	blockNumber uint64,
) MsgTopup {
	return MsgTopup{
		FromAddress: fromAddr,
		User:        user,
		Fee:         fee,
		TxHash:      txhash,
		LogIndex:    logIndex,
		BlockNumber: blockNumber,
	}
}

// Route Implements Msg.
func (msg MsgTopup) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgTopup) Type() string {
	return ModuleName
}

// ValidateBasic Implements Msg.
func (msg MsgTopup) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.FromAddress.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.FromAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTopup) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgTopup) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.FromAddress)}
}

// GetTxHash Returns tx hash
func (msg MsgTopup) GetTxHash() types.HeimdallHash {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgTopup) GetLogIndex() uint64 {
	return msg.LogIndex
}

func (msg MsgTopup) GetSideSignBytes() []byte {
	return nil
}

//
// Fee token withdrawal
//

// MsgWithdrawFee - high level transaction of the fee coin withdrawal module
type MsgWithdrawFee struct {
	UserAddress types.HeimdallAddress `json:"from_address"`
	Amount      sdk.Int               `json:"amount"`
}

var _ sdk.Msg = MsgWithdrawFee{}

// NewMsgWithdrawFee - construct arbitrary fee withdraw msg
func NewMsgWithdrawFee(
	fromAddr types.HeimdallAddress,
	amount sdk.Int,
) MsgWithdrawFee {
	return MsgWithdrawFee{
		UserAddress: fromAddr,
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgWithdrawFee) Route() string {
	return RouterKey
}

// Type Implements Msg.
func (msg MsgWithdrawFee) Type() string {
	return "withdraw"
}

// ValidateBasic Implements Msg.
func (msg MsgWithdrawFee) ValidateBasic() sdk.Error {
	if msg.UserAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.UserAddress.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.UserAddress.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawFee) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.UserAddress)}
}
