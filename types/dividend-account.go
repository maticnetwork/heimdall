package types

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/cbergoon/merkletree"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/bor/crypto"
)

// DividendAccount contains burned Fee amount
type DividendAccount struct {
	User      HeimdallAddress `json:"user"`
	FeeAmount string          `json:"feeAmount"` // string representation of big.Int
}

func NewDividendAccount(user HeimdallAddress, fee string) DividendAccount {
	return DividendAccount{
		User:      user,
		FeeAmount: fee,
	}
}

func (da *DividendAccount) String() string {
	if da == nil {
		return "nil-DividendAccount"
	}

	return fmt.Sprintf("DividendAccount{%s %v}",
		da.User.EthAddress(),
		da.FeeAmount)
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

// SortDividendAccountByAddress - Sorts DividendAccounts  By  Address
func SortDividendAccountByAddress(dividendAccounts []DividendAccount) []DividendAccount {
	sort.Slice(dividendAccounts, func(i, j int) bool {
		return bytes.Compare(dividendAccounts[i].User.Bytes(), dividendAccounts[j].User.Bytes()) < 0
	})
	return dividendAccounts
}

//CalculateHash hashes the values of a DividendAccount
func (da DividendAccount) CalculateHash() ([]byte, error) {
	fee, _ := big.NewInt(0).SetString(da.FeeAmount, 10)
	divAccountHash := crypto.Keccak256(appendBytes32(
		da.User.Bytes(),
		fee.Bytes(),
	))

	return divAccountHash, nil
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
	return da.User.Equals(other.(DividendAccount).User), nil
}
