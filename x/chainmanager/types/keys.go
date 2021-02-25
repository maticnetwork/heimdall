package types

const (
	// ModuleName defines the module name
	ModuleName = "chainmanager"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)
var (
	// ProposerKeyPrefix prefix for proposer
	ProposerKeyPrefix = []byte("proposer")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// ProposerKey returns proposer key
func ProposerKey() []byte {
	return ProposerKeyPrefix
}
