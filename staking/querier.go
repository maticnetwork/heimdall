package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/staking/types"
)

// NewQuerier returns querier for staking Rest endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryValidatorStatus:
			return handleQueryValidatorStatus(ctx, req, keeper)
		case types.QueryProposerBonusPercent:
			return handleQueryProposerBonusPercent(ctx, req, keeper)
		case types.QueryCurrentValidatorSet:
			return handleQueryCurrentValidatorSet(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handleQueryValidatorStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryValidatorStatusParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator status by signer address
	status := keeper.IsCurrentValidatorByAddress(ctx, params.SignerAddress)

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, status)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryProposerBonusPercent(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// GetProposerBonusPercent
	proposerBonusPercent := keeper.GetProposerBonusPercent(ctx)

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposerBonusPercent)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryCurrentValidatorSet(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// get validator set
	validatorSet := keeper.GetValidatorSet(ctx)

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, validatorSet)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
