package types

import (
	"fmt"
	"sort"
)

// Checkpoint block header struct
type Checkpoint struct {
	Proposer   HeimdallAddress `json:"proposer"`
	StartBlock uint64          `json:"start_block"`
	EndBlock   uint64          `json:"end_block"`
	RootHash   HeimdallHash    `json:"root_hash"`
	BorChainID string          `json:"bor_chain_id"`
	TimeStamp  uint64          `json:"timestamp"`
}

// Milestone block header struct
type Milestone struct {
	Proposer    HeimdallAddress `json:"proposer"`
	StartBlock  uint64          `json:"start_block"`
	EndBlock    uint64          `json:"end_block"`
	Hash        HeimdallHash    `json:"hash"`
	BorChainID  string          `json:"bor_chain_id"`
	MilestoneID string          `json:"milestone_id"`
	TimeStamp   uint64          `json:"timestamp"`
}

// CreateBlock generate new block
func CreateBlock(
	start uint64,
	end uint64,
	rootHash HeimdallHash,
	proposer HeimdallAddress,
	borChainID string,
	timestamp uint64,
) Checkpoint {
	return Checkpoint{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
		BorChainID: borChainID,
		TimeStamp:  timestamp,
	}
}

// CreateBlock generate new block
func CreateMilestone(
	start uint64,
	end uint64,
	hash HeimdallHash,
	proposer HeimdallAddress,
	borChainID string,
	milestoneID string,
	timestamp uint64,
) Milestone {
	return Milestone{
		StartBlock:  start,
		EndBlock:    end,
		Hash:        hash,
		Proposer:    proposer,
		BorChainID:  borChainID,
		MilestoneID: milestoneID,
		TimeStamp:   timestamp,
	}
}

// SortHeaders sorts array of headers on the basis for timestamps
func SortHeaders(headers []Checkpoint) []Checkpoint {
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].TimeStamp < headers[j].TimeStamp
	})

	return headers
}

// String returns human readable string
func (m Checkpoint) String() string {
	return fmt.Sprintf(
		"Checkpoint {%v (%d:%d) %v %v %v}",
		m.Proposer.String(),
		m.StartBlock,
		m.EndBlock,
		m.RootHash.Hex(),
		m.BorChainID,
		m.TimeStamp,
	)
}

// String returns human readable string
func (m Milestone) String() string {
	return fmt.Sprintf(
		"Milestone {%v (%d:%d) %v %v %v %v}",
		m.Proposer.String(),
		m.StartBlock,
		m.EndBlock,
		m.Hash.Hex(),
		m.BorChainID,
		m.MilestoneID,
		m.TimeStamp,
	)
}
