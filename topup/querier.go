package topup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper, contractCaller helper.IContractCaller) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QuerySequence:
			return querySequence(ctx, req, k, contractCaller)
		case types.QueryDividendAccount:
			return handleQueryDividendAccount(ctx, req, k)
		case types.QueryDividendAccountRoot:
			return handleDividendAccountRoot(ctx, req, k)
		case types.QueryAccountProof:
			return handleQueryAccountProof(ctx, req, k, contractCaller)
		case types.QueryVerifyAccountProof:
			return handleQueryVerifyAccountProof(ctx, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown topup query endpoint")
		}
	}
}

func querySequence(ctx sdk.Context, req abci.RequestQuery, k Keeper, contractCallerObj helper.IContractCaller) ([]byte, sdk.Error) {
	var params types.QuerySequenceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	chainParams := k.chainKeeper.GetParams(ctx)

	// get main tx receipt
	receipt, err := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(params.TxHash).EthHash(), chainParams.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Transaction is not confirmed yet. Please wait for sometime and try again"))
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(params.LogIndex))

	// check if incoming tx already exists
	if !k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("No sequence exist: %s %s", params.TxHash, params.LogIndex)
		return nil, nil
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sequence)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryDividendAccount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDividendAccountParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get dividend account info
	dividendAccount, err := keeper.GetDividendAccountByAddress(ctx, params.UserAddress)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("No dividend account found")
	}

	// json record
	bz, err := json.Marshal(dividendAccount)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleDividendAccountRoot(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// Calculate new account root hash
	dividendAccounts := keeper.GetAllDividendAccounts(ctx)
	accountRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch accountroothash ", err.Error()))
	}
	return accountRoot, nil
}

func handleQueryAccountProof(ctx sdk.Context, req abci.RequestQuery, keeper Keeper, contractCallerObj helper.IContractCaller) ([]byte, sdk.Error) {
	// 1. Fetch AccountRoot a1 present on RootChainContract
	// 2. Fetch AccountRoot a2 from current account
	// 3. if a1 == a2, Calculate merkle path using GetAllDividendAccounts

	var params types.QueryAccountProofParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	chainParams := keeper.chainKeeper.GetParams(ctx)

	stakingInfoAddress := chainParams.ChainParams.StakingInfoAddress.EthAddress()
	stakingInfoInstance, _ := contractCallerObj.GetStakingInfoInstance(stakingInfoAddress)

	accountRootOnChain, err := contractCallerObj.CurrentAccountStateRoot(stakingInfoInstance)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch account root from onchain ", err.Error()))
	}

	dividendAccounts := keeper.GetAllDividendAccounts(ctx)
	currentStateAccountRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)

	if bytes.Equal(accountRootOnChain[:], currentStateAccountRoot) {
		// Calculate new account root hash
		merkleProof, index, err := checkpointTypes.GetAccountProof(dividendAccounts, params.UserAddress)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could fetch account proof", err.Error()))
		}

		accountProof := hmTypes.NewDividendAccountProof(params.UserAddress, merkleProof, index)

		// json record
		bz, err := json.Marshal(accountProof)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil

	} else {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch merkle proof ", err.Error()))
	}
}

func handleQueryVerifyAccountProof(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryVerifyAccountProofParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	dividendAccounts := keeper.GetAllDividendAccounts(ctx)

	// Verify account proof
	accountProofStatus, err := checkpointTypes.VerifyAccountProof(dividendAccounts, params.UserAddress, params.AccountProof)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not verify merkle proof ", err.Error()))
	}

	// json record
	bz, err := json.Marshal(accountProofStatus)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
