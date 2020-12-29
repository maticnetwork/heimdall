package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewEventRecord creates new record
func NewEventRecord(
	txHash []byte,
	logIndex uint64,
	id uint64,
	contract sdk.AccAddress,
	data []byte,
	chainID string,
	recordTime time.Time,
) EventRecord {
	return EventRecord{
		Id:         id,
		Contract:   contract,
		Data:       data,
		TxHash:     txHash,
		LogIndex:   logIndex,
		ChainId:    chainID,
		RecordTime: recordTime,
	}
}
