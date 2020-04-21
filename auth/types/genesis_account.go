package types

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyExported "github.com/maticnetwork/heimdall/supply/exported"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisAccount is a struct for account initialization used exclusively during genesis
type GenesisAccount struct {
	Address       hmTypes.HeimdallAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins               `json:"coins" yaml:"coins"`
	Sequence      uint64                  `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber uint64                  `json:"account_number" yaml:"account_number"`

	// module account fields
	ModuleName        string   `json:"module_name" yaml:"module_name"`               // name of the module account
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

// Validate checks for errors on the vesting and module account parameters
func (ga GenesisAccount) Validate() error {
	// don't allow blank (i.e just whitespaces) on the module name
	if ga.ModuleName != "" && strings.TrimSpace(ga.ModuleName) == "" {
		return errors.New("module account name cannot be blank")
	}

	return nil
}

// NewGenesisAccountRaw creates a new GenesisAccount object
func NewGenesisAccountRaw(
	address hmTypes.HeimdallAddress,
	coins sdk.Coins,
	module string,
	permissions ...string,
) GenesisAccount {

	return GenesisAccount{
		Address:           address,
		Coins:             coins,
		Sequence:          0,
		AccountNumber:     0, // ignored set by the account keeper during InitGenesis
		ModuleName:        module,
		ModulePermissions: permissions,
	}
}

// NewGenesisAccount creates a GenesisAccount instance from a BaseAccount.
func NewGenesisAccount(acc *BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}
}

// NewGenesisAccountI creates a GenesisAccount instance from an Account interface.
func NewGenesisAccountI(acc Account) (GenesisAccount, error) {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	if err := gacc.Validate(); err != nil {
		return gacc, err
	}

	switch acc := acc.(type) {
	case supplyExported.ModuleAccountI:
		gacc.ModuleName = acc.GetName()
		gacc.ModulePermissions = acc.GetPermissions()
	}

	return gacc, nil
}

// ToAccount converts a GenesisAccount to an Account interface
func (ga *GenesisAccount) ToAccount() Account {
	bacc := NewBaseAccount(ga.Address, ga.Coins.Sort(), nil, ga.AccountNumber, ga.Sequence)
	return bacc
}

// ------------------------------------------
//

// GenesisAccounts list of genesis account
type GenesisAccounts []GenesisAccount

// Contains checks if genesis accounts contain an address
func (gaccs GenesisAccounts) Contains(acc hmTypes.HeimdallAddress) bool {
	for _, gacc := range gaccs {
		if gacc.Address.Equals(acc) {
			return true
		}
	}
	return false
}
