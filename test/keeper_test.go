package test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/staking"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"time"

	"github.com/maticnetwork/heimdall/checkpoint"
	hmcmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/types"
)

// TODO use table testing as much as possible

func TestUpdateAck(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.UpdateACKCount(ctx)
	ack := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(2), ack, "Ack Count Not Equal")
}

func TestCheckpointACK(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	prevACK := keeper.GetACKCount(ctx)

	// create random header block
	headerBlock, err := GenRandCheckpointHeader(265)
	require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	keeper.AddCheckpoint(ctx, 20000, headerBlock)
	require.Empty(t, err, "Unable to store checkpoint, Error: %v", err)

	keeper.UpdateACKCount(ctx)
	keeper.FlushCheckpointBuffer(ctx)
	acksCount := keeper.GetACKCount(ctx)

	// fetch last checkpoint key (NumberOfACKs * ChildBlockInterval)
	lastCheckpointKey := 10000 * acksCount

	storedHeader, err := keeper.GetCheckpointByIndex(ctx, lastCheckpointKey)
	// TODO uncomment when config is loading properly
	//storedHeader, err := keeper.GetLastCheckpoint(ctx)
	require.Empty(t, err, "Unable to retrieve checkpoint, Error: %v", err)
	require.Equal(t, headerBlock, storedHeader, "Header Blocks dont match")

	currentACK := keeper.GetACKCount(ctx)
	require.Equal(t, prevACK+1, currentACK, "ACK count should have been incremented by 1")

	// flush and check if its flushed
	keeper.FlushCheckpointBuffer(ctx)
	storedHeader, err = keeper.GetCheckpointFromBuffer(ctx)
	require.NotEmpty(t, err, "HeaderBlock should not exist after flush")

}

// tests setter/getters for validatorSignerMaps , validator set/get
func TestValidator(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	vals := GenRandomVal(1, 0, 0, 10, true, 1)
	validator := vals[0]

	err := keeper.AddValidator(ctx, validator)
	require.Empty(t, err, "Unable to set validator, Error: %v", err)

	storedVal, err := keeper.GetValidatorInfo(ctx, validator.Signer.Bytes())
	require.Empty(t, err, "Unable to fetch validator")
	require.Equal(t, validator, storedVal, "Unable to fetch validator from val address")

	storedSigner, ok := keeper.GetSignerFromValidatorID(ctx, validator.ID)
	require.Equal(t, true, ok, "Validator<=>Signer not mapped")
	require.Equal(t, validator.Signer, storedSigner, "Signer doesnt match")

	storedValidator, ok := keeper.GetValidatorFromValID(ctx, validator.ID)
	require.Equal(t, true, ok, "Validator<=>Signer not mapped")
	require.Equal(t, validator, storedValidator, "Unable to fetch validator from val address")

	//valToSignerMap := keeper.GetSignerFromValidatorID(ctx)
	//mappedSigner := valToSignerMap[hex.EncodeToString(validator.Signer.Bytes())]
	//require.Equal(t, validator.Signer, mappedSigner, "GetValidatorToSignerMap doesnt give right signer")
}

func TestValidatorSet(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	valSet := LoadValidatorSet(4, t, keeper, ctx, true, 10)

	storedValSet := keeper.GetValidatorSet(ctx)
	require.Equal(t, valSet, storedValSet, "Validator Set in state doesnt match ")

	//keeper.IncreamentAccum(ctx, 1)
	//initialProposer := keeper.GetCurrentProposer(ctx)
	//
	//keeper.IncreamentAccum(ctx, 1)
	//newProposer := keeper.GetCurrentProposer(ctx)
	//fmt.Printf("Prev :%#v  , New : %#v", initialProposer, newProposer)
}

