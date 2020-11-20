package types

const (
	// ModuleName defines the module name
	ModuleName = "auth"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"

	FeeToken = "matic"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	// AddressStoreKeyPrefix prefix for account-by-address store
	AddressStoreKeyPrefix = []byte{0x01}

	// ProposerKeyPrefix prefix for proposer
	ProposerKeyPrefix = []byte("proposer")

	// GlobalAccountNumberKey param key for global account number
	GlobalAccountNumberKey = []byte("globalAccountNumber")
)

// AddressStoreKey turn an address to key used to get it from the account store
func AddressStoreKey(addr string) []byte {
	return append(AddressStoreKeyPrefix, addr...)
}

// ProposerKey returns proposer key
func ProposerKey() []byte {
	return ProposerKeyPrefix
}
