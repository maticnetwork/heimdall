package types

import (
	"fmt"

	"github.com/maticnetwork/heimdall/types"
)

// EventRecord represents state record
type EventRecord struct {
	ID       uint64                `json:"id" yaml:"id"`
	Contract types.HeimdallAddress `json:"contract" yaml:"contract"`
	Data     types.HexBytes        `json:"data" yaml:"data"`
	TxHash   types.HeimdallHash    `json:"tx_hash" yaml:"tx_hash"`
	LogIndex uint64                `json:"log_index" yaml:"log_index"`
}

// NewEventRecord creates new record
func NewEventRecord(
	txHash types.HeimdallHash,
	logIndex uint64,
	id uint64,
	contract types.HeimdallAddress,
	data types.HexBytes,
) EventRecord {
	return EventRecord{
		ID:       id,
		Contract: contract,
		Data:     data,
		TxHash:   txHash,
		LogIndex: logIndex,
	}
}

// String returns the string representatin of span
func (s *EventRecord) String() string {
	return fmt.Sprintf(
		"EventRecord: id %v, contract %v, data: %v, txHash: %v, logIndex: %v",
		s.ID,
		s.Contract.String(),
		s.Data.String(),
		s.TxHash.Hex(),
		s.LogIndex,
	)
}
