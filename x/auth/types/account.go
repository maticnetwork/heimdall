package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/x/auth/exported"
)

func init() {
	ModuleCdc.RegisterConcrete(secp256k1.PubKey{}, secp256k1.PubKeyName, nil)
	ModuleCdc.RegisterConcrete(secp256k1.PrivKey{}, secp256k1.PrivKeyName, nil)
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

// var _ Account = (*BaseAccount)(nil)

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(
	address string,
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

	// if acc.PubKey != nil {
	// 	// pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	// 	var pubObject secp256k1.PubKey
	// 	ModuleCdc.MustUnmarshalBinaryBare(acc.PubKey.Bytes(), pubObject)
	// 	pubkey = "0x" + hex.EncodeToString(pubObject[:])
	// }

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
// func ProtoBaseAccount() Account {
// 	return BaseAccount{}
// }

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr string) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() string {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr string) error {
	if len(acc.Address) != 0 && acc.Address != "" {
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
	if acc.PubKey != nil && acc.Address != "" &&
		acc.PubKey.Address().String() != acc.Address {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &acc)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

//
// light base account
//

// LightBaseAccount - a base account structure.
type LightBaseAccount struct {
	Address       string `json:"address" yaml:"address"`
	AccountNumber uint64 `json:"account_number" yaml:"account_number"`
	Sequence      uint64 `json:"sequence" yaml:"sequence"`
}
