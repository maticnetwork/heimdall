package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// CheckpointBlockHeader block header struct
type CheckpointBlockHeader struct {
	Proposer   HeimdallAddress `json:"proposer"`
	StartBlock uint64          `json:"startBlock"`
	EndBlock   uint64          `json:"endBlock"`
	RootHash   common.Hash     `json:"rootHash"`
	TimeStamp  uint64          `json:"timestamp"`
}

// CreateBlock generate new block
func CreateBlock(start uint64, end uint64, rootHash common.Hash, proposer HeimdallAddress, timestamp uint64) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
		TimeStamp:  timestamp,
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
		m.TimeStamp,
	)
}
