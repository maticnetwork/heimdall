package checkpoint

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/types"
	"github.com/mn/heimdall/checkpoint"
	"github.com/stretchr/testify/require"
)

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
			got := checkpoint.HandleMsgCheckpoint(ctx, msgCheckpoint, keeper, &contractCallerObj)
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
