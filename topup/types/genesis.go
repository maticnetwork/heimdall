package types

import (
	"encoding/json"
	"errors"

	"github.com/maticnetwork/heimdall/bor/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

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

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil, nil)
}

// ValidateGenesis performs basic validation of topup genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
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
