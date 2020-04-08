package types

// ValidatorSigningInfo defines the signing info for a validator
type ValidatorSigningInfo struct {
	// Address github_com_cosmos_cosmos_sdk_types.ConsAddress `protobuf:"bytes,1,opt,name=address,proto3,casttype=github.com/cosmos/cosmos-sdk/types.ConsAddress" json:"address,omitempty"`
	// // height at which validator was first a candidate OR was unjailed
	// StartHeight int64 `protobuf:"varint,2,opt,name=start_height,json=startHeight,proto3" json:"start_height,omitempty" yaml:"start_height"`
	// // index offset into signed block bit array
	// IndexOffset int64 `protobuf:"varint,3,opt,name=index_offset,json=indexOffset,proto3" json:"index_offset,omitempty" yaml:"index_offset"`
	// // timestamp validator cannot be unjailed until
	// JailedUntil time.Time `protobuf:"bytes,4,opt,name=jailed_until,json=jailedUntil,proto3,stdtime" json:"jailed_until" yaml:"jailed_until"`
	// // whether or not a validator has been tombstoned (killed out of validator set)
	// Tombstoned bool `protobuf:"varint,5,opt,name=tombstoned,proto3" json:"tombstoned,omitempty"`
	// // missed blocks counter (to avoid scanning the array every time)
	// MissedBlocksCounter int64 `protobuf:"varint,6,opt,name=missed_blocks_counter,json=missedBlocksCounter,proto3" json:"missed_blocks_counter,omitempty" yaml:"missed_blocks_counter"`
}
