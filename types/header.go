package types

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// CheckpointBlockHeader block header struct
type CheckpointBlockHeader struct {
	Proposer   common.Address
	StartBlock uint64
	EndBlock   uint64
	RootHash   common.Hash
	TimeStamp  time.Time
}

// CreateBlock generate new block
func CreateBlock(start uint64, end uint64, rootHash common.Hash, proposer common.Address) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
		TimeStamp:  time.Now().UTC(),
	}
}

// String returns human redable string
func (m CheckpointBlockHeader) String() string {
	return fmt.Sprintf(
		"CheckpointBlockHeader {%v (%d:%d) %v %v}",
		m.Proposer.Hex(),
		m.StartBlock,
		m.EndBlock,
		m.RootHash.Hex(),
		m.TimeStamp.String(),
	)
}
