package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
)

var cdc = codec.NewLegacyAmino()

var _ sdk.Msg = &MsgEventRecordRequest{}

// NewMsgEventRecord - construct state msg
func NewMsgEventRecord(
	from sdk.AccAddress,
	txHash hmCommon.HeimdallHash,
	logIndex uint64,
	blockNumber uint64,
	id uint64,
	contractAddress sdk.AccAddress,
	data []byte,
	chainID string,

) MsgEventRecordRequest {
	fromStr := strings.ToLower(from.String())
	contractAddressStr := strings.ToLower(contractAddress.String())

	return MsgEventRecordRequest{
		From:            fromStr,
		TxHash:          txHash.String(),
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		Id:              id,
		ContractAddress: contractAddressStr,
		Data:            data,
		ChainId:         chainID,
	}
}

// Route Implements Msg.
func (msg MsgEventRecordRequest) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgEventRecordRequest) Type() string { return "event-record" }

// ValidateBasic Implements Msg.
func (msg MsgEventRecordRequest) ValidateBasic() error {
	if msg.From == "" {
		return sdkerrors.ErrUnknownRequest
	}

	if msg.TxHash == "" {
		return sdkerrors.ErrInvalidAddress
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEventRecordRequest) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgEventRecordRequest) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromHex(msg.From)
	return []sdk.AccAddress{from}
}

// GetTxHash Returns tx hash
func (msg MsgEventRecordRequest) GetTxHash() hmCommon.HeimdallHash {
	return hmCommon.HexToHeimdallHash(msg.TxHash)
}

// GetLogIndex Returns log index
func (msg MsgEventRecordRequest) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgEventRecordRequest) GetSideSignBytes() []byte {
	return nil
}
