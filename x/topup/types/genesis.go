package types

import (
	"encoding/json"
	"errors"

	"github.com/cosmos/cosmos-sdk/x/auth/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, sq := range gs.TopupSequences {
		if sq == "" {
			return errors.New("Invalid Sequence")
		}
	}
	return nil
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(topupSequence []string, dividentAccounts []*hmTypes.DividendAccount) GenesisState {
	return GenesisState{
		TopupSequences:   topupSequence,
		DividendAccounts: dividentAccounts,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil, nil)
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		types.ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, dividendAccounts []*hmTypes.DividendAccount) (map[string]json.RawMessage, error) {
	// set state to staking state
	topupState := GetGenesisStateFromAppState(appState)
	topupState.DividendAccounts = dividendAccounts

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(&topupState)
	return appState, nil
}