func TestValUpdates(t *testing.T) {

	// create sub test to check if validator remove
	t.Run("remove", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)

		// load 4 validators to state
		LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)

		currentValSet := initValSet.Copy()
		prevValidatorSet := initValSet.Copy()

		// remove validator (making IsCurrentValidator return false)
		prevValidatorSet.Validators[0].StartEpoch = 20

		t.Log("Updated Validators in state")
		for _, v := range prevValidatorSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}

		err := keeper.AddValidator(ctx, *prevValidatorSet.Validators[0])
		require.Empty(t, err, "Unable to update validator set")

		// apply updates
		helper.UpdateValidators(
			currentValSet,                // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx), // All validators
			5, // ack count
		)
		updatedValSet := currentValSet
		t.Log("Validators in updated validator set")
		for _, v := range updatedValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}
		// check if 1 validator is removed
		require.Equal(t, len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators), "Validator set should be reduced by one ")
		// remove first validator from initial validator set and equate with new
		require.Equal(t, append(prevValidatorSet.Validators[:0], prevValidatorSet.Validators[1:]...), updatedValSet.Validators, "Validator at 0 index should be deleted")
	})

	t.Run("add", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)

		// load 4 validators to state
		LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)

		validators := GenRandomVal(1, 0, 10, 10, false, 1)
		prevValSet := initValSet.Copy()
		valToBeAdded := validators[0]
		currentValSet := initValSet.Copy()
		//prevValidatorSet := initValSet.Copy()
		keeper.AddValidator(ctx, valToBeAdded)

		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}
		t.Log("Val to be Added")
		t.Log("-->", "Address", valToBeAdded.Signer.String(), "StartEpoch", valToBeAdded.StartEpoch, "EndEpoch", valToBeAdded.EndEpoch, "Power", valToBeAdded.Power)

		helper.UpdateValidators(
			currentValSet,                // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx), // All validators
			5, // ack count
		)
		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}

		require.Equal(t, len(prevValSet.Validators)+1, len(currentValSet.Validators), "Number of validators should be increased by 1")
		require.Equal(t, true, currentValSet.HasAddress(valToBeAdded.Signer.Bytes()), "New Validator should be added")
		require.Equal(t, prevValSet.TotalVotingPower()+int64(valToBeAdded.Power), currentValSet.TotalVotingPower(), "Total power should be increased")
	})

	t.Run("update", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)

		// load 4 validators to state
		LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)
		keeper.IncreamentAccum(ctx, 2)
		prevValSet := initValSet.Copy()
		currentValSet := keeper.GetValidatorSet(ctx)
		valToUpdate := currentValSet.Validators[0]
		newSigner := GenRandomVal(1, 0, 10, 10, false, 1)
		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "Accum", v.Accum, "Signer", v.Signer.String(), "Total power", currentValSet.TotalVotingPower())
		}
		keeper.UpdateSigner(ctx, newSigner[0].Signer, newSigner[0].PubKey, valToUpdate.Signer)
		helper.UpdateValidators(
			&currentValSet,               // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx), // All validators
			5, // ack count
		)
		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "Accum", v.Accum, "Signer", v.Signer.String(), "Total power", currentValSet.TotalVotingPower())
		}

		require.Equal(t, len(prevValSet.Validators), len(currentValSet.Validators), "Number of validators should remain same")
		_, val := currentValSet.GetByAddress(valToUpdate.Signer.Bytes())
		require.Equal(t, newSigner[0].Signer, val.Signer, "Signer address should change")
		require.Equal(t, newSigner[0].PubKey, val.PubKey, "Signer pubkey should change")
		require.Equal(t, valToUpdate.Accum, val.Accum, "Validator accum should not change")
		require.Equal(t, prevValSet.TotalVotingPower(), currentValSet.TotalVotingPower(), "Total power should not change")
		// TODO not sure if proposer check is needed
		//require.Equal(t, &initValSet.Proposer.Address, &currentValSet.Proposer.Address, "Proposer should not change")
	})

}

func LoadValidatorSet(count int, t *testing.T, keeper hmcmn.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	// create 4 validators
	validators := GenRandomVal(4, 0, 10, uint64(timeAlive), randomise, 1)

	var valSet types.ValidatorSet

	// add validators to new Validator set and state
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.Empty(t, err, "Unable to set validator, Error: %v", err)
		// add validator to validator set
		valSet.Add(&validator)
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")
	return valSet
}

