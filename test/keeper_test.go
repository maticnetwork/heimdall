package test

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdateAck(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)
	keeper.UpdateACKCount(ctx)
	ack := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(1), ack, "Ack Count Not Equal")
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
	require.Empty(t, err, "Unable to retrieve checkpoint, Error: %v", err)
	require.Equal(t, headerBlock, storedHeader, "Header Blocks dont match")

	// flush and check if its flushed
	keeper.FlushCheckpointBuffer(ctx)
	storedHeader, err = keeper.GetCheckpointFromBuffer(ctx)
	require.NotEmpty(t, err, "HeaderBlock should not exist after flush")

	//TODO add this check for handler test
	//err = keeper.SetCheckpointBuffer(ctx, headerBlock)
	//if err == nil {
	//	require.Fail(t, "Checkpoint should not be stored if checkpoint already exists in buffer")
	//}
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

func TestValidatorAdd(t *testing.T) {
	ctx, keeper := CreateTestInput(t, false)

	validator := GenRandomVal()

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

//TODO add tests for validator set changes on update/signer
// TODO add mocks for contract calls/tx to test
