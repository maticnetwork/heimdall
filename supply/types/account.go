package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	authExported "github.com/maticnetwork/heimdall/auth/exported"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/supply/exported"
	"github.com/maticnetwork/heimdall/types"
)

//
// Module account
//

var _ exported.ModuleAccountI = (*ModuleAccount)(nil)
var _ authTypes.Account = (*ModuleAccount)(nil)

type (
	// ModuleAccountInterface exported module account interface
	ModuleAccountInterface = exported.ModuleAccountI
)

// ModuleAccount defines an account for modules that holds coins on a pool
type ModuleAccount struct {
	*authTypes.BaseAccount

	Name        string   `json:"name" yaml:"name"`               // name of the module
	Permissions []string `json:"permissions" yaml:"permissions"` // permissions of module account
}

// NewModuleAddress creates an AccAddress from the hash of the module's name
func NewModuleAddress(name string) types.HeimdallAddress {
	return types.BytesToHeimdallAddress(crypto.AddressHash([]byte(name)).Bytes())
}

// NewEmptyModuleAccount creates empty module account
func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := authTypes.NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: &baseAcc,
		Name:        name,
		Permissions: permissions,
	}
}

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba *authTypes.BaseAccount, name string, permissions ...string) *ModuleAccount {
	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: ba,
		Name:        name,
		Permissions: permissions,
	}
}

// AddPermissions adds the permissions to the module account's list of granted
// permissions.
func (ma *ModuleAccount) AddPermissions(permissions ...string) {
	ma.Permissions = append(ma.Permissions, permissions...)
}

// RemovePermission removes the permission from the list of granted permissions
// or returns an error if the permission is has not been granted.
func (ma *ModuleAccount) RemovePermission(permission string) error {
	for i, perm := range ma.Permissions {
		if perm == permission {
			ma.Permissions = append(ma.Permissions[:i], ma.Permissions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cannot remove non granted permission %s", permission)
}

// HasPermission returns whether or not the module account has permission.
func (ma ModuleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma ModuleAccount) GetName() string {
	return ma.Name
}

// GetPermissions returns permissions granted to the module account
func (ma ModuleAccount) GetPermissions() []string {
	return ma.Permissions
}

// SetPubKey - Implements Account
func (ma ModuleAccount) SetPubKey(pubKey crypto.PubKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// SetSequence - Implements Account
func (ma ModuleAccount) SetSequence(seq uint64) error {
	return fmt.Errorf("not supported for module accounts")
}

// String follows stringer interface
func (ma ModuleAccount) String() string {
	b, err := yaml.Marshal(ma)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (ma ModuleAccount) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		Address       types.HeimdallAddress
		Coins         sdk.Coins
		PubKey        string
		AccountNumber uint64
		Sequence      uint64
		Name          string
		Permissions   []string
	}{
		Address:       ma.Address,
		Coins:         ma.Coins,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Sequence:      ma.Sequence,
		Name:          ma.Name,
		Permissions:   ma.Permissions,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

//
// Account processor
//

// AccountProcessor process supply's module account
func AccountProcessor(ga *authTypes.GenesisAccount, ba *authTypes.BaseAccount) authExported.Account {
	// module accounts
	if ga.ModuleName != "" {
		return NewModuleAccount(ba, ga.ModuleName, ga.ModulePermissions...)
	}

	return ba
}
