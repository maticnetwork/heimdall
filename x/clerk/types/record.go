package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
)

// NewEventRecord creates new record
func NewEventRecord(
	txHash hmCommon.HeimdallHash,
	logIndex uint64,
	id uint64,
	contract sdk.AccAddress,
	data []byte,
	chainID string,
	recordTime time.Time,
) EventRecord {
	return EventRecord{
		Id:         id,
		Contract:   contract.String(),
		Data:       data,
		TxHash:     txHash.String(),
		LogIndex:   logIndex,
		ChainId:    chainID,
		RecordTime: recordTime,
	}
}
