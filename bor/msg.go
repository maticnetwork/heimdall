package bor

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/common"
	"time"
)

var cdc = codec.New()

// BorRoute represents route in app
const BorRoute = "bor"

//
// Propose Span Msg
//

var _ sdk.Msg = &MsgProposeSpan{}

type MsgProposeSpan struct {
	StartBlock uint64 `json:"startBlock"`

	// Timestamp only exits to allow submission of multiple transactions without bringing in nonce
	TimeStamp uint64 `json:"timestamp"`
}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpan(startBlock uint64, timestamp uint64) MsgProposeSpan {
	return MsgProposeSpan{
		StartBlock: startBlock,
		TimeStamp:  timestamp,
	}
}

// Type returns message type
func (msg MsgProposeSpan) Type() string {
	return "ProposeSpan"
}

// Route returns route for message
func (msg MsgProposeSpan) Route() string {
	return BorRoute
}

// GetSigners returns address of the signer
func (msg MsgProposeSpan) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	return addrs
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeSpan) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeSpan) ValidateBasic() sdk.Error {
	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		return common.ErrInvalidMsg(common.DefaultCodespace, "Invalid timestamp %d", msg.TimeStamp)
	}
	return nil
}
