package types

import (
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type CheckpointBlockHeader struct {
	Proposer   common.Address
	StartBlock uint64
	EndBlock   uint64
	RootHash   common.Hash
	TimeStamp  time.Time
}

func CreateBlock(start uint64, end uint64, rootHash common.Hash, proposer common.Address) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
		TimeStamp:  time.Now().UTC(),
	}
}

// add JSON marshaller and Unmarshaller here
