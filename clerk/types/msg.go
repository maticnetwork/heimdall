package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/types"
)

// MsgEventRecord - state msg
type MsgEventRecord struct {
	From     types.HeimdallAddress `json:"from"`
	TxHash   types.HeimdallHash    `json:"tx_hash"`
	ID       uint64                `json:"id"`
	Contract types.HeimdallAddress `json:"contract"`
	Data     []byte                `json:"data"`
}

var _ sdk.Msg = MsgEventRecord{}

// NewMsgEventRecord - construct state msg
func NewMsgEventRecord(
	from types.HeimdallAddress,
	txHash types.HeimdallHash,
	id uint64,
	contract types.HeimdallAddress,
	data []byte,
) MsgEventRecord {
	return MsgEventRecord{From: from, TxHash: txHash, ID: id, Contract: contract, Data: data}
}

// Route Implements Msg.
func (msg MsgEventRecord) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgEventRecord) Type() string { return "event-record" }

// ValidateBasic Implements Msg.
func (msg MsgEventRecord) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.TxHash.Empty() {
		return sdk.ErrInvalidAddress("missing tx hash")
	}

	if msg.Contract.Empty() {
		return sdk.ErrInvalidAddress("missing contract address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEventRecord) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgEventRecord) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}
