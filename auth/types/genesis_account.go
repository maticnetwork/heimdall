package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"

	"github.com/maticnetwork/heimdall/types"
)

//-----------------------------------------------------------------------------
// GenesisAccount

var _ Account = (*GenesisAccount)(nil)

// GenesisAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type GenesisAccount struct {
	Address       types.HeimdallAddress `json:"address" yaml:"address"`
	Coins         types.Coins           `json:"coins" yaml:"coins"`
	AccountNumber uint64                `json:"account_number" yaml:"account_number"`
	Sequence      uint64                `json:"sequence" yaml:"sequence"`
}

// NewGenesisAccount creates new genesis account
func NewGenesisAccount(acc Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	return gacc
}

// BaseToGenesisAccount converts base account to genesis account
func BaseToGenesisAccount(acc BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		Sequence:      acc.Sequence,
		AccountNumber: acc.AccountNumber,
	}
}

// String implements fmt.Stringer
func (acc GenesisAccount) String() string {
	return fmt.Sprintf(`Account:
  Address:       %s
  Coins:         %s
  AccountNumber: %d
  Sequence:      %d`,
		acc.Address, acc.Coins, acc.AccountNumber, acc.Sequence,
	)
}

// GetAddress - Implements sdk.Account.
func (acc GenesisAccount) GetAddress() types.HeimdallAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *GenesisAccount) SetAddress(addr types.HeimdallAddress) error {
	if len(acc.Address) != 0 && !acc.Address.Empty() {
		return errors.New("cannot override GenesisAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc GenesisAccount) GetPubKey() crypto.PubKey {
	return nil
}

// SetPubKey - Implements sdk.Account.
func (acc *GenesisAccount) SetPubKey(pubKey crypto.PubKey) error {
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc GenesisAccount) GetCoins() types.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *GenesisAccount) SetCoins(coins types.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc GenesisAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *GenesisAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc GenesisAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *GenesisAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *GenesisAccount) SpendableCoins(_ time.Time) types.Coins {
	return acc.GetCoins()
}

// Validate checks for errors on the account fields
func (acc GenesisAccount) Validate() error {
	return nil
}

// MarshalYAML returns the YAML representation of an account.
func (acc GenesisAccount) MarshalYAML() (interface{}, error) {
	var bs []byte
	var err error

	bs, err = yaml.Marshal(struct {
		Address       types.HeimdallAddress
		Coins         types.Coins
		AccountNumber uint64
		Sequence      uint64
	}{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	})
	if err != nil {
		return nil, err
	}

	return string(bs), err
}
