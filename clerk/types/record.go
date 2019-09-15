package types

import (
	"encoding/hex"
	"fmt"

	"github.com/maticnetwork/heimdall/types"
)

// Record represents record
type Record struct {
	ID       uint64                `json:"id" yaml:"id"`
	Contract types.HeimdallAddress `json:"contract" yaml:"contract"`
	Data     []byte                `json:"data" yaml:"data"`
}

// NewRecord creates new record
func NewRecord(id uint64, contract types.HeimdallAddress, data []byte) Record {
	return Record{
		ID:       id,
		Contract: contract,
		Data:     data,
	}
}

// String returns the string representatin of span
func (s *Record) String() string {
	return fmt.Sprintf(
		"Record: id %v, contract %v, data: %v",
		s.ID,
		s.Contract,
		hex.EncodeToString(s.Data),
	)
}
