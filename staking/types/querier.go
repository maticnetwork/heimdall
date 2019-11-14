package types

// query endpoints supported by the staking Querier
const (
	QueryValStatus = "val-status"
)

// QueryValStatusParams defines the params for querying val status.
type QueryValStatusParams struct {
	SignerAddress []byte
}

// NewQueryValStatusParams creates a new instance of NewQueryValStatusParams.
func NewQueryValStatusParams(signerAddress []byte) QueryValStatusParams {
	return QueryValStatusParams{SignerAddress: signerAddress}
}
