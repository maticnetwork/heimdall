package types

import (
	"fmt"
	"sort"
)

// CheckpointBlockHeader block header struct
type CheckpointBlockHeader struct {
	Proposer        HeimdallAddress `json:"proposer"`
	StartBlock      uint64          `json:"startBlock"`
	EndBlock        uint64          `json:"endBlock"`
	RootHash        HeimdallHash    `json:"rootHash"`
	AccountRootHash HeimdallHash    `json:"accountRootHash"`
	TimeStamp       uint64          `json:"timestamp"`
}

// CreateBlock generate new block
func CreateBlock(start uint64, end uint64, rootHash HeimdallHash, accountRootHash HeimdallHash, proposer HeimdallAddress, timestamp uint64) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock:      start,
		EndBlock:        end,
		RootHash:        rootHash,
		AccountRootHash: accountRootHash,
		Proposer:        proposer,
		TimeStamp:       timestamp,
	}
}

// SortHeaders sorts array of headers on the basis for timestamps
func SortHeaders(headers []CheckpointBlockHeader) []CheckpointBlockHeader {
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].TimeStamp < headers[j].TimeStamp
	})
	return headers
}

// String returns human redable string
func (m CheckpointBlockHeader) String() string {
	return fmt.Sprintf(
		"CheckpointBlockHeader {%v (%d:%d) %v %v %v}",
		m.Proposer.String(),
		m.StartBlock,
		m.EndBlock,
		m.RootHash.Hex(),
		m.AccountRootHash.Hex(),
		m.TimeStamp,
	)
}
