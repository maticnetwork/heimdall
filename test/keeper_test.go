package test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	hmcmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

func TestUpdateAck(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.UpdateACKCount(ctx)
	ack := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(2), ack, "Ack Count Not Equal")
}

func TestCheckpointBuffer(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	// create random header block
	headerBlock, err := GenRandCheckpointHeader()
	require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	// set checkpoint
	err = keeper.SetCheckpointBuffer(ctx, headerBlock)
	require.Empty(t, err, "Unable to store checkpoint, Error: %v", err)

	// check if we are able to get checkpoint after set
	storedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	t.Log("Checkpoint", storedHeader)
	require.Empty(t, err, "Unable to retrieve checkpoint, Error: %v", err)
	require.Equal(t, headerBlock, storedHeader, "Header Blocks dont match")

	// flush and check if its flushed
	keeper.FlushCheckpointBuffer(ctx)
	storedHeader, err = keeper.GetCheckpointFromBuffer(ctx)
	require.NotEmpty(t, err, "HeaderBlock should not exist after flush")
}

func TestCheckpointACK(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	prevACK := keeper.GetACKCount(ctx)

	// create random header block
	headerBlock, err := GenRandCheckpointHeader()
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

	vals := GenRandomVal(1, 0, 0, 10, true)
	validator := vals[0]

	err := keeper.AddValidator(ctx, validator)
	require.Empty(t, err, "Unable to set validator, Error: %v", err)

	storedVal, err := keeper.GetValidatorInfo(ctx, validator.Signer.Bytes())
	require.Empty(t, err, "Unable to fetch validator")
	require.Equal(t, validator, storedVal, "Unable to fetch validator from val address")

	storedSigner, ok := keeper.GetSignerFromValidator(ctx, validator.Address)
	require.Equal(t, true, ok, "Validator<=>Signer not mapped")
	require.Equal(t, validator.Signer, storedSigner, "Signer doesnt match")

	storedValidator, ok := keeper.GetValidatorFromValAddr(ctx, validator.Address)
	require.Equal(t, true, ok, "Validator<=>Signer not mapped")
	require.Equal(t, validator, storedValidator, "Unable to fetch validator from val address")

	valToSignerMap := keeper.GetValidatorToSignerMap(ctx)
	mappedSigner := valToSignerMap[hex.EncodeToString(validator.Address.Bytes())]
	require.Equal(t, validator.Signer, mappedSigner, "GetValidatorToSignerMap doesnt give right signer")
}

func TestValidatorSet(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	valSet := LoadValidatorSet(4, t, keeper, ctx, true)

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
		LoadValidatorSet(4, t, keeper, ctx, false)
		initValSet := keeper.GetValidatorSet(ctx)

		currentValSet := initValSet.Copy()
		prevValidatorSet := initValSet.Copy()

		// remove validator (making IsCurrentValidator return false)
		prevValidatorSet.Validators[0].StartEpoch = 20

		t.Log("Updated Validators in state")
		for _, v := range prevValidatorSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}

		err := keeper.AddValidator(ctx, *prevValidatorSet.Validators[0])
		require.Empty(t, err, "Unable to update validator set")

		// apply updates
		helper.UpdateValidators(
			currentValSet,                       // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx),        // All validators
			keeper.GetValidatorToSignerMap(ctx), // validator to signer map
			5,                                   // ack count
		)
		updatedValSet := currentValSet
		t.Log("Validators in updated validator set")
		for _, v := range updatedValSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}
		// check if 1 validator is removed
		require.Equal(t, len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators), "Validator set should be reduced by one ")
		// remove first validator from initial validator set and equate with new
		require.Equal(t, append(prevValidatorSet.Validators[:0], prevValidatorSet.Validators[1:]...), updatedValSet.Validators, "Validator at 0 index should be deleted")
	})

	t.Run("add", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)

		// load 4 validators to state
		LoadValidatorSet(4, t, keeper, ctx, false)
		initValSet := keeper.GetValidatorSet(ctx)

		validators := GenRandomVal(1, 0, 10, 10, false)
		prevValSet := initValSet.Copy()
		valToBeAdded := validators[0]
		currentValSet := initValSet.Copy()
		//prevValidatorSet := initValSet.Copy()
		keeper.AddValidator(ctx, valToBeAdded)

		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}
		t.Log("Val to be Added")
		t.Log("-->", "Address", valToBeAdded.Address.String(), "StartEpoch", valToBeAdded.StartEpoch, "EndEpoch", valToBeAdded.EndEpoch, "Power", valToBeAdded.Power)

		helper.UpdateValidators(
			currentValSet,                       // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx),        // All validators
			keeper.GetValidatorToSignerMap(ctx), // validator to signer map
			10,                                  // ack count
		)
		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "Power", v.Power)
		}

		require.Equal(t, len(prevValSet.Validators)+1, len(currentValSet.Validators), "Number of validators should be increased by 1")
		require.Equal(t, true, currentValSet.HasAddress(valToBeAdded.Address.Bytes()), "New Validator should be added")
		require.Equal(t, prevValSet.TotalVotingPower()+int64(valToBeAdded.Power), currentValSet.TotalVotingPower(), "Total power should be increased")
	})

	t.Run("update", func(t *testing.T) {
		ctx, keeper := CreateTestInput(t, false)

		// load 4 validators to state
		LoadValidatorSet(4, t, keeper, ctx, false)
		initValSet := keeper.GetValidatorSet(ctx)
		keeper.IncreamentAccum(ctx, 2)
		prevValSet := initValSet.Copy()
		currentValSet := keeper.GetValidatorSet(ctx)
		valToUpdate := currentValSet.Validators[0]
		newSigner := GenRandomVal(1, 0, 10, 10, false)
		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "Accum", v.Accum, "Signer", v.Signer.String(), "Total power", currentValSet.TotalVotingPower())
		}
		keeper.UpdateSigner(ctx, newSigner[0].Signer, newSigner[0].PubKey, valToUpdate.Signer)
		helper.UpdateValidators(
			&currentValSet,                      // pointer to current validator set -- UpdateValidators will modify it
			keeper.GetAllValidators(ctx),        // All validators
			keeper.GetValidatorToSignerMap(ctx), // validator to signer map
			10,                                  // ack count
		)
		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Address.String(), "Accum", v.Accum, "Signer", v.Signer.String(), "Total power", currentValSet.TotalVotingPower())
		}

		require.Equal(t, len(prevValSet.Validators), len(currentValSet.Validators), "Number of validators should remain same")
		_, val := currentValSet.GetByAddress(valToUpdate.Address.Bytes())
		require.Equal(t, newSigner[0].Signer, val.Signer, "Signer address should change")
		require.Equal(t, newSigner[0].PubKey, val.PubKey, "Signer pubkey should change")
		require.Equal(t, valToUpdate.Accum, val.Accum, "Validator accum should not change")
		require.Equal(t, prevValSet.TotalVotingPower(), currentValSet.TotalVotingPower(), "Total power should not change")
		// TODO not sure if proposer check is needed
		//require.Equal(t, &initValSet.Proposer.Address, &currentValSet.Proposer.Address, "Proposer should not change")
	})

}

func LoadValidatorSet(count int, t *testing.T, keeper hmcmn.Keeper, ctx sdk.Context, randomise bool) types.ValidatorSet {
	// create 4 validators
	validators := GenRandomVal(4, 0, 10, 10, randomise)

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

//TODO add tests for validator set changes on update/signer
// TODO add mocks for contract calls/tx to test
