package test

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/maticnetwork/bor/common"
// 	"github.com/maticnetwork/heimdall/helper"
// 	"github.com/maticnetwork/heimdall/helper/mocks"
// 	"github.com/maticnetwork/heimdall/staking"
// 	"github.com/stretchr/testify/mock"
// 	"github.com/stretchr/testify/require"
// )

// func TestHandleMsgValidatorUpdate(t *testing.T) {
// 	contractCallerObj := mocks.IContractCaller{}
// 	ctx, keeper := CreateTestInput(t, false)
// 	// pass 0 as time alive to generate non de-activated validators
// 	LoadValidatorSet(4, t, keeper, ctx, false, 0)
// 	oldValSet := keeper.GetValidatorSet(ctx)
// 	PrintVals(t, &oldValSet)
// 	// vals := oldValSet.(*Validators)
// 	oldSigner := oldValSet.Validators[0]
// 	newSigner := GenRandomVal(1, 0, 10, 10, false, 1)
// 	newSigner[0].ID = oldSigner.ID
// 	newSigner[0].Power = oldSigner.Power
// 	t.Log("To be Updated ===>", "Validator", newSigner[0].String())
// 	// gen msg
// 	msgTxHash := common.HexToHash("123")
// 	msg := staking.NewMsgSignerUpdate(uint64(newSigner[0].ID), newSigner[0].PubKey, json.Number("10"), msgTxHash)
// 	contractCallerObj.On("IsTxConfirmed", msgTxHash).Return(true)
// 	contractCallerObj.On("DecodeSignerUpdateEvent", msgTxHash).Return(uint64(oldSigner.ID), newSigner[0].PubKey.Address(), oldSigner.Signer, nil)
// 	got := staking.HandleMsgSignerUpdate(ctx, msg, keeper, &contractCallerObj)
// 	require.True(t, got.IsOK(), "expected validator update to be ok, got %v", got)
// 	newValidators := keeper.GetCurrentValidators(ctx)
// 	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of validators should be equal")
// 	// apply updates
// 	helper.UpdateValidators(
// 		&oldValSet,                   // pointer to current validator set -- UpdateValidators will modify it
// 		keeper.GetAllValidators(ctx), // All validators
// 		keeper.GetACKCount(ctx)+1,    // ack count
// 	)
// 	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)
// 	newValSet := keeper.GetValidatorSet(ctx)
// 	PrintVals(t, &newValSet)
// 	_, mutatedVal := newValSet.GetByAddress(newSigner[0].PubKey.Address().Bytes())
// 	require.Equal(t, oldSigner.Accum, (*mutatedVal).Accum, "Accum should remain the same")
// 	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
// 	require.True(t, ok, "new signer should be found, got %v", ok)
// 	require.Equal(t, ValFrmID.Signer.Bytes(), newSigner[0].Signer.Bytes(), "New Signer should be mapped to old validator ID")
// 	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.Signer.Bytes())
// 	require.Empty(t, err, "deleted validator should be found, got %v", err)
// 	t.Log("Deleted validator ===>", "Validator", removedVal.String())
// }

// func TestHandleMsgValidatorExit(t *testing.T) {
// 	contractCallerObj := mocks.IContractCaller{}
// 	ctx, keeper := CreateTestInput(t, false)
// 	// pass 0 as time alive to generate non de-activated validators
// 	LoadValidatorSet(4, t, keeper, ctx, false, 0)
// 	validators := keeper.GetCurrentValidators(ctx)

// 	msgTxHash := common.HexToHash("123")
// 	contractCallerObj.On("IsTxConfirmed", msgTxHash).Return(true)
// 	validators[0].EndEpoch = 10
// 	msg := staking.NewMsgValidatorExit(uint64(validators[0].ID), msgTxHash)
// 	contractCallerObj.On("GetValidatorInfo", validators[0].ID).Return(validators[0], nil)
// 	got := staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
// 	require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)
// 	got = staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
// 	require.True(t, !got.IsOK(), "expected validator exit to be ok, got %v", got)
// 	validator, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
// 	require.True(t, found, "Validator should be present even after deactivation")
// 	require.Equal(t, 10, int(validator.EndEpoch), "end epoch should be set to 10")

// 	keeper.UpdateACKCountWithValue(ctx, 20)
// 	currentVals := keeper.GetCurrentValidators(ctx)
// 	require.Equal(t, 3, len(currentVals), "No validators should exist after epoch passes")

// 	found = FindSigner(validators[0].Signer, currentVals)
// 	require.True(t, !found, "Validator should not exist in current val set")
// }

// func TestHandleMsgValidatorJoin(t *testing.T) {
// 	contractCallerObj := mocks.IContractCaller{}
// 	ctx, keeper := CreateTestInput(t, false)
// 	mockVals := GenRandomVal(1, 0, 10, 10, false, 1)
// 	// select first validator from slice
// 	mockVal := mockVals[0]
// 	t.Log("Inserting ===>", "Validator", mockVal.Signer.String())
// 	contractCallerObj.On("GetValidatorInfo", mock.Anything).Return(mockVal, nil)
// 	// insert new validator
// 	msgTxHash := common.HexToHash("123")
// 	contractCallerObj.On("IsTxConfirmed", msgTxHash).Return(true)
// 	msgValJoin := staking.NewMsgValidatorJoin(uint64(mockVal.ID), mockVal.PubKey, msgTxHash)
// 	got := staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
// 	require.True(t, got.IsOK(), "expected validator join to be ok, got %v", got)
// 	// validator is stored properly and signer is created properly
// 	storedVal, err := keeper.GetValidatorInfo(ctx, mockVal.Signer.Bytes())
// 	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", mockVal.Signer.String(), err)
// 	require.Equal(t, mockVal.PubKey.Address(), storedVal.Signer, "Signer address should match")
// 	t.Log("Stored ===>", "Validator", storedVal.String())
// 	// signer to validator mapping should exist properly
// 	storedSigner, found := keeper.GetSignerFromValidatorID(ctx, mockVal.ID)
// 	require.True(t, found, "signer and validator address should be mapped, got %v", found)
// 	require.Equal(t, mockVal.Signer.Bytes(), storedSigner.Bytes(), "Signer address in signer=>validator map should be same")
// 	t.Log("Mapped validator ID and Signer ===>", "ID", mockVal.ID, "Signer", storedSigner.String())
// 	// insert validator again
// 	got = staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
// 	require.True(t, !got.IsOK(), "expected validator join to be not-ok, got %v", got)
// 	// check if new validator gets added in validator set
// }
