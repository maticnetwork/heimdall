package common

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -------------- KEYS/CONSTANTS

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ACKCountKey             = []byte{0x11} // key to store ACK count
	BufferCheckpointKey     = []byte{0x12} // Key to store checkpoint in buffer
	HeaderBlockKey          = []byte{0x13} // prefix key for when storing header after ACK
	CheckpointCacheKey      = []byte{0x14} // key to store Cache for checkpoint
	CheckpointACKCacheKey   = []byte{0x15} // key to store Cache for checkpointACK
	CheckpointNoACKCacheKey = []byte{0x16} // key to store last no-ack

	ValidatorsKey          = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey        = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey = []byte{0x23} // Key to store current validator set

	SpanDurationKey       = []byte{0x24} // Key to store span duration for Bor
	LastSpanStartBlockKey = []byte{0x25} // Key to store last span start block
	SpanPrefixKey         = []byte{0x26} // prefix key to store span
)

//
// master keeper
//

// Keeper stores all related data
type Keeper struct {
	MasterKey     sdk.StoreKey
	cdc           *codec.Codec
	CheckpointKey sdk.StoreKey
	StakingKey    sdk.StoreKey
	BorKey        sdk.StoreKey
	// codespace
	Codespace sdk.CodespaceType
}

// NewKeeper create new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, stakingKey sdk.StoreKey, checkpointKey sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		MasterKey:     key,
		cdc:           cdc,
		Codespace:     codespace,
		CheckpointKey: checkpointKey,
		StakingKey:    stakingKey,
	}
	return keeper
}
