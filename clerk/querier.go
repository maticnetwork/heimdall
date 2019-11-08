package clerk

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(clerkKeeper Keeper, stakingKeeper staking.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case clerkTypes.QueryRecord:
			return HandlerQueryRecord(ctx, req, clerkKeeper)
		case clerkTypes.QueryStateSyncer:
			return HandlerQueryStateSyncer(ctx, req, clerkKeeper, stakingKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func HandlerQueryRecord(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params clerkTypes.QueryRecordParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get state record by record id
	record, err := keeper.GetEventRecord(ctx, params.RecordID)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get state record", err.Error()))
	}

	// return error if record doesn't exist
	if record == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("record %v does not exist", params.RecordID))
	}

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, record)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

// HandlerQueryStateSyncer will select next 3 state syncers
func HandlerQueryStateSyncer(ctx sdk.Context, req abci.RequestQuery, clerkKeeper Keeper, stakingKeeper staking.Keeper) ([]byte, sdk.Error) {

	// Total no of state sync events
	eventCount := clerkKeeper.GetStateSyncEventCount(ctx)

	// Total no of validators
	validatorCount := len(stakingKeeper.GetAllValidators(ctx))
	stateSyncerList := []types.Validator{}

	// No of state syncers to return
	syncerCount := 3

	// if no of validators is less than three, return existing validators
	if validatorCount < 3 {
		syncerCount = validatorCount
	}

	// Select next 3 active validators
	for i := 0; i < syncerCount; {
		valIndex := eventCount % uint64(validatorCount)
		validator, _ := stakingKeeper.GetValidatorFromValID(ctx, types.ValidatorID(valIndex+1))
		if stakingKeeper.IsCurrentValidatorByAddress(ctx, validator.Signer.Bytes()) {
			stateSyncerList = append(stateSyncerList, validator)
			i++
			eventCount++
		}
	}

	res, err := codec.MarshalJSONIndent(clerkKeeper.cdc, stateSyncerList)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal syncer list to JSON", err.Error()))
	}

	return res, nil
}
