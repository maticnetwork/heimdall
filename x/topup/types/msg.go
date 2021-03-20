package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	common "github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
)

//
// Fee token
//

var _ sdk.Msg = &MsgTopup{}

// NewMsgTopup - construct arbitrary multi-in, multi-out send msg.
func NewMsgTopup(
	fromAddr sdk.AccAddress,
	user sdk.AccAddress,
	fee sdk.Int,
	txhash hmCommon.HeimdallHash,
	logIndex uint64,
	blockNumber uint64,
) MsgTopup {
	return MsgTopup{
		FromAddress: fromAddr.String(),
		User:        user.String(),
		Fee:         &fee,
		TxHash:      txhash.String(),
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
func (msg MsgTopup) ValidateBasic() error {
	if msg.FromAddress == "" {
		return common.ErrEmptyAddr
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTopup) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgTopup) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHex(msg.FromAddress)
	return []sdk.AccAddress{addr}
}

// GetTxHash Returns tx hash
func (msg MsgTopup) GetTxHash() hmCommon.HeimdallHash {
	return hmCommon.HexToHeimdallHash(msg.TxHash)
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

var _ sdk.Msg = &MsgWithdrawFee{}

// NewMsgWithdrawFee - construct arbitrary fee withdraw msg
func NewMsgWithdrawFee(
	fromAddr sdk.AccAddress,
	amount sdk.Int,
) MsgWithdrawFee {
	return MsgWithdrawFee{
		UserAddress: fromAddr.String(),
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
func (msg MsgWithdrawFee) ValidateBasic() error {
	if msg.UserAddress == "" {
		return common.ErrEmptyAddr
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawFee) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawFee) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHex(msg.UserAddress)
	return []sdk.AccAddress{addr}
}
