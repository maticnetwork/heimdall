package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SubmitProposal(goCtx context.Context, msg *types.MsgSubmitProposal) (*types.MsgSubmitProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := getValidValidator(ctx, k.Keeper, msg.Proposer, msg.Validator); err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	proposal, err := k.Keeper.SubmitProposal(ctx, msg.GetContent())
	if err != nil {
		return nil, err
	}

	err, votingStarted := k.Keeper.AddDeposit(ctx, proposal.ProposalId, msg.GetProposer(), msg.GetInitialDeposit(), msg.Validator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.GetProposer().String()),
		),
	)

	submitEvent := sdk.NewEvent(types.EventTypeSubmitProposal, sdk.NewAttribute(types.AttributeKeyProposalType, msg.GetContent().ProposalType()))
	if votingStarted {
		submitEvent = submitEvent.AppendAttributes(
			sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalId)),
		)
	}

	ctx.EventManager().EmitEvent(submitEvent)
	return &types.MsgSubmitProposalResponse{
		ProposalId: proposal.ProposalId,
	}, nil
}

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := getValidValidator(ctx, k.Keeper, msg.Depositor, msg.Validator); err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	err, votingStarted := k.Keeper.AddDeposit(ctx, msg.ProposalId, msg.Depositor, msg.Amount, msg.Validator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeProposalDeposit,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalId)),
			),
		)
	}

	return &types.MsgDepositResponse{}, nil
}

func (k msgServer) Vote(goCtx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := getValidValidator(ctx, k.Keeper, msg.Voter, msg.Validator); err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	err := k.Keeper.AddVote(ctx, msg.ProposalId, msg.Voter, msg.Option, msg.Validator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return &types.MsgVoteResponse{}, nil
}

//
// Internal methods
//

// checks if validator is active validator by signer and checks if incoming validator id matches with stored validator
func getValidValidator(ctx sdk.Context, keeper Keeper, signer sdk.AccAddress, validator hmTypes.ValidatorID) (hmTypes.Validator, error) {
	v, err := keeper.sk.GetActiveValidatorInfo(ctx, signer.Bytes())
	if err != nil {
		keeper.Logger(ctx).Info("No active validator by signer", "signer", signer.String())
		return v, err
	}

	// check if validator id matches with incoming validator and signer
	if v.ID.Uint64() != validator.Uint64() {
		keeper.Logger(ctx).Info(
			"Validator id mismatch",
			"expectedValidator", validator.String(),
			"storedValidator", v.ID.String(),
			"signer", signer.String(),
		)

		return v, fmt.Errorf(
			"Validator id mismatch. Expected validator %s, stored validator %s, signer %s",
			validator.String(),
			v.ID.String(),
			signer.String(),
		)
	}

	return v, nil
}
