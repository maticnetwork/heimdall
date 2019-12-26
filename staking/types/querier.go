package types

import "github.com/maticnetwork/heimdall/types"

// query endpoints supported by the staking Querier
const (
	QueryCurrentValidatorSet  = "current-validator-set"
	QuerySigner               = "signer"
	QueryValidator            = "validator"
	QueryValidatorStatus      = "validator-status"
	QueryProposer             = "proposer"
	QueryCurrentProposer      = "current-proposer"
	QueryProposerBonusPercent = "proposer-bonus-percent"
)

// QuerySignerParams defines the params for querying by address
type QuerySignerParams struct {
	SignerAddress []byte `json:"signer_address"`
}

// NewQuerySignerParams creates a new instance of QuerySignerParams.
func NewQuerySignerParams(signerAddress []byte) QuerySignerParams {
	return QuerySignerParams{SignerAddress: signerAddress}
}

// QueryValidatorParams defines the params for querying val status.
type QueryValidatorParams struct {
	ValidatorID types.ValidatorID `json:"validator_id"`
}

// NewQueryValidatorParams creates a new instance of QueryValidatorParams.
func NewQueryValidatorParams(validatorID types.ValidatorID) QueryValidatorParams {
	return QueryValidatorParams{ValidatorID: validatorID}
}

// QueryProposerParams defines the params for querying val status.
type QueryProposerParams struct {
	Times uint64 `json:"times"`
}

// NewQueryProposerParams creates a new instance of QueryProposerParams.
func NewQueryProposerParams(times uint64) QueryProposerParams {
	return QueryProposerParams{Times: times}
}
