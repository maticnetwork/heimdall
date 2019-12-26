package types

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DividendAccount contains Rewards, Shares, Slashed Amount
type DividendAccount struct {
	ID            DividendAccountID `json:"ID"`
	Shares        float32           `json:"shares"`
	RewardAmount  string            `json:"rewardAmount"`
	SlashedAmount string            `json:"slashedAmount"`
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
		da.Shares,
		da.RewardAmount,
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
