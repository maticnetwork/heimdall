package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/common"
)

var cdc = codec.NewLegacyAmino()

var _ sdk.Msg = &MsgProposeSpan{}

// NewMsgProposeSpan creates new propose span message
func NewMsgProposeSpan(
	spanId uint64,
	proposer string,
	startBlock uint64,
	endBlock uint64,
	chainID string,
	seed string,
) MsgProposeSpan {
	return MsgProposeSpan{
		SpanId:     spanId,
		Proposer:   proposer,
		StartBlock: startBlock,
		EndBlock:   endBlock,
		ChainId:    chainID,
		Seed:       seed,
	}
}

func (m MsgProposeSpan) Route() string {
	return RouterKey
}

func (m MsgProposeSpan) Type() string {
	return "propose-span"
}

func (m MsgProposeSpan) ValidateBasic() error {
	if len(m.Proposer) == 0 {
		return common.ErrInvalidMsg
	}
	return nil
}

func (m MsgProposeSpan) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(m)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (m *MsgProposeSpan) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromHex(m.Proposer)
	return []sdk.AccAddress{addr}
}

// GetSideSignBytes returns side sign bytes
func (m MsgProposeSpan) GetSideSignBytes() []byte {
	return nil
}
