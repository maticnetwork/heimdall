package topup

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QuerySequence:
			return querySequence(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown topup query endpoint")
		}
	}
}

func querySequence(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QuerySequenceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf(err.Error()))
	}

	// get main tx receipt
	receipt, _ := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(params.TxHash).EthHash())
	if err != nil || receipt == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Transaction is not confirmed yet. Please for sometime and try again"))
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + params.LogIndex

	// check if incoming tx already exists
	if !k.HasTopupSequence(ctx, sequence) {
		k.Logger(ctx).Error("No sequence exist: %s %s", params.TxHash, params.LogIndex)
		return nil, sdk.ErrInternal(fmt.Sprintf("no sequence exist:: %s", params.TxHash))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sequence)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
