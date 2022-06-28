package types

import (
	"fmt"

	"github.com/maticnetwork/heimdall/types"
)

var (
	NQuerier NodeQuerier
)

// NodeQuerier is an interface that is satisfied by types that provide the QueryWithData method
//go:generate mockgen -destination=./mocks/node_querier_mock.go -package=mocks . NodeQuerier
type NodeQuerier interface {
	// QueryWithData performs a query to a Tendermint node with the provided path
	// and a data payload. It returns the result and height of the query upon success
	// or an error if the query fails.
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// AccountRetriever defines the properties of a type that can be used to
// retrieve accounts.
type AccountRetriever struct {
	querier NodeQuerier
}

// NewAccountRetriever initialises a new AccountRetriever instance.
func NewAccountRetriever(querier NodeQuerier) AccountRetriever {
	NQuerier = querier
	return AccountRetriever{querier: querier}
}

// GetAccount queries for an account given an address and a block height. An
// error is returned if the query or decoding fails.
func (ar AccountRetriever) GetAccount(addr types.HeimdallAddress) (Account, error) {
	account, _, err := ar.GetAccountWithHeight(addr)
	return account, err
}

// GetAccountWithHeight queries for an account given an address. Returns the
// height of the query with the account. An error is returned if the query
// or decoding fails.
func (ar AccountRetriever) GetAccountWithHeight(addr types.HeimdallAddress) (Account, int64, error) {
	bs, err := ModuleCdc.MarshalJSON(NewQueryAccountParams(addr))
	if err != nil {
		return nil, 0, err
	}

	res, height, err := NQuerier.QueryWithData(fmt.Sprintf("custom/%s/%s", QuerierRoute, QueryAccount), bs)
	if err != nil {
		return nil, height, err
	}

	var account Account
	if err := ModuleCdc.UnmarshalJSON(res, &account); err != nil {
		return nil, height, err
	}

	return account, height, nil
}

// EnsureExists returns an error if no account exists for the given address else nil.
func (ar AccountRetriever) EnsureExists(addr types.HeimdallAddress) error {
	if _, err := ar.GetAccount(addr); err != nil {
		return err
	}

	return nil
}

// GetAccountNumberSequence returns sequence and account number for the given address.
// It returns an error if the account couldn't be retrieved from the state.
func (ar AccountRetriever) GetAccountNumberSequence(addr types.HeimdallAddress) (uint64, uint64, error) {
	acc, err := ar.GetAccount(addr)
	if err != nil {
		return 0, 0, err
	}

	return acc.GetAccountNumber(), acc.GetSequence(), nil
}
