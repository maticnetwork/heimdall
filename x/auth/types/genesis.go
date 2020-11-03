package types

import (
	"encoding/json"
	"sort"

	"github.com/maticnetwork/heimdall/x/blog/types"
)

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, accounts []*GenesisAccount) GenesisState {
	return GenesisState{
		Params:   params,
		Accounts: SanitizeGenesisAccounts(accounts),
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []*GenesisAccount{})
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
func SetGenesisStateToAppState(appState map[string]json.RawMessage, accounts []*GenesisAccount) (map[string]json.RawMessage, error) {
	authState := GetGenesisStateFromAppState(appState)
	authState.Accounts = accounts

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(&authState)
	return appState, nil
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	// if err := data.Params.Validate(); err != nil {
	// 	return err
	// }

	// return ValidateGenAccounts(data.Accounts)
	return nil
}

// SanitizeGenesisAccounts sorts accounts and coin sets.
func SanitizeGenesisAccounts(genAccs []*GenesisAccount) []*GenesisAccount {
	sort.Slice(genAccs, func(i, j int) bool {
		return genAccs[i].AccountNumber < genAccs[j].AccountNumber
	})

	for _, acc := range genAccs {
		acc.Coins = acc.Coins.Sort()
	}

	return genAccs
}

// ValidateGenAccounts validates an array of GenesisAccounts and checks for duplicates
func ValidateGenAccounts(accounts []GenesisAccount) error {
	// addrMap := make(map[string]bool, len(accounts))
	// for _, acc := range accounts {

	// 	// check for duplicated accounts
	// 	addrStr := acc.Address.String()
	// 	if _, ok := addrMap[addrStr]; ok {
	// 		return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
	// 	}

	// 	addrMap[addrStr] = true

	// 	// check account specific validation
	// 	if err := acc.Validate(); err != nil {
	// 		return fmt.Errorf("invalid account found in genesis state; address: %s, error: %s", addrStr, err.Error())
	// 	}
	// }
	return nil
}

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
