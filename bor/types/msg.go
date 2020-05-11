package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/common"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Propose Span Msg
//

var _ sdk.Msg = &MsgProposeSpan{}

// MsgProposeSpan creates msg propose span
type MsgProposeSpan struct {
	ID         uint64                  `json:"span_id"`
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	ChainID    string                  `json:"bor_chain_id"`
	Seed       common.Hash             `json:"seed"`
}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpan(
	id uint64,
	proposer hmTypes.HeimdallAddress,
	startBlock uint64,
	endBlock uint64,
	chainID string,
	seed common.Hash,
) MsgProposeSpan {
	return MsgProposeSpan{
		ID:         id,
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		ChainID:    chainID,
		Seed:       seed,
	}
}

// Type returns message type
func (msg MsgProposeSpan) Type() string {
	return "propose-span"
}

// Route returns route for message
func (msg MsgProposeSpan) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgProposeSpan) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeSpan) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeSpan) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}

	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgProposeSpan) GetSideSignBytes() []byte {
	return nil
}
