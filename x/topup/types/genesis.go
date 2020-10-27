package types

import (
	"encoding/json"
	"errors"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/auth/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	TopupSequences   []string                  `json:"tx_sequences" yaml:"tx_sequences"`
	DividentAccounts []hmTypes.DividendAccount `json:"dividend_accounts" yaml:"dividend_accounts"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(topupSequence []string, dividentAccounts []hmTypes.DividendAccount) GenesisState {
	return GenesisState{
		TopupSequences:   topupSequence,
		DividentAccounts: dividentAccounts,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() GenesisState {
	return NewGenesisState(nil, nil)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate(data GenesisState) error {
	for _, sq := range data.TopupSequences {
		if sq == "" {
			return errors.New("Invalid Sequence")
		}
	}
	return nil
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
func SetGenesisStateToAppState(appState map[string]json.RawMessage, dividendAccounts []hmTypes.DividendAccount) (map[string]json.RawMessage, error) {
	// set state to staking state
	topupState := GetGenesisStateFromAppState(appState)
	topupState.DividentAccounts = dividendAccounts

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(topupState)
	return appState, nil
}
