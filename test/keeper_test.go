package test

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	hmcmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"testing"
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

	vals := GenRandomVal(1)
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
	var valSet types.ValidatorSet
	LoadValidatorSet(4, t, keeper, ctx, &valSet)

	storedValSet := keeper.GetValidatorSet(ctx)
	require.Equal(t, valSet, storedValSet, "Validator Set in state doesnt match ")

	keeper.IncreamentAccum(ctx, 1)
	initialProposer := keeper.GetCurrentProposer(ctx)

	keeper.IncreamentAccum(ctx, 1)
	newProposer := keeper.GetCurrentProposer(ctx)
	fmt.Printf("Prev :%#v  , New : %#v", initialProposer, newProposer)

}

func TestValUpdates(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	var valSet types.ValidatorSet
	// load 4 validators to state
	LoadValidatorSet(4, t, keeper, ctx, &valSet)

	// select proposer
	keeper.IncreamentAccum(ctx, 1)

	currentValidatorSet := keeper.GetValidatorSet(ctx)

	for _, v := range currentValidatorSet.Validators {
		t.Logf("===>>>>>>", v.Address.Hex())
	}

	// remove first validator
	currentValidatorSet.Validators[0].StartEpoch = keeper.GetACKCount(ctx) + 10
	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")

	// apply updates
	helper.UpdateValidators(
		&currentValidatorSet,                // pointer to current validator set -- UpdateValidators will modify it
		keeper.GetAllValidators(ctx),        // All validators
		keeper.GetValidatorToSignerMap(ctx), // validator to signer map
		keeper.GetACKCount(ctx),             // ack count
	)

	keeper.IncreamentAccum(ctx, 1)
	newProposer := keeper.GetCurrentProposer(ctx)
	newValSet := keeper.GetValidatorSet(ctx)
	for _, v := range newValSet.Validators {
		t.Logf("===>>>>>>", v.Address.Hex())
	}

}

func LoadValidatorSet(count int, t *testing.T, keeper hmcmn.Keeper, ctx sdk.Context, validatorSet *types.ValidatorSet) {
	// create 4 validators
	validators := GenRandomVal(4)

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
}

//TODO add tests for validator set changes on update/signer
// TODO add mocks for contract calls/tx to test
