package types

import (
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// query endpoints supported by the governance Querier
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

// Params for queries:
// - 'custom/gov/proposal'
// - 'custom/gov/deposits'
// - 'custom/gov/tally'
// - 'custom/gov/votes'
type QueryProposalParams struct {
	ProposalID uint64
}

// creates a new instance of QueryProposalParams
func NewQueryProposalParams(proposalID uint64) QueryProposalParams {
	return QueryProposalParams{
		ProposalID: proposalID,
	}
}

// Params for query 'custom/gov/deposit'
type QueryDepositParams struct {
	ProposalID uint64
	Depositor  hmTypes.ValidatorID
}

// creates a new instance of QueryDepositParams
func NewQueryDepositParams(proposalID uint64, depositor hmTypes.ValidatorID) QueryDepositParams {
	return QueryDepositParams{
		ProposalID: proposalID,
		Depositor:  depositor,
	}
}

// Params for query 'custom/gov/vote'
type QueryVoteParams struct {
	ProposalID uint64
	Voter      hmTypes.ValidatorID
}

// creates a new instance of QueryVoteParams
func NewQueryVoteParams(proposalID uint64, voter hmTypes.ValidatorID) QueryVoteParams {
	return QueryVoteParams{
		ProposalID: proposalID,
		Voter:      voter,
	}
}

// Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Voter          hmTypes.ValidatorID
	Depositor      hmTypes.ValidatorID
	ProposalStatus ProposalStatus
	Limit          uint64
}

// creates a new instance of QueryProposalsParams
func NewQueryProposalsParams(status ProposalStatus, limit uint64, voter hmTypes.ValidatorID, depositor hmTypes.ValidatorID) QueryProposalsParams {
	return QueryProposalsParams{
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: status,
		Limit:          limit,
	}
}
