package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

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

//
// Propose Span v2 Msg
//

var _ sdk.Msg = &MsgProposeSpanV2{}

// MsgProposeSpanV2 creates msg propose span
type MsgProposeSpanV2 struct {
	ID         uint64                  `json:"span_id"`
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	ChainID    string                  `json:"bor_chain_id"`
	Seed       common.Hash             `json:"seed"`
	SeedAuthor common.Address          `json:"seed_author"`
}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpanV2(
	id uint64,
	proposer hmTypes.HeimdallAddress,
	startBlock uint64,
	endBlock uint64,
	chainID string,
	seed common.Hash,
	seedAuthor common.Address,
) MsgProposeSpanV2 {
	return MsgProposeSpanV2{
		ID:         id,
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		ChainID:    chainID,
		Seed:       seed,
		SeedAuthor: seedAuthor,
	}
}

// Type returns message type
func (msg MsgProposeSpanV2) Type() string {
	return "propose-span-v2"
}

// Route returns route for message
func (msg MsgProposeSpanV2) Route() string {
	return RouterKey
}

// GetSigners returns address of the signer
func (msg MsgProposeSpanV2) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.Proposer)}
}

// GetSignBytes returns sign bytes for proposeSpan message type
func (msg MsgProposeSpanV2) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic validates the message and returns error
func (msg MsgProposeSpanV2) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if msg.Seed.Cmp(common.Hash{}) == 0 {
		return sdk.ErrUnknownRequest("Seed cannot be empty")
	}
	if msg.SeedAuthor.Cmp(common.Address{}) == 0 {
		return sdk.ErrUnknownRequest("SeedAuthor cannot be empty")
	}
	return nil
}

// GetSideSignBytes returns side sign bytes
func (msg MsgProposeSpanV2) GetSideSignBytes() []byte {
	return nil
}

type MsgBackfillSpans struct {
	Proposer        hmTypes.HeimdallAddress `json:"proposer"`
	ChainID         string                  `json:"chain_id"`
	LatestSpanID    uint64                  `json:"latest_span_id"`
	LatestBorSpanID uint64                  `json:"latest_bor_span_id"`
}

func NewMsgBackfillSpans(proposer hmTypes.HeimdallAddress, chainID string, latestSpanID, latestBorSpanID uint64) MsgBackfillSpans {
	return MsgBackfillSpans{
		Proposer:        proposer,
		ChainID:         chainID,
		LatestSpanID:    latestSpanID,
		LatestBorSpanID: latestBorSpanID,
	}
}

func (msg MsgBackfillSpans) Type() string {
	return "backfill-spans"
}

func (msg MsgBackfillSpans) Route() string {
	return RouterKey
}

func (msg MsgBackfillSpans) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(hmTypes.HeimdallAddress(msg.Proposer))}
}

func (msg MsgBackfillSpans) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

func (msg MsgBackfillSpans) ValidateBasic() sdk.Error {
	if msg.Proposer.Empty() {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if msg.ChainID == "" {
		return sdk.ErrUnknownRequest("ChainID cannot be empty")
	}
	if msg.LatestSpanID == 0 {
		return sdk.ErrUnknownRequest("LatestSpanID cannot be zero")
	}
	if msg.LatestBorSpanID == 0 {
		return sdk.ErrUnknownRequest("LatestBorSpanID cannot be zero")
	}
	return nil
}

func (msg MsgBackfillSpans) GetSideSignBytes() []byte {
	return nil
}
