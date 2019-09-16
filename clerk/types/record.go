package types

import (
	"encoding/hex"
	"fmt"

	"github.com/maticnetwork/heimdall/types"
)

// EventRecord represents state record
type EventRecord struct {
	ID       uint64                `json:"id" yaml:"id"`
	Contract types.HeimdallAddress `json:"contract" yaml:"contract"`
	Data     []byte                `json:"data" yaml:"data"`
}

// NewEventRecord creates new record
func NewEventRecord(id uint64, contract types.HeimdallAddress, data []byte) EventRecord {
	return EventRecord{
		ID:       id,
		Contract: contract,
		Data:     data,
	}
}

// String returns the string representatin of span
func (s *EventRecord) String() string {
	return fmt.Sprintf(
		"EventRecord: id %v, contract %v, data: %v",
		s.ID,
		s.Contract.String(),
		hex.EncodeToString(s.Data),
	)
}
