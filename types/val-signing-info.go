package types

import (
	"fmt"
)

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	ValID ValidatorID `json:"valID"`

	// height at which validator was first a candidate OR was unjailed
	StartHeight int64 `json:"startHeight"`
	// index offset into signed block bit array
	IndexOffset int64 `json:"indexOffset"`
	// timestamp validator cannot be unjailed until
	// JailedUntil time.Time `json:"jailedUntil"`
	// whether or not a validator has been tombstoned (killed out of validator set)
	// Tombstoned bool `protobuf:"varint,5,opt,name=tombstoned,proto3" json:"tombstoned,omitempty"`
	// missed blocks counter (to avoid scanning the array every time)
	MissedBlocksCounter int64 `json:"missed_blocks_counter,omitempty"`
}

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
func NewValidatorSigningInfo(
	valID ValidatorID, startHeight, indexOffset int64,
	missedBlocksCounter int64,
) ValidatorSigningInfo {
	return ValidatorSigningInfo{
		ValID:       valID,
		StartHeight: startHeight,
		IndexOffset: indexOffset,
		// JailedUntil: jailedUntil,
		// Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  valID:               %d
  Start Height:          %d
  Index Offset:          %d  
  Missed Blocks Counter: %d`,
		i.ValID, i.StartHeight, i.IndexOffset,
		i.MissedBlocksCounter)
}
