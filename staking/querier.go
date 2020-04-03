package staking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
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
		case types.QueryDividendAccount:
			return handleQueryDividendAccount(ctx, req, keeper)
		case types.QueryDividendAccountRoot:
			return handleDividendAccountRoot(ctx, req, keeper)
		case types.QueryAccountProof:
			return handleQueryAccountProof(ctx, req, keeper)
		case types.QueryVerifyAccountProof:
			return handleQueryVerifyAccountProof(ctx, req, keeper)
		case types.QueryStakingSequence:
			return handleQueryStakingSequence(ctx, req, keeper)

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

func handleQueryDividendAccount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryDividendAccountParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	// get dividend account info
	dividendAccount, err := keeper.GetDividendAccountByID(ctx, params.DividendAccountID)
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

func handleQueryAccountProof(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	// 1. Fetch AccountRoot a1 present on RootChainContract
	// 2. Fetch AccountRoot a2 from current account
	// 3. if a1 == a2, Calculate merkle path using GetAllDividendAccounts

	var params types.QueryAccountProofParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	contractCallerObj, err := helper.NewContractCaller()
	accountRootOnChain, err := contractCallerObj.CurrentAccountStateRoot()
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch account root from onchain ", err.Error()))
	}

	dividendAccounts := keeper.GetAllDividendAccounts(ctx)
	currentStateAccountRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)

	if bytes.Compare(accountRootOnChain[:], currentStateAccountRoot) == 0 {
		// Calculate new account root hash
		merkleProof, index, _ := checkpointTypes.GetAccountProof(dividendAccounts, params.DividendAccountID)
		accountProof := hmTypes.NewDividendAccountProof(params.DividendAccountID, merkleProof, index)
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
	accountProofStatus, err := checkpointTypes.VerifyAccountProof(dividendAccounts, params.DividendAccountID, params.AccountProof)
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

func handleQueryStakingSequence(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryStakingSequenceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	contractCallerObj, err := helper.NewContractCaller()
	if err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf(err.Error()))
	}

	// get main tx receipt
	receipt, _ := contractCallerObj.GetConfirmedTxReceipt(time.Now().UTC(), hmTypes.HexToHeimdallHash(params.TxHash).EthHash())
	if err != nil || receipt == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("Transaction is not confirmed yet. Please for sometime and try again"))
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(params.LogIndex))

	// check if incoming tx already exists
	if !keeper.HasStakingSequence(ctx, sequence.String()) {
		keeper.Logger(ctx).Error("No staking sequence exist: %s %s", params.TxHash, params.LogIndex)
		return nil, sdk.ErrInternal(fmt.Sprintf("no sequence exist:: %s", params.TxHash))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, sequence)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
