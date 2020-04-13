package types

import "time"

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	Signer HeimdallAddress `json:"signer"`

	// height at which validator was first a candidate OR was unjailed
	StartHeight int64 `json:"startHeight"`
	// index offset into signed block bit array
	IndexOffset int64 `json:"indexOffset"`
	// timestamp validator cannot be unjailed until
	JailedUntil time.Time `json:"jailedUntil"`
	// whether or not a validator has been tombstoned (killed out of validator set)
	// Tombstoned bool `protobuf:"varint,5,opt,name=tombstoned,proto3" json:"tombstoned,omitempty"`
	// missed blocks counter (to avoid scanning the array every time)
	MissedBlocksCounter int64 `json:"missed_blocks_counter,omitempty"`
}
