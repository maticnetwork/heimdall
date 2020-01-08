package types

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/cbergoon/merkletree"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/bor/crypto"
)

// DividendAccount contains Rewards, Shares, Slashed Amount
type DividendAccount struct {
	ID            DividendAccountID `json:"ID"`
	FeeAmount     string            `json:"feeAmount"`     // string representation of big.Int
	SlashedAmount string            `json:"slashedAmount"` // string representation of big.Int
}

// DividendAccountID  dividend account ID and helper functions
type DividendAccountID uint64

// NewDividendAccountID generate new dividendAccount ID
func NewDividendAccountID(id uint64) DividendAccountID {
	return DividendAccountID(id)
}

// Bytes get bytes of dividendAccountID
func (dividendAccountID DividendAccountID) Bytes() []byte {
	return []byte(strconv.Itoa(dividendAccountID.Int()))
}

// Int converts dividendAccount ID to int
func (dividendAccountID DividendAccountID) Int() int {
	return int(dividendAccountID)
}

// Uint64 converts dividendAccount ID to int
func (dividendAccountID DividendAccountID) Uint64() uint64 {
	return uint64(dividendAccountID)
}

func (da *DividendAccount) String() string {
	if da == nil {
		return "nil-DividendAccount"
	}

	return fmt.Sprintf("DividendAccount{%v %v %v %v}",
		da.ID,
		da.FeeAmount,
		da.SlashedAmount)
}

// MarshallDividendAccount - amino Marshall DividendAccount
func MarshallDividendAccount(cdc *codec.Codec, dividendAccount DividendAccount) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(dividendAccount)
	if err != nil {
		return bz, err
	}

	return bz, nil
}

// UnMarshallDividendAccount - amino Unmarshall DividendAccount
func UnMarshallDividendAccount(cdc *codec.Codec, value []byte) (DividendAccount, error) {

	var dividendAccount DividendAccount
	err := cdc.UnmarshalBinaryBare(value, &dividendAccount)
	if err != nil {
		return dividendAccount, err
	}
	return dividendAccount, nil
}

// SortDividendAccountByID - Sorts DividendAccounts  By  ID
func SortDividendAccountByID(dividendAccounts []DividendAccount) []DividendAccount {
	sort.Slice(dividendAccounts, func(i, j int) bool { return dividendAccounts[i].ID < dividendAccounts[j].ID })
	return dividendAccounts
}

//CalculateHash hashes the values of a DividendAccount
func (da DividendAccount) CalculateHash() ([]byte, error) {
	h := sha256.New()
	reward, _ := big.NewInt(0).SetString(da.FeeAmount, 10)
	slashAmount, _ := big.NewInt(0).SetString(da.SlashedAmount, 10)
	divAccountHash := crypto.Keccak256(appendBytes32(
		new(big.Int).SetUint64(uint64(da.ID)).Bytes(),
		reward.Bytes(),
		slashAmount.Bytes(),
	))
	var arr [32]byte
	copy(arr[:], divAccountHash)

	if _, err := h.Write(arr[:]); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func appendBytes32(data ...[]byte) []byte {
	var result []byte
	for _, v := range data {
		paddedV, err := convertTo32(v)
		if err == nil {
			result = append(result, paddedV[:]...)
		}
	}
	return result
}

func convertTo32(input []byte) (output [32]byte, err error) {
	l := len(input)
	if l > 32 || l == 0 {
		return
	}
	copy(output[32-l:], input[:])
	return
}

//Equals tests for equality of two Contents
func (da DividendAccount) Equals(other merkletree.Content) (bool, error) {
	return da.ID == other.(DividendAccount).ID, nil
}
