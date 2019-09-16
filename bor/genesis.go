package bor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/types"
)

// GenesisState is the bor state that must be provided at genesis.
type GenesisState struct {
	SprintDuration uint64        `json:"sprint_duration" yaml:"sprint_duration"` // sprint duration
	SpanDuration   uint64        `json:"span_duration" yaml:"span_duration"`     // span duration ie number of blocks for which val set is frozen on heimdall
	ProducerCount  uint64        `json:"producer_count" yaml:"producer_count"`   // producer count per span
	Spans          []*types.Span `json:"spans" yaml:"spans"`                     // list of spans
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(sprintDuration uint64, spanDuration uint64, producerCount uint64, spans []*types.Span) GenesisState {
	return GenesisState{
		SprintDuration: sprintDuration,
		SpanDuration:   spanDuration,
		ProducerCount:  producerCount,
		Spans:          spans,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultSprintDuration, DefaultSpanDuration, DefaultProducerCount, make([]*types.Span, 0))
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetSprintDuration(ctx, data.SprintDuration)
	keeper.SetSpanDuration(ctx, data.SpanDuration)
	if len(data.Spans) > 0 {
		// add new span
		for _, span := range data.Spans {
			keeper.AddNewRawSpan(ctx, *span)
		}

		// update last span
		keeper.UpdateLastSpan(ctx, data.Spans[len(data.Spans)-1].StartBlock)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(
		keeper.GetSprintDuration(ctx),
		keeper.GetSpanDuration(ctx),
		keeper.GetProducerCount(ctx),
		// TODO think better way to export all spans
		keeper.GetAllSpans(ctx),
	)
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }
