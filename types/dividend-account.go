package types

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/cbergoon/merkletree"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/crypto"
)

// DividendAccount contains burned Fee amount
// type DividendAccount struct {
// 	User      HeimdallAddress `json:"user"`
// 	FeeAmount string          `json:"feeAmount"` // string representation of big.Int
// }

func NewDividendAccount(user sdk.AccAddress, fee string) DividendAccount {
	return DividendAccount{
		User:      user.String(),
		FeeAmount: fee,
	}
}

func (da *DividendAccount) String() string {
	if da == nil {
		return "nil-DividendAccount"
	}

	return fmt.Sprintf("DividendAccount{%s %v}", //nolint
		da.User,
		da.FeeAmount)
}

// MarshallDividendAccount - amino Marshall DividendAccount
func MarshallDividendAccount(cdc codec.BinaryMarshaler, dividendAccount *DividendAccount) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(dividendAccount)
	if err != nil {
		return bz, err
	}

	return bz, nil
}

// UnMarshallDividendAccount - amino Unmarshall DividendAccount
func UnMarshallDividendAccount(cdc codec.BinaryMarshaler, value []byte) (DividendAccount, error) {

	var dividendAccount DividendAccount
	err := cdc.UnmarshalBinaryBare(value, &dividendAccount)
	if err != nil {
		return dividendAccount, err
	}
	return dividendAccount, nil
}

// SortDividendAccountByAddress - Sorts DividendAccounts  By  Address
func SortDividendAccountByAddress(dividendAccounts []*DividendAccount) []*DividendAccount {
	sort.Slice(dividendAccounts, func(i, j int) bool {
		return dividendAccounts[i].User == dividendAccounts[j].User
	})
	return dividendAccounts
}

//CalculateHash hashes the values of a DividendAccount
func (da DividendAccount) CalculateHash() ([]byte, error) {
	fee, _ := big.NewInt(0).SetString(da.FeeAmount, 10)
	divAccountHash := crypto.Keccak256(appendBytes32(
		[]byte(da.User),
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
	return da.User == other.(DividendAccount).User, nil
}
