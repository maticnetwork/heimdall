package staking

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewQuerier returns querier for staking Rest endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryCurrentValidatorSet:
			return handleQueryCurrentValidatorSet(ctx, req, keeper)
		case types.QuerySigner:
			return handleQuerySigner(ctx, req, keeper)
		case types.QueryValidator:
			return handleQueryValidator(ctx, req, keeper)
		case types.QueryValidatorStatus:
			return handleQueryValidatorStatus(ctx, req, keeper)
		case types.QueryProposer:
			return handleQueryProposer(ctx, req, keeper)
		case types.QueryCurrentProposer:
			return handleQueryCurrentProposer(ctx, req, keeper)
		case types.QueryProposerBonusPercent:
			return handleQueryProposerBonusPercent(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func handleQueryCurrentValidatorSet(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// get validator set
	validatorSet := keeper.GetValidatorSet(ctx)

	// json record
	bz, err := json.Marshal(validatorSet)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQuerySigner(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySignerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator info
	validator, err := keeper.GetValidatorInfo(ctx, params.SignerAddress)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Error while getting validator by signer", err.Error()))
	}

	// json record
	bz, err := json.Marshal(validator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryValidator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator info
	validator, ok := keeper.GetValidatorFromValID(ctx, params.ValidatorID)
	if !ok {
		return nil, sdk.ErrUnknownRequest("No validator found")
	}

	// json record
	bz, err := json.Marshal(validator)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryValidatorStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySignerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator status by signer address
	status := keeper.IsCurrentValidatorByAddress(ctx, params.SignerAddress)

	// json record
	bz, err := json.Marshal(status)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryProposer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryProposerParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator set
	validatorSet := keeper.GetValidatorSet(ctx)

	times := int(params.Times)
	if times > len(validatorSet.Validators) {
		times = len(validatorSet.Validators)
	}

	// init proposers
	var proposers []hmTypes.Validator

	// get proposers
	for index := 0; index < times; index++ {
		proposers = append(proposers, *(validatorSet.GetProposer()))
		validatorSet.IncrementProposerPriority(1)
	}

	// json record
	bz, err := json.Marshal(proposers)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryCurrentProposer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	proposer := keeper.GetCurrentProposer(ctx)

	bz, err := json.Marshal(proposer)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryProposerBonusPercent(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// GetProposerBonusPercent
	proposerBonusPercent := keeper.GetProposerBonusPercent(ctx)

	// json record
	bz, err := json.Marshal(proposerBonusPercent)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
