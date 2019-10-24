package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/auth/types"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the staking Querier
const (
	QuerySlashValidator = "slash-validator"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QuerySlashValidator:
			return querySlashValidator(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

func querySlashValidator(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params stakingTypes.ValidatorSlashParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	err := keeper.SlashValidator(ctx, params.ValID, params.SlashAmount)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not Slash validator", err.Error()))
	}
	return nil, nil
}
