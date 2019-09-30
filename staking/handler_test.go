package staking_test

import (
	"math/big"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/staking"
	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgValidatorJoin(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	mockVals := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
	// select first validator from slice
	mockVal := mockVals[0]
	t.Log("Inserting ===>", "Validator", mockVal.Signer.String())
	contractCallerObj.On("GetValidatorInfo", mock.Anything).Return(mockVal, nil)
	// insert new validator
	// msgTxHash := types.HeimdallHash("123")
	msgTxHash := types.HexToHeimdallHash("123")
	contractCallerObj.On("IsTxConfirmed", msgTxHash.EthHash()).Return(true)
	msgValJoin := staking.NewMsgValidatorJoin(mockVal.Signer, uint64(mockVal.ID), mockVal.PubKey, msgTxHash, 0)
	t.Log("msg val join", msgValJoin)
	got := staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator join to be ok, got %v", got)
	// validator is stored properly and signer is created properly
	storedVal, err := keeper.GetValidatorInfo(ctx, mockVal.Signer.Bytes())
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", mockVal.Signer.String(), err)
	require.Equal(t, mockVal.Signer, storedVal.Signer, "Signer address should match")
	t.Log("Stored ===>", "Validator", storedVal.String())
	// signer to validator mapping should exist properly
	storedSigner, found := keeper.GetSignerFromValidatorID(ctx, mockVal.ID)
	require.True(t, found, "signer and validator address should be mapped, got %v", found)
	require.Equal(t, mockVal.Signer.Bytes(), storedSigner.Bytes(), "Signer address in signer=>validator map should be same")
	t.Log("Mapped validator ID and Signer ===>", "ID", mockVal.ID, "Signer", storedSigner.String())
	// insert validator again
	got = staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
	require.True(t, !got.IsOK(), "expected validator join to be not-ok, got %v", got)
	// check if new validator gets added in validator set
}

func TestHandleMsgValidatorUpdate(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper, _ := cmn.CreateTestInput(t, false)

	// pass 0 as time alive to generate non de-activated validators
	LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	// vals := oldValSet.(*Validators)
	oldSigner := oldValSet.Validators[0]
	newSigner := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	t.Log("To be Updated ===>", "Validator", newSigner[0].String())
	// gen msg
	msgTxHash := types.HexToHeimdallHash("123")
	msg := staking.NewMsgSignerUpdate(newSigner[0].Signer, uint64(newSigner[0].ID), newSigner[0].PubKey, msgTxHash, 0)
	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	contractCallerObj.On("GetConfirmedTxReceipt", msgTxHash.EthHash()).Return(txreceipt, nil)
	signerUpdateEvent := &stakemanager.StakemanagerSignerChange{
		ValidatorId: new(big.Int).SetUint64(oldSigner.ID.Uint64()),
		OldSigner:   oldSigner.Signer.EthAddress(),
		NewSigner:   newSigner[0].Signer.EthAddress(),
	}
	contractCallerObj.On("DecodeSignerUpdateEvent", txreceipt, uint64(0)).Return(signerUpdateEvent, nil)

	got := staking.HandleMsgSignerUpdate(ctx, msg, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator update to be ok, got %v", got)
	newValidators := keeper.GetCurrentValidators(ctx)
	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")
	// apply updates
	// helper.UpdateValidators(
	// 	&oldValSet,                          // pointer to current validator set -- UpdateValidators will modify it
	// 	keeper.GetAllValidators(ctx),        // All validators
	// 	checkpointkeeper.GetACKCount(ctx)+1, // ack count
	// )
	setUpdates := helper.GetUpdatedValidators(&oldValSet, keeper.GetAllValidators(ctx), 5)
	oldValSet.UpdateWithChangeSet(setUpdates)
	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
	require.True(t, ok, "new signer should be found, got %v", ok)
	require.Equal(t, ValFrmID.Signer.Bytes(), newSigner[0].Signer.Bytes(), "New Signer should be mapped to old validator ID")
	require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)
	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.Signer.Bytes())
	require.Empty(t, err, "deleted validator should be found, got %v", err)
	require.Equal(t, removedVal.VotingPower, int64(0), "removed validator VotingPower should be zero")
	t.Log("Deleted validator ===>", "Validator", removedVal.String())

}

func TestHandleMsgValidatorExit(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper, checkpointkeeper := cmn.CreateTestInput(t, false)
	// pass 0 as time alive to generate non de-activated validators
	LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := types.HexToHeimdallHash("123")
	contractCallerObj.On("IsTxConfirmed", msgTxHash.EthHash()).Return(true)

	validators[0].EndEpoch = 10
	msg := staking.NewMsgValidatorExit(validators[0].Signer, uint64(validators[0].ID), msgTxHash, 0)
	contractCallerObj.On("GetValidatorInfo", validators[0].ID).Return(validators[0], nil)
	got := staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)
	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].Signer.String(), err)
	require.Equal(t, updatedValInfo.EndEpoch, validators[0].EndEpoch, "deactivation epoch should be set correctly")
	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")
	got = staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
	require.True(t, !got.IsOK(), "validator already exited. cannot exit again")

	currentVals := keeper.GetCurrentValidators(ctx)
	require.Equal(t, 4, len(currentVals), "No of current validators should exist before epoch passes")

	checkpointkeeper.UpdateACKCountWithValue(ctx, 20)
	currentVals = keeper.GetCurrentValidators(ctx)
	require.Equal(t, 3, len(currentVals), "No of current validators should reduce after epoch passes")
}

func TestHandleMsgStakeUpdate(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper, _ := cmn.CreateTestInput(t, false)

	// pass 0 as time alive to generate non de-activated validators
	LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	t.Log("To be Updated ===>", "Validator", oldVal.String())
	// gen msg
	msgTxHash := types.HexToHeimdallHash("123")
	msg := staking.NewMsgStakeUpdate(oldVal.Signer, oldVal.ID.Uint64(), msgTxHash, 0)
	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	contractCallerObj.On("GetConfirmedTxReceipt", msgTxHash.EthHash()).Return(txreceipt, nil)
	stakeUpdateEvent := &stakemanager.StakemanagerStakeUpdate{
		ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
		OldAmount:   new(big.Int).SetInt64(oldVal.VotingPower),
		NewAmount:   new(big.Int).SetInt64(2000000000000000000),
	}

	contractCallerObj.On("DecodeValidatorStakeUpdateEvent", txreceipt, uint64(0)).Return(stakeUpdateEvent, nil)

	got := staking.HandleMsgStakeUpdate(ctx, msg, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator stake update to be ok, got %v", got)
	updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.Signer.Bytes())
	require.Empty(t, err, "unable to fetch validator info %v-", err)
	require.Equal(t, stakeUpdateEvent.NewAmount.Int64(), updatedVal.VotingPower, "Validator VotingPower should be updated to %v", stakeUpdateEvent.NewAmount.Uint64())
}
