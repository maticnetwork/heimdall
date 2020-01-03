package checkpoint

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryAckCount:
			return handleQueryAckCount(ctx, req, keeper)
		case types.QueryInitialAccountRoot:
			return handleInitialAccountRoot(ctx, req, keeper)
		case types.QueryCheckpoint:
			return handleQueryCheckpoint(ctx, req, keeper)
		case types.QueryCheckpointBuffer:
			return handleQueryCheckpointBuffer(ctx, req, keeper)
		case types.QueryLastNoAck:
			return handleQueryLastNoAck(ctx, req, keeper)
		case types.QueryCheckpointList:
			return handleQueryCheckpointList(ctx, req, keeper)
		case types.QueryAccountProof:
			return handleQueryAccountProof(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handleQueryAckCount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetACKCount(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleInitialAccountRoot(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Calculate new account root hash
	dividendAccounts := keeper.sk.GetAllDividendAccounts(ctx)
	accountRoot, err := types.GetAccountRootHash(dividendAccounts)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch genesis accountroothash ", err.Error()))
	}
	return accountRoot, nil
}

func handleQueryAccountProof(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// 1. Fetch AccountRoot a1 present on RootChainContract
	// 2. Fetch AccountRoot a2 present in latest checkpoint
	// 3. if a1 == a2, Calculate merkle path using GetAllDividendAccounts
	// 4. if a1 != a2, Calculate merkle path using GetAllPrevDividendAccounts

	// contractCallerObj, err := helper.NewContractCaller()
	// if err != nil {

	// }

	// accountRootOnChain, err := contractCallerObj.CurrentAccountStateRoot()
	// if err != nil {
	// 	RestLogger.Error("Unable to get current account state root caller object ", "Error", err.Error())
	// }

	// lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	// var dividendAccounts []hmTypes.DividendAccount

	// if accountRootOnChain == lastCheckpoint.AccountRootHash {
	// 	dividendAccounts = keeper.sk.GetAllDividendAccounts(ctx)
	// } else {
	// 	dividendAccounts = keeper.sk.GetAllPrevDividendAccounts(ctx)
	// }
	// // Calculate new account root hash
	// merkleProof, index, err := types.GetAccountProof(dividendAccounts)
	// if err != nil {
	// 	return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch merkle proof ", err.Error()))
	// }
	// return merkleProof, index, err

	return nil, nil
}

func handleQueryCheckpoint(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryCheckpointParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetCheckpointByIndex(ctx, params.HeaderIndex)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch genesis rewardroothash", err.Error()))
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
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch genesis rewardroothash", err.Error()))
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
	var params types.QueryCheckpointListParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetCheckpointList(ctx, params.Page, params.Limit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch genesis rewardroothash", err.Error()))
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
