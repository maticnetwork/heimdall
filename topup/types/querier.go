package types

import "github.com/maticnetwork/heimdall/types"

const (
	QuerySequence            = "sequence"
	QueryDividendAccount     = "dividend-account"
	QueryDividendAccountRoot = "dividend-account-root"
	QueryAccountProof        = "dividend-account-proof"
	QueryVerifyAccountProof  = "verify-account-proof"
)

// QuerySequenceParams defines the params for querying an account Sequence.
type QuerySequenceParams struct {
	TxHash   string
	LogIndex uint64
}

// NewQuerySequenceParams creates a new instance of QuerySequenceParams.
func NewQuerySequenceParams(txHash string, logIndex uint64) QuerySequenceParams {
	return QuerySequenceParams{TxHash: txHash, LogIndex: logIndex}
}

// QueryDividendAccountParams defines the params for querying dividend account status.
type QueryDividendAccountParams struct {
	DividendAccountID types.DividendAccountID `json:"dividend_account_id"`
}

// NewQueryDividendAccountParams creates a new instance of QueryDividendAccountParams.
func NewQueryDividendAccountParams(dividendAccountID types.DividendAccountID) QueryDividendAccountParams {
	return QueryDividendAccountParams{DividendAccountID: dividendAccountID}
}

// QueryAccountProofParams defines the params for querying account proof.
type QueryAccountProofParams struct {
	DividendAccountID types.DividendAccountID `json:"dividend_account_id"`
}

// NewQueryAccountProofParams creates a new instance of QueryAccountProofParams.
func NewQueryAccountProofParams(dividendAccountID types.DividendAccountID) QueryAccountProofParams {
	return QueryAccountProofParams{DividendAccountID: dividendAccountID}
}

// QueryVerifyAccountProofParams defines the params for verifying account proof.
type QueryVerifyAccountProofParams struct {
	DividendAccountID types.DividendAccountID `json:"dividend_account_id"`
	AccountProof      string                  `json:"account_proof"`
}

// NewQueryVerifyAccountProofParams creates a new instance of QueryVerifyAccountProofParams.
func NewQueryVerifyAccountProofParams(dividendAccountID types.DividendAccountID, accountProof string) QueryVerifyAccountProofParams {
	return QueryVerifyAccountProofParams{DividendAccountID: dividendAccountID, AccountProof: accountProof}
}
