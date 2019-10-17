package types

import (
	"fmt"
)

// CheckpointBlockHeader block header struct
type CheckpointBlockHeader struct {
	Proposer       HeimdallAddress `json:"proposer"`
	StartBlock     uint64          `json:"startBlock"`
	EndBlock       uint64          `json:"endBlock"`
	RootHash       HeimdallHash    `json:"rootHash"`
	RewardRootHash HeimdallHash    `json:"rewardRootHash"`
	TimeStamp      uint64          `json:"timestamp"`
}

// CreateBlock generate new block
func CreateBlock(start uint64, end uint64, rootHash HeimdallHash, rewardRootHash HeimdallHash, proposer HeimdallAddress, timestamp uint64) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock:     start,
		EndBlock:       end,
		RootHash:       rootHash,
		RewardRootHash: rewardRootHash,
		Proposer:       proposer,
		TimeStamp:      timestamp,
	}
}

// String returns human redable string
func (m CheckpointBlockHeader) String() string {
	return fmt.Sprintf(
		"CheckpointBlockHeader {%v (%d:%d) %v %v}",
		m.Proposer.String(),
		m.StartBlock,
		m.EndBlock,
		m.RootHash.Hex(),
		m.RewardRootHash.Hex(),
		m.TimeStamp,
	)
}