// test handler for message
func TestHandleMsgCheckpoint(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}

	// check valid checkpoint
	t.Run("validCheckpoint", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)
		// generate proposer for validator set
		LoadValidatorSet(4, t, keeper, ctx, false, 10)
		keeper.IncreamentAccum(ctx, 1)
		header, err := GenRandCheckpointHeader(10)
		require.Empty(t, err, "Unable to create random header block, Error:%v", err)
		// make sure proposer has min ether
		contractCallerObj.On("GetBalance", keeper.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)
		SentValidCheckpoint(header, keeper, ctx, contractCallerObj, t)
	})

	// check invalid proposer
	t.Run("invalidProposer", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)
		// generate proposer for validator set
		LoadValidatorSet(4, t, keeper, ctx, false, 10)
		keeper.IncreamentAccum(ctx, 1)
		header, err := GenRandCheckpointHeader(10)
		require.Empty(t, err, "Unable to create random header block, Error:%v", err)

		// add wrong proposer to header
		header.Proposer = keeper.GetValidatorSet(ctx).Validators[2].Signer
		// make sure proposer has min ether
		contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)

		// create checkpoint msg
		msgCheckpoint := checkpoint.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, uint64(time.Now().Unix()))

		// send checkpoint to handler
		got := checkpoint.HandleMsgCheckpoint(ctx, msgCheckpoint, keeper, &contractCallerObj)
		require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got)
	})

	t.Run("multipleCheckpoint", func(t *testing.T) {
		t.Run("afterTimeout", func(t *testing.T) {
			ctx, keeper := CreateTestInput(t, false)
			// generate proposer for validator set
			LoadValidatorSet(4, t, keeper, ctx, false, 10)
			keeper.IncreamentAccum(ctx, 1)
			// gen random checkpoint
			header, err := GenRandCheckpointHeader(10)
			require.Empty(t, err, "Unable to create random header block, Error:%v", err)
			// add current proposer to header
			header.Proposer = keeper.GetValidatorSet(ctx).Proposer.Signer
			// make sure proposer has min ether
			contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
			// create checkpoint 257 seconds prev to current time
			header.TimeStamp = uint64(time.Now().Add(-(helper.CheckpointBufferTime + time.Second)).Unix())
			t.Log("Sending checkpoint with timestamp", "Timestamp", header.TimeStamp, "Current", time.Now().Unix())
			// send old checkpoint
			SentValidCheckpoint(header, keeper, ctx, contractCallerObj, t)
			header, err = GenRandCheckpointHeader(10)
			header.Proposer = keeper.GetValidatorSet(ctx).Proposer.Signer
			// create new checkpoint with current time
			header.TimeStamp = uint64(time.Now().Unix())

			msgCheckpoint := checkpoint.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.TimeStamp)
			// send new checkpoint which should replace old one
			got := checkpoint.HandleMsgCheckpoint(ctx, msgCheckpoint, keeper, &contractCallerObj)
			require.True(t, got.IsOK(), "expected send-checkpoint to be  ok, got %v", got)
		})

		t.Run("beforeTimeout", func(t *testing.T) {
			ctx, keeper := CreateTestInput(t, false)
			// generate proposer for validator set
			LoadValidatorSet(4, t, keeper, ctx, false, 10)
			keeper.IncreamentAccum(ctx, 1)
			header, err := GenRandCheckpointHeader(10)
			require.Empty(t, err, "Unable to create random header block, Error:%v", err)
			// add current proposer to header
			header.Proposer = keeper.GetValidatorSet(ctx).Proposer.Signer
			// make sure proposer has min ether
			contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
			// add current proposer to header
			header.Proposer = keeper.GetValidatorSet(ctx).Proposer.Signer
			SentValidCheckpoint(header, keeper, ctx, contractCallerObj, t)
			// create checkpoint msg
			msgCheckpoint := checkpoint.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, uint64(time.Now().Unix()))
			// send checkpoint to handler
			got := checkpoint.HandleMsgCheckpoint(ctx, msgCheckpoint, keeper, contractCallerObj)
			require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got)
		})
	})

}

func SentValidCheckpoint(header types.CheckpointBlockHeader, keeper hmcmn.Keeper, ctx sdk.Context, contractCallerObj mocks.IContractCaller, t *testing.T) {
	// add current proposer to header
	header.Proposer = keeper.GetValidatorSet(ctx).Proposer.Signer

	// create checkpoint msg
	msgCheckpoint := checkpoint.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.TimeStamp)

	// send checkpoint to handler
	got := checkpoint.HandleMsgCheckpoint(ctx, msgCheckpoint, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)

	// check if cache is set
	found := keeper.GetCheckpointCache(ctx, hmcmn.CheckpointCacheKey)
	require.Equal(t, true, found, "Checkpoint cache should exist")

	// check if checkpoint matches
	storedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)

	// ignoring time difference
	header.TimeStamp = storedHeader.TimeStamp
	require.Equal(t, header, storedHeader, "Header block Doesnt Match")
}

