package types

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/maticnetwork/heimdall/types"
)

//
// Genesis accounts
//

// GenesisAccounts defines a slice of GenesisAccount objects
type GenesisAccounts []GenesisAccount

// Contains returns true if the given address exists in a slice of GenesisAccount
// objects.
func (accounts GenesisAccounts) Contains(addr types.HeimdallAddress) bool {
	for _, acc := range accounts {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
}

//
// Gensis state
//

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	Params   Params          `json:"params" yaml:"params"`
	Accounts GenesisAccounts `json:"accounts" yaml:"accounts"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, accounts GenesisAccounts) GenesisState {
	return GenesisState{
		Params:   params,
		Accounts: SanitizeGenesisAccounts(accounts),
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), GenesisAccounts{})
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := ModuleCdc.UnmarshalJSON(appState[ModuleName], &genesisState)
		if err != nil {
			panic(err)
		}
	}

	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, accounts []GenesisAccount) (map[string]json.RawMessage, error) {
	authState := GetGenesisStateFromAppState(appState)
	authState.Accounts = accounts
	var err error
	appState[ModuleName], err = ModuleCdc.MarshalJSON(authState)
	if err != nil {
		return appState, err
	}
	return appState, nil
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return ValidateGenAccounts(data.Accounts)
}

// SanitizeGenesisAccounts sorts accounts and coin sets.
func SanitizeGenesisAccounts(genAccs GenesisAccounts) GenesisAccounts {
	sort.Slice(genAccs, func(i, j int) bool {
		return genAccs[i].GetAccountNumber() < genAccs[j].GetAccountNumber()
	})

	for _, acc := range genAccs {
		if err := acc.SetCoins(acc.GetCoins().Sort()); err != nil {
			panic(err)
		}
	}

	return genAccs
}

// ValidateGenAccounts validates an array of GenesisAccounts and checks for duplicates
func ValidateGenAccounts(accounts GenesisAccounts) error {
	addrMap := make(map[string]bool, len(accounts))
	for _, acc := range accounts {

		// check for duplicated accounts
		addrStr := acc.GetAddress().String()
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		addrMap[addrStr] = true

		// check account specific validation
		if err := acc.Validate(); err != nil {
			return fmt.Errorf("invalid account found in genesis state; address: %s, error: %s", addrStr, err.Error())
		}
	}
	return nil
}

// GenesisAccountIterator implements genesis account iteration.
type GenesisAccountIterator struct{}

// IterateGenesisAccounts iterates over all the genesis accounts found in
// appGenesis and invokes a callback on each genesis account. If any call
// returns true, iteration stops.
func (GenesisAccountIterator) IterateGenesisAccounts(appGenesis map[string]json.RawMessage, cb func(Account) (stop bool)) {
	for _, genAcc := range GetGenesisStateFromAppState(appGenesis).Accounts {
		if cb(&genAcc) {
			break
		}
	}
}
