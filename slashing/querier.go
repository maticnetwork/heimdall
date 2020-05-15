package slashing

import (
	"encoding/json"
	"fmt"
	"math/big"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewQuerier creates a new querier for slashing clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, k)

		case types.QueryTickCount:
			return handleQueryTickCount(ctx, req, k)

		case types.QuerySigningInfo:
			return querySigningInfo(ctx, req, k)

		case types.QuerySigningInfos:
			return querySigningInfos(ctx, req, k)

		case types.QuerySlashingInfo:
			return querySlashingInfo(ctx, req, k)

		case types.QuerySlashingInfos:
			return querySlashingInfos(ctx, req, k)

		case types.QuerySlashingInfoBytes:
			return querySlashingInfoBytes(ctx, req, k)

		case types.QueryTickSlashingInfos:
			return queryTickSlashingInfos(ctx, req, k)

		case types.QuerySlashingSequence:
			return querySlashingSequence(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown slashing query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySigningInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySigningInfoParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator signing info
	signingInfo, found := k.GetValidatorSigningInfo(ctx, params.ValidatorID)
	if !found {
		return nil, sdk.ErrInternal("Error while getting validator signing info")
	}

	// json record
	bz, err := json.Marshal(signingInfo)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySigningInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySigningInfosParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	var signingInfos []hmTypes.ValidatorSigningInfo

	k.IterateValidatorSigningInfos(ctx, func(valID hmTypes.ValidatorID, info hmTypes.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})

	start, end := client.Paginate(len(signingInfos), params.Page, params.Limit, len(signingInfos))
	if start < 0 || end < 0 {
		signingInfos = []hmTypes.ValidatorSigningInfo{}
	} else {
		signingInfos = signingInfos[start:end]
	}

	// json record
	bz, err := json.Marshal(signingInfos)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryTickCount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetTickCount(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySlashingInfo(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySlashingInfoParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get validator slashing info
	slashingInfo, found := k.GetBufferValSlashingInfo(ctx, params.ValidatorID)
	if !found {
		return nil, sdk.ErrInternal(" slashing info not found for given val")
	}

	// json record
	bz, err := json.Marshal(slashingInfo)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySlashingInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySlashingInfosParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	var slashingInfos []hmTypes.ValidatorSlashingInfo

	k.IterateBufferValSlashingInfos(ctx, func(info hmTypes.ValidatorSlashingInfo) (stop bool) {
		slashingInfos = append(slashingInfos, info)
		return false
	})

	start, end := client.Paginate(len(slashingInfos), params.Page, params.Limit, len(slashingInfos))
	if start < 0 || end < 0 {
		slashingInfos = []hmTypes.ValidatorSlashingInfo{}
	} else {
		slashingInfos = slashingInfos[start:end]
	}

	// json record
	bz, err := json.Marshal(slashingInfos)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySlashingInfoBytes(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Calculate new slashInfo bytes
	slashingInfos, err := keeper.GetBufferValSlashingInfos(ctx)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("no slash infos in buffer", err.Error()))
	}

	slashingInfoBytes, err := types.SortAndRLPEncodeSlashInfos(slashingInfos)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch slashingInfoBytes ", err.Error()))
	}
	return slashingInfoBytes, nil
}

func queryTickSlashingInfos(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryTickSlashingInfosParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	var slashingInfos []hmTypes.ValidatorSlashingInfo

	k.IterateTickValSlashingInfos(ctx, func(info hmTypes.ValidatorSlashingInfo) (stop bool) {
		slashingInfos = append(slashingInfos, info)
		return false
	})

	start, end := client.Paginate(len(slashingInfos), params.Page, params.Limit, len(slashingInfos))
	if start < 0 || end < 0 {
		slashingInfos = []hmTypes.ValidatorSlashingInfo{}
	} else {
		slashingInfos = slashingInfos[start:end]
	}

	// json record
	bz, err := json.Marshal(slashingInfos)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func querySlashingSequence(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySlashingSequenceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	chainParams := keeper.chainKeeper.GetParams(ctx)

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf(err.Error()))
	}

	// get main tx receipt
	receipt, err := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(params.TxHash).EthHash(), chainParams.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Transaction is not confirmed yet. Please wait for sometime and try again"))
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(params.LogIndex))

	// check if incoming tx already exists
	if !keeper.HasSlashingSequence(ctx, sequence.String()) {
		keeper.Logger(ctx).Error("No slashing sequence exist: %s %s", params.TxHash, params.LogIndex)
		return nil, nil
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sequence)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
