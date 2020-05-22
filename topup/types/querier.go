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
	UserAddress types.HeimdallAddress `json:"user_addr"`
}

// NewQueryDividendAccountParams creates a new instance of QueryDividendAccountParams.
func NewQueryDividendAccountParams(userAddress types.HeimdallAddress) QueryDividendAccountParams {
	return QueryDividendAccountParams{UserAddress: userAddress}
}

// QueryAccountProofParams defines the params for querying account proof.
type QueryAccountProofParams struct {
	UserAddress types.HeimdallAddress `json:"user_addr"`
}

// NewQueryAccountProofParams creates a new instance of QueryAccountProofParams.
func NewQueryAccountProofParams(userAddress types.HeimdallAddress) QueryAccountProofParams {
	return QueryAccountProofParams{UserAddress: userAddress}
}

// QueryVerifyAccountProofParams defines the params for verifying account proof.
type QueryVerifyAccountProofParams struct {
	UserAddress  types.HeimdallAddress `json:"user_addr"`
	AccountProof string                `json:"account_proof"`
}

// NewQueryVerifyAccountProofParams creates a new instance of QueryVerifyAccountProofParams.
func NewQueryVerifyAccountProofParams(userAddress types.HeimdallAddress, accountProof string) QueryVerifyAccountProofParams {
	return QueryVerifyAccountProofParams{UserAddress: userAddress, AccountProof: accountProof}
}
