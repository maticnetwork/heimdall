package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	yaml "gopkg.in/yaml.v2"

	"github.com/maticnetwork/heimdall/auth/exported"
	"github.com/maticnetwork/heimdall/types"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.PrivKeyAminoName, nil)
}

// Account is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements Account.
type (
	Account = exported.Account
)

//-----------------------------------------------------------------------------
// BaseAccount

var _ Account = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       types.HeimdallAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins             `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64                `json:"account_number" yaml:"account_number"`
	Sequence      uint64                `json:"sequence" yaml:"sequence"`
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(
	address types.HeimdallAddress,
	coins sdk.Coins,
	pubKey crypto.PubKey,
	accountNumber uint64,
	sequence uint64,
) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		// pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
		var pubObject secp256k1.PubKeySecp256k1
		cdc.MustUnmarshalBinaryBare(acc.PubKey.Bytes(), &pubObject)
		pubkey = "0x" + hex.EncodeToString(pubObject[:])
	}

	return fmt.Sprintf(`Account:
  Address:       %s
  Pubkey:        %s
  Coins:         %s
  AccountNumber: %d
  Sequence:      %d`,
		acc.Address, pubkey, acc.Coins, acc.AccountNumber, acc.Sequence,
	)
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() Account {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr types.HeimdallAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() types.HeimdallAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr types.HeimdallAddress) error {
	if len(acc.Address) != 0 && !acc.Address.Empty() {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && !acc.Address.Empty() &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	var bs []byte
	var err error
	var pubkey string

	if acc.PubKey != nil {
		pubkey, err = sdk.Bech32ifyAccPub(acc.PubKey)
		if err != nil {
			return nil, err
		}
	}

	bs, err = yaml.Marshal(struct {
		Address       types.HeimdallAddress
		Coins         sdk.Coins
		PubKey        string
		AccountNumber uint64
		Sequence      uint64
	}{
		Address:       acc.Address,
		Coins:         acc.Coins,
		PubKey:        pubkey,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	})
	if err != nil {
		return nil, err
	}

	return string(bs), err
}

//
// light base account
//

// LightBaseAccount - a base account structure.
type LightBaseAccount struct {
	Address       types.HeimdallAddress `json:"address" yaml:"address"`
	AccountNumber uint64                `json:"account_number" yaml:"account_number"`
	Sequence      uint64                `json:"sequence" yaml:"sequence"`
}
