package types

import (
	"encoding/binary"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "gov"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_capability"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// Keys for governance store
// Items are stored with the following key: values
//
// - 0x00<proposalID_Bytes>: Proposal
//
// - 0x01<endTime_Bytes><proposalID_Bytes>: activeProposalID
//
// - 0x02<endTime_Bytes><proposalID_Bytes>: inactiveProposalID
//
// - 0x03: nextProposalID
//
// - 0x10<proposalID_Bytes><depositorAddr_Bytes>: Deposit
//
// - 0x20<proposalID_Bytes><voterAddr_Bytes>: Voter
var (
	ProposalsKeyPrefix          = []byte{0x00}
	ActiveProposalQueuePrefix   = []byte{0x01}
	InactiveProposalQueuePrefix = []byte{0x02}
	ProposalIDKey               = []byte{0x03}

	DepositsKeyPrefix = []byte{0x10}

	VotesKeyPrefix = []byte{0x20}
)

var lenTime = len(sdk.FormatTimeBytes(time.Now()))

// ProposalKey gets a specific proposal from the store
func ProposalKey(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)
	return append(ProposalsKeyPrefix, bz...)
}

// GetProposalIDBytes returns the byte representation of the proposalID
func GetProposalIDBytes(proposalID uint64) (proposalIDBz []byte) {
	proposalIDBz = make([]byte, 8)
	binary.BigEndian.PutUint64(proposalIDBz, proposalID)
	return
}

// GetProposalIDFromBytes returns proposalID in uint64 format from a byte array
func GetProposalIDFromBytes(bz []byte) (proposalID uint64) {
	return binary.BigEndian.Uint64(bz)
}

// ActiveProposalByTimeKey gets the active proposal queue key by endTime
func ActiveProposalByTimeKey(endTime time.Time) []byte {
	return append(ActiveProposalQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// ActiveProposalQueueKey returns the key for a proposalID in the activeProposalQueue
func ActiveProposalQueueKey(proposalID uint64, endTime time.Time) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)

	return append(ActiveProposalByTimeKey(endTime), bz...)
}

// InactiveProposalByTimeKey gets the inactive proposal queue key by endTime
func InactiveProposalByTimeKey(endTime time.Time) []byte {
	return append(InactiveProposalQueuePrefix, sdk.FormatTimeBytes(endTime)...)
}

// InactiveProposalQueueKey returns the key for a proposalID in the inactiveProposalQueue
func InactiveProposalQueueKey(proposalID uint64, endTime time.Time) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)

	return append(InactiveProposalByTimeKey(endTime), bz...)
}

// DepositsKey gets the first part of the deposits key based on the proposalID
func DepositsKey(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)
	return append(DepositsKeyPrefix, bz...)
}

// DepositKey key of a specific deposit from the store
func DepositKey(proposalID uint64, validator hmTypes.ValidatorID) []byte {
	return append(DepositsKey(proposalID), validator.Bytes()...)
}

// VotesKey gets the first part of the votes key based on the proposalID
func VotesKey(proposalID uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, proposalID)
	return append(VotesKeyPrefix, bz...)
}

// VoteKey key of a specific vote from the store
func VoteKey(proposalID uint64, validator hmTypes.ValidatorID) []byte {
	return append(VotesKey(proposalID), validator.Bytes()...)
}

// Split keys function; used for iterators

// SplitProposalKey split the proposal key and returns the proposal id
func SplitProposalKey(key []byte) (proposalID uint64) {
	if len(key[1:]) != 8 {
		panic(fmt.Sprintf("unexpected key length (%d ≠ 8)", len(key[1:])))
	}

	return binary.LittleEndian.Uint64(key[1:])
}

// SplitActiveProposalQueueKey split the active proposal key and returns the proposal id and endTime
func SplitActiveProposalQueueKey(key []byte) (proposalID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// SplitInactiveProposalQueueKey split the inactive proposal key and returns the proposal id and endTime
func SplitInactiveProposalQueueKey(key []byte) (proposalID uint64, endTime time.Time) {
	return splitKeyWithTime(key)
}

// SplitKeyDeposit split the deposits key and returns the proposal id and depositor address
func SplitKeyDeposit(key []byte) (proposalID uint64, depositorAddr hmCommon.HeimdallAddress) {
	return splitKeyWithAddress(key)
}

// SplitKeyVote split the votes key and returns the proposal id and voter address
func SplitKeyVote(key []byte) (proposalID uint64, voterAddr hmCommon.HeimdallAddress) {
	return splitKeyWithAddress(key)
}

// private functions

func splitKeyWithTime(key []byte) (proposalID uint64, endTime time.Time) {
	if len(key[1:]) != 8+lenTime {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key[1:]), lenTime+8))
	}

	endTime, err := sdk.ParseTimeBytes(key[1 : 1+lenTime])
	if err != nil {
		panic(err)
	}
	proposalID = binary.LittleEndian.Uint64(key[1+lenTime:])
	return
}

func splitKeyWithAddress(key []byte) (proposalID uint64, addr hmCommon.HeimdallAddress) {
	if len(key[1:]) != 8+sdk.AddrLen {
		panic(fmt.Sprintf("unexpected key length (%d ≠ %d)", len(key), 8+sdk.AddrLen))
	}

	proposalID = binary.LittleEndian.Uint64(key[1:9])
	addr = hmCommon.BytesToHeimdallAddress(key[9:])
	return
}
