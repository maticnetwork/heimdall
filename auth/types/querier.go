package types

import (
	"github.com/maticnetwork/heimdall/types"
)

// query endpoints supported by the auth Querier
const (
	QueryParams  = "params"
	QueryAccount = "account"
)

// QueryAccountParams defines the params for querying accounts.
type QueryAccountParams struct {
	Address types.HeimdallAddress
}

// NewQueryAccountParams creates a new instance of QueryAccountParams.
func NewQueryAccountParams(addr types.HeimdallAddress) QueryAccountParams {
	return QueryAccountParams{Address: addr}
}
