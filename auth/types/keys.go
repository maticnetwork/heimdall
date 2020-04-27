package types

import (
	"github.com/maticnetwork/heimdall/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "auth"

	// StoreKey is the store key string for auth
	StoreKey = ModuleName

	// RouterKey is the message route for auth
	RouterKey = ModuleName

	// QuerierRoute is the querier route for auth
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// FeeStoreKey is a string representation of the store key for fees
	FeeStoreKey = "fee"

	// FeeCollectorName the root string for the fee collector account address
	FeeCollectorName = "fee_collector"

	// FeeToken fee token name
	FeeToken = "matic"
)

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	// ProposerKeyPrefix prefix for proposer
	ProposerKeyPrefix = []byte("proposer")

	// GlobalAccountNumberKey param key for global account number
	GlobalAccountNumberKey = []byte("globalAccountNumber")
)

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr types.HeimdallAddress) []byte {
	return append(AddressStoreKeyPrefix, addr.Bytes()...)
}

// ProposerKey returns proposer key
func ProposerKey() []byte {
	return ProposerKeyPrefix
}
