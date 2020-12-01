package types

import (
	"bytes"
	"errors"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	proto "github.com/gogo/protobuf/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
// type (
// 	AccountI = exported.AccountI
// )

//-----------------------------------------------------------------------------
// // BaseAccount
// var (
// 	_ AccountI = (*BaseAccount)(nil)
// )

// NewBaseAccount creates a new BaseAccount object
//nolint:interfacer
func NewBaseAccount(address string, pubKey crypto.PubKey, accountNumber, sequence uint64) *BaseAccount {
	acc := &BaseAccount{
		Address:       address,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}

	err := acc.SetPubKey(pubKey)
	if err != nil {
		panic(err)
	}

	return acc
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() AccountI {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr string) *BaseAccount {
	return &BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.AccountI.
func (acc BaseAccount) GetAddress() string {
	return acc.Address
}

// SetAddress - Implements sdk.AccountI.
func (acc *BaseAccount) SetAddress(addr string) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.AccountI.
func (acc BaseAccount) GetPubKey() (pk crypto.PubKey) {
	return acc.PubKey
}

// SetPubKey - Implements sdk.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// SetCoins - Implements sdk.AccountI.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetCoins - Implements sdk.AccountI.
func (acc *BaseAccount) GetCoins() (coins sdk.Coins) {
	return acc.Coins
}

// GetAccountNumber - Implements AccountI
func (acc BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements AccountI
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.AccountI.
func (acc BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.AccountI.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.Address == "" || acc.PubKey == nil {
		return nil
	}

	accAddr, err := sdk.AccAddressFromBech32(acc.Address)
	if err != nil {
		return err
	}

	if !bytes.Equal(acc.GetPubKey().Address().Bytes(), accAddr.Bytes()) {
		return errors.New("account address and pubkey address do not match")
	}

	return nil
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
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

// AccountI is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements AccountI.
type AccountI interface {
	proto.Message

	GetAddress() string
	SetAddress(string) error // errors if already set.

	GetCoins() sdk.Coins
	SetCoins(sdk.Coins) error

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	// Ensure that account implements stringer
	String() string
}
