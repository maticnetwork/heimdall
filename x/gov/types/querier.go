package types

import (
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	QueryParams    = "params"
	QueryProposals = "proposals"
	QueryProposal  = "proposal"
	QueryDeposits  = "deposits"
	QueryDeposit   = "deposit"
	QueryVotes     = "votes"
	QueryVote      = "vote"
	QueryTally     = "tally"

	ParamDeposit  = "deposit"
	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

// QueryProposalsParams Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Limit          uint64
	VoterID        hmTypes.ValidatorID
	DepositorID    hmTypes.ValidatorID
	ProposalStatus ProposalStatus
}

// NewQueryProposalsParams creates a new instance of QueryProposalsParams
func NewQueryProposalsParams(limit uint64, status ProposalStatus, voterID, depositorID hmTypes.ValidatorID) QueryProposalsParams {
	return QueryProposalsParams{
		Limit:          limit,
		VoterID:        voterID,
		DepositorID:    depositorID,
		ProposalStatus: status,
	}
}
