package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/gov/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler handle all "gov" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgDeposit:
			return handleMsgDeposit(ctx, keeper, msg)

		case types.MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)

		case types.MsgVote:
			return handleMsgVote(ctx, keeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized gov message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg types.MsgSubmitProposal) sdk.Result {
	if _, err := getValidValidator(ctx, keeper, msg.Proposer, msg.Validator); err != nil {
		return hmCommon.ErrInvalidMsg(keeper.Codespace(), "No active validator by proposer").Result()
	}

	proposal, err := keeper.SubmitProposal(ctx, msg.Content)
	if err != nil {
		return err.Result()
	}

	err, votingStarted := keeper.AddDeposit(ctx, proposal.ProposalID, msg.Proposer, msg.InitialDeposit, msg.Validator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer.String()),
		),
	)

	if votingStarted {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSubmitProposal,
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", proposal.ProposalID)),
			),
		)
	}

	return sdk.Result{
		Data:   keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}
}

func handleMsgDeposit(ctx sdk.Context, keeper Keeper, msg types.MsgDeposit) sdk.Result {
	if _, err := getValidValidator(ctx, keeper, msg.Depositor, msg.Validator); err != nil {
		return hmCommon.ErrInvalidMsg(keeper.Codespace(), "No active validator by depositor").Result()
	}

	err, votingStarted := keeper.AddDeposit(ctx, msg.ProposalID, msg.Depositor, msg.Amount, msg.Validator)
	if err != nil {
		return err.Result()
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
				sdk.NewAttribute(types.AttributeKeyVotingPeriodStart, fmt.Sprintf("%d", msg.ProposalID)),
			),
		)
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg types.MsgVote) sdk.Result {
	if _, err := getValidValidator(ctx, keeper, msg.Voter, msg.Validator); err != nil {
		return hmCommon.ErrInvalidMsg(keeper.Codespace(), "No active validator by voter").Result()
	}

	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option, msg.Validator)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

//
// Internal methods
//

// checks if validator is active validator by signer and checks if incoming validator id matches with stored validator
func getValidValidator(ctx sdk.Context, keeper Keeper, signer hmTypes.HeimdallAddress, validator hmTypes.ValidatorID) (hmTypes.Validator, error) {
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
