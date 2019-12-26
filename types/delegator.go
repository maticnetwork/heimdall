package types

import (
	"strconv"
)

// DelegatorID  delegator ID and helper functions
type DelegatorID uint64

// NewDelegatorID generate new delegator ID
func NewDelegatorID(id uint64) DelegatorID {
	return DelegatorID(id)
}

// Bytes get bytes of delegatorID
func (delegatorID DelegatorID) Bytes() []byte {
	return []byte(strconv.Itoa(delegatorID.Int()))
}

// Int converts delegator ID to int
func (delegatorID DelegatorID) Int() int {
	return int(delegatorID)
}

// Uint64 converts delegator ID to int
func (delegatorID DelegatorID) Uint64() uint64 {
	return uint64(delegatorID)
}
