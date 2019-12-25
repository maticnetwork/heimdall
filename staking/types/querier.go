package types

// query endpoints supported by the staking Querier
const (
	QueryValidatorStatus      = "validator-status"
	QueryProposerBonusPercent = "proposer-bonus-percent"
	QueryCurrentValidatorSet  = "current-validator-set"
)

// QueryValidatorStatusParams defines the params for querying val status.
type QueryValidatorStatusParams struct {
	SignerAddress []byte
}

// NewQueryValidatorStatusParams creates a new instance of QueryValidatorStatusParams.
func NewQueryValidatorStatusParams(signerAddress []byte) QueryValidatorStatusParams {
	return QueryValidatorStatusParams{SignerAddress: signerAddress}
}
