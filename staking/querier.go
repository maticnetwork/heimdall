package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
)

// NewQuerier returns querier for staking Rest endpoints
func NewQuerier(keeper Keeper) sdk.Querier {

	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case stakingTypes.QueryValStatus:
			return handlerQueryValStatus(ctx, req, keeper)
		case stakingTypes.QueryCheckpointReward:
			return handlerQueryCheckpointReward(ctx, req, keeper)
		case stakingTypes.QueryProposerBonusPercent:
			return handlerQueryProposerBonusPercent(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handlerQueryValStatus(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params stakingTypes.QueryValStatusParams
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

func handlerQueryCheckpointReward(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// GetCheckpointReward
	checkpointReward := keeper.GetCheckpointReward(ctx)

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, checkpointReward)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handlerQueryProposerBonusPercent(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// GetProposerBonusPercent
	proposerBonusPercent := keeper.GetProposerBonusPercent(ctx)

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, proposerBonusPercent)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
