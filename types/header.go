package types

import (
	"fmt"
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

func GenEmptyCheckpointBlockHeader() CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: 0,
		EndBlock:   0,
	}
}

// add JSON marshaller and Unmarshaller here

func (m CheckpointBlockHeader) HumanReadableString() string {
	resp := "Checkpoint \n"

	resp += fmt.Sprintf("Proposer : %s\n", m.Proposer.String())
	resp += fmt.Sprintf("StartBlock: %d\n", m.StartBlock)
	resp += fmt.Sprintf("EndBlock: %d\n", m.EndBlock)
	resp += fmt.Sprintf("RootHash: %v\n", m.RootHash)
	resp += fmt.Sprintf("CreationTime: %v", m.TimeStamp.String())
	return resp
}
