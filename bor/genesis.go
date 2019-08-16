package bor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState is the bor state that must be provided at genesis.
type GenesisState struct {
	SprintDuration uint64 `json:"sprintDuration" yaml:"sprintDuration"` // sprint duration
	SpanDuration   uint64 `json:"spanDuration" yaml:"spanDuration"`     // span duration ie number of blocks for which val set is frozen on heimdall
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(sprintDuration uint64, spanDuration uint64) GenesisState {
	return GenesisState{
		SprintDuration: sprintDuration,
		SpanDuration:   spanDuration,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultSprintDuration, DefaultSpanDuration)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetSprintDuration(ctx, data.SprintDuration)
	keeper.SetSpanDuration(ctx, data.SpanDuration)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(
		keeper.GetSprintDuration(ctx),
		keeper.GetSpanDuration(ctx),
	)
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }
