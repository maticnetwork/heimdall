package clerk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	EventRecords []*types.EventRecord `json:"event_records"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(eventRecords []*types.EventRecord) GenesisState {
	return GenesisState{EventRecords: eventRecords}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make([]*types.EventRecord, 0))
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// add checkpoint headers
	if len(data.EventRecords) != 0 {
		for _, record := range data.EventRecords {
			keeper.SetEventRecord(ctx, *record)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetAllEventRecords(ctx))
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
