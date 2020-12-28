package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

var _ types.QueryServer = Keeper{}

// Proposal returns proposal details based on ProposalID
func (q Keeper) Proposal(c context.Context, req *types.QueryProposalRequest) (*types.QueryProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	proposal, found := q.GetProposal(ctx, req.ProposalId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
	}

	return &types.QueryProposalResponse{Proposal: proposal}, nil
}

// TallyResult queries the tally of a proposal vote
func (q Keeper) TallyResult(c context.Context, req *types.QueryTallyResultRequest) (*types.QueryTallyResultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	proposal, ok := q.GetProposal(ctx, req.ProposalId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "proposal %d doesn't exist", req.ProposalId)
	}

	var tallyResult types.TallyResult

	switch {
	case proposal.Status == types.StatusDepositPeriod:
		tallyResult = types.EmptyTallyResult()

	case proposal.Status == types.StatusPassed || proposal.Status == types.StatusRejected:
		tallyResult = proposal.FinalTallyResult

	default:
		// proposal is in voting period
		_, _, tallyResult = q.Tally(ctx, proposal)
	}

	return &types.QueryTallyResultResponse{Tally: tallyResult}, nil
}

// Deposits returns single proposal's all deposits
func (q Keeper) Deposits(c context.Context, req *types.QueryDepositsRequest) (*types.QueryDepositsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)
	deposits := q.GetDeposits(ctx, req.ProposalId)

	return &types.QueryDepositsResponse{Deposits: deposits}, nil
}

// Deposit queries single deposit information based proposalID, depositAddr
func (q Keeper) Deposit(c context.Context, req *types.QueryDepositRequest) (*types.QueryDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	if req.Depositor == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty depositor address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	deposit, found := q.GetDeposit(ctx, req.ProposalId, req.Depositor)
	if !found {
		return nil, status.Errorf(codes.InvalidArgument,
			"depositer: %v not found for proposal: %v", req.Depositor, req.ProposalId)
	}

	return &types.QueryDepositResponse{Deposit: deposit}, nil
}

// Params queries all params
func (q Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	switch req.ParamsType {
	case types.ParamDeposit:
		depositParmas := q.GetDepositParams(ctx)
		return &types.QueryParamsResponse{DepositParams: depositParmas}, nil

	case types.ParamVoting:
		votingParmas := q.GetVotingParams(ctx)
		return &types.QueryParamsResponse{VotingParams: votingParmas}, nil

	case types.ParamTallying:
		tallyParams := q.GetTallyParams(ctx)
		return &types.QueryParamsResponse{TallyParams: tallyParams}, nil

	default:
		return nil, status.Errorf(codes.InvalidArgument,
			"%s is not a valid parameter type", req.ParamsType)
	}
}

// Votes returns single proposal's votes
func (q Keeper) Votes(c context.Context, req *types.QueryVotesRequest) (*types.QueryVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	ctx := sdk.UnwrapSDKContext(c)

	votes := q.GetVotes(ctx, req.ProposalId)

	return &types.QueryVotesResponse{Votes: votes}, nil
}

// Vote returns Voted information based on proposalID, voterAddr
func (q Keeper) Vote(c context.Context, req *types.QueryVoteRequest) (*types.QueryVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.ProposalId == 0 {
		return nil, status.Error(codes.InvalidArgument, "proposal id can not be 0")
	}

	if req.Voter == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty voter address")
	}

	ctx := sdk.UnwrapSDKContext(c)

	vote, found := q.GetVote(ctx, req.ProposalId, req.Voter)
	if !found {
		return nil, status.Errorf(codes.InvalidArgument,
			"voter: %v not found for proposal: %v", req.Voter, req.ProposalId)
	}

	return &types.QueryVoteResponse{Vote: vote}, nil
}

func (q Keeper) Proposals(c context.Context, req *types.QueryProposalsRequest) (*types.QueryProposalsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	proposals := q.GetProposalsFiltered(ctx, req.Voter, req.Depositor, req.ProposalStatus, req.NumLimit)
	return &types.QueryProposalsResponse{Proposals: proposals}, nil
}