// test condition where checkpoint on mainchain gets confirmed after someone send no-ack on heimdall
func TestACKAfterNoACK(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper := CreateTestInput(t, false)
	// generate proposer for validator set
	LoadValidatorSet(4, t, keeper, ctx, false, 10)
	keeper.IncreamentAccum(ctx, 1)
	header, err := GenRandCheckpointHeader(10)
	require.Empty(t, err, "Unable to create random header block, Error:%v", err)
	contractCallerObj.On("GetBalance", keeper.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)
	SentValidCheckpoint(header, keeper, ctx, contractCallerObj, t)
	msgNoACK := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Add(-(helper.CheckpointBufferTime + time.Second)).Unix()))
	got := checkpoint.HandleMsgCheckpointNoAck(ctx, msgNoACK, keeper)
	require.True(t, got.IsOK(), "expected send-no-ack to be ok, got %v", got)

	contractCallerObj.On("GetHeaderInfo", uint64(10000)).Return(header.RootHash, header.StartBlock, header.EndBlock, nil)
	// create ack msg
	msgACK := checkpoint.NewMsgCheckpointAck(uint64(10000), uint64(time.Now().Unix()))
	// send ack to handler
	got = checkpoint.HandleMsgCheckpointAck(ctx, msgACK, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected send-ack to be ok, got %v", got)
}

func TestFirstNoACK(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	LoadValidatorSet(4, t, keeper, ctx, false, 10)
	keeper.IncreamentAccum(ctx, 1)
	msgNoACK := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Unix()))
	got := checkpoint.HandleMsgCheckpointNoAck(ctx, msgNoACK, keeper)
	require.True(t, got.IsOK(), "expected send-no-ack to be ok, got %v", got)

}

// Test Staking handlers
// ------

func TestHandleMsgValidatorJoin(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper := CreateTestInput(t, false)
	mockVals := GenRandomVal(1, 0, 10, 10, false, 1)
	// select first validator from slice
	mockVal := mockVals[0]
	t.Log("Mock validator generated", "Validator", mockVal.Signer.String())
	contractCallerObj.On("GetValidatorInfo", mock.Anything).Return(mockVal, nil)

	// insert new validator
	msgValJoin := staking.NewMsgValidatorJoin(uint64(mockVal.ID), mockVal.PubKey, mockVal.StartEpoch, mockVal.EndEpoch, json.Number(strconv.Itoa(int(mockVal.Accum))))
	got := staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator join to be ok, got %v", got)
	// validator is stored properly and signer is created properly
	storedVal, err := keeper.GetValidatorInfo(ctx, mockVal.Signer.Bytes())
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", mockVal.Signer.String(), err)
	require.Equal(t, mockVal.PubKey.Address(), storedVal.Signer, "Signer address should match")
	// signer to validator mapping should exist properly
	storedSigner, found := keeper.GetSignerFromValidatorID(ctx, mockVal.ID)
	require.True(t, found, "signer and validator address should be mapped, got %v", found)
	require.Equal(t, mockVal.Signer.Bytes(), storedSigner.Bytes(), "Signer address in signer=>validator map should be same")
	// insert validator again
	got = staking.HandleMsgValidatorJoin(ctx, msgValJoin, keeper, &contractCallerObj)
	require.True(t, !got.IsOK(), "expected validator join to be not-ok, got %v", got)

}

func TestHandleMsgValidatorExit(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}
	ctx, keeper := CreateTestInput(t, false)
	// pass 0 as time alive to generate non de-activated validators
	LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	validators[0].EndEpoch = 10
	msg := staking.NewMsgValidatorExit(uint64(validators[0].ID))
	contractCallerObj.On("GetValidatorInfo", validators[0].Signer).Return(validators[0], nil)

	got := staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
	require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)

	got = staking.HandleMsgValidatorExit(ctx, msg, keeper, &contractCallerObj)
	require.True(t, !got.IsOK(), "expected validator exit to be ok, got %v", got)
	validator, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)

	require.True(t, found, "Validator should be present even after deactivation")
	require.Equal(t, 10, int(validator.EndEpoch), "end epoch should be set to 10")

	keeper.UpdateACKCountWithValue(ctx, 20)
	currentVals := keeper.GetCurrentValidators(ctx)
	require.Equal(t, 3, len(currentVals), "No validators should exist after epoch passes")

	found = FindSigner(validators[0].Signer, currentVals)
	require.True(t, !found, "Validator should not exist in current val set")
}

func TestHandleMsgValidatorUpdate(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	// pass 0 as time alive to generate non de-activated validators
	LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	newSigner := GenRandomVal(1, 0, 10, 10, false, 1)
	msg := staking.NewMsgValidatorUpdate(uint64(validators[0].ID), newSigner[0].PubKey, json.Number(int(newSigner[0].Accum)))
	got := staking.HandleMsgSignerUpdate(ctx, msg, keeper)
	require.True(t, got.IsOK(), "expected validator update to be ok, got %v", got)
}

// ------

// finds address in give validator slice
func FindSigner(signer common.Address, vals []types.Validator) bool {
	for _, val := range vals {
		if bytes.Compare(signer.Bytes(), val.Signer.Bytes()) == 0 {
			return true
		}
	}
	return false
}
