package checkpoint

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/topup"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper, stakingKeeper staking.Keeper, topupKeeper topup.Keeper, contractCaller helper.IContractCaller) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return handleQueryParams(ctx, req, keeper)
		case types.QueryAckCount:
			return handleQueryAckCount(ctx, req, keeper)
		case types.QueryCheckpoint:
			return handleQueryCheckpoint(ctx, req, keeper)
		case types.QueryCheckpointBuffer:
			return handleQueryCheckpointBuffer(ctx, req, keeper)
		case types.QueryLastNoAck:
			return handleQueryLastNoAck(ctx, req, keeper)
		case types.QueryCheckpointList:
			return handleQueryCheckpointList(ctx, req, keeper)
		case types.QueryNextCheckpoint:
			return handleQueryNextCheckpoint(ctx, req, keeper, stakingKeeper, topupKeeper, contractCaller)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handleQueryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryAckCount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetACKCount(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryCheckpoint(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryCheckpointParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetCheckpointByNumber(ctx, params.Number)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch checkpoint by index %v", params.Number), err.Error()))
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryCheckpointBuffer(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res, err := keeper.GetCheckpointFromBuffer(ctx)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch checkpoint buffer", err.Error()))
	}

	if res == nil {
		return nil, common.ErrNoCheckpointBufferFound(keeper.Codespace())
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryLastNoAck(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// get last no ack
	res := keeper.GetLastNoAck(ctx)
	// sed result
	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryCheckpointList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params hmTypes.QueryPaginationParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetCheckpointList(ctx, params.Page, params.Limit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch checkpoint list with page %v and limit %v", params.Page, params.Limit), err.Error()))
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryNextCheckpoint(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, sk staking.Keeper, tk topup.Keeper, contractCaller helper.IContractCaller) ([]byte, sdk.Error) {
	var queryParams types.QueryBorChainID
	if err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse query params: %s", err))
	}

	// get validator set
	validatorSet := sk.GetValidatorSet(ctx)
	proposer := validatorSet.GetProposer()
	ackCount := keeper.GetACKCount(ctx)
	params := keeper.GetParams(ctx)

	var start uint64

	if ackCount != 0 {
		checkpointNumber := ackCount
		lastCheckpoint, err := keeper.GetCheckpointByNumber(ctx, checkpointNumber)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch checkpoint by index %v", checkpointNumber), err.Error()))
		}
		start = lastCheckpoint.EndBlock + 1
	}

	end := start + params.AvgCheckpointLength

	rootHash, err := contractCaller.GetRootHash(start, end, params.MaxCheckpointLength)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch roothash for start:%v end:%v error:%v", start, end, err), err.Error()))
	}

	accs := tk.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(accs)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not get generate account root hash. Error:%v", err), err.Error()))
	}

	checkpointMsg := types.NewMsgCheckpointBlock(
		proposer.Signer,
		start,
		start+params.AvgCheckpointLength,
		hmTypes.BytesToHeimdallHash(rootHash),
		hmTypes.BytesToHeimdallHash(accRootHash),
		queryParams.BorChainID,
	)
	bz, err := json.Marshal(checkpointMsg)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not marshall checkpoint msg. Error:%v", err), err.Error()))
	}
	return bz, nil
}
