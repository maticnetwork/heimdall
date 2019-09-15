package clerk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	StateRecords []*types.Record `json:"state_records"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(stateRecords []*types.Record) GenesisState {
	return GenesisState{StateRecords: stateRecords}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make([]*types.Record, 0))
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// add checkpoint headers
	if len(data.StateRecords) != 0 {
		for _, record := range data.StateRecords {
			keeper.SetStateRecord(ctx, *record)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetAllStateRecords(ctx))
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
