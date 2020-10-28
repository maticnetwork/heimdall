package types

import (
	"fmt"
	"sort"

	"github.com/maticnetwork/heimdall/types/common"
)

// CreateBlock generate new block
func CreateBlock(
	start uint64,
	end uint64,
	rootHash common.HeimdallHash,
	proposer common.HeimdallAddress,
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

// SortHeaders sorts array of headers on the basis for timestamps
func SortHeaders(headers []Checkpoint) []Checkpoint {
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].TimeStamp < headers[j].TimeStamp
	})
	return headers
}

// String returns human redable string
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
