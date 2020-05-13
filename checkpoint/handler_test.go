package checkpoint_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	cmn "github.com/maticnetwork/heimdall/test"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app    *app.HeimdallApp
	ctx    sdk.Context
	cliCtx context.CLIContext

	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = checkpoint.NewHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

// test handler for message
func (suite *HandlerTestSuite) TestHandleMsgCheckpoint() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		ID:        hmTypes.NewDividendAccountID(1),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)
	keeper.FlushCheckpointBuffer(ctx)
	// check valid checkpoint
	suite.Run("validCheckpoint", func() {
		// generate proposer for validator set
		cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
		stakingKeeper.IncrementAccum(ctx, 1)
		start := uint64(0)
		end := uint64(256)
		lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
		if err == nil {
			start = start + lastCheckpoint.EndBlock + 1
		}
		dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
		accRootHash, err := types.GetAccountRootHash(dividendAccounts)
		rootHash := hmTypes.BytesToHeimdallHash(accRootHash)

		// add current proposer to header
		proposer := stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

		header, err := cmn.GenRandCheckpointHeader(start, end, rootHash, proposer, params.MaxCheckpointLength, suite.contractCaller)
		// keeper.SetCheckpointBuffer(ctx, header)
		// require.Empty(t, err, "Unable to create random header block, Error:%v", err)
		// make sure proposer has min ether
		suite.contractCaller.On("GetBalance", stakingKeeper.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)

		suite.SentValidCheckpoint(header)
	})

	// check invalid proposer
	// suite.Run("invalidProposer", func() {
	// 	// generate proposer for validator set
	// 	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	// 	stakingKeeper.IncrementAccum(ctx, 1)
	// 	header, err := GenRandCheckpointHeader(0, 10)
	// 	require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	// 	// add wrong proposer to header
	// 	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Validators[2].Signer

	// 	accs := stakingKeeper.GetAllDividendAccounts(ctx)
	// 	root, err := types.GetAccountRootHash(accs)

	// 	header.AccountRootHash = hmTypes.BytesToHeimdallHash(root)

	// 	// create checkpoint msg
	// 	msgCheckpoint := types.NewMsgCheckpointBlock(header.Proposer,
	// 		header.StartBlock,
	// 		header.EndBlock,
	// 		header.RootHash,
	// 		header.AccountRootHash,
	// 		header.TimeStamp) // send checkpoint to handler
	// 	got := suite.handler(ctx, msgCheckpoint)
	// 	require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got.IsOK())
	// })

	// suite.Run("multipleCheckpoint", func() {
	// 	suite.Run("afterTimeout", func() {
	// 		ctx, stakingKeeper, ck := CreateTestInput(t, false)
	// 		// generate proposer for validator set
	// 		cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	// 		stakingKeeper.IncrementAccum(ctx, 1)
	// 		header, err := GenRandCheckpointHeader(0, 10)
	// 		require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	// 		// add current proposer to header
	// 		header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// 		// make sure proposer has min ether
	// 		contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
	// 		// create checkpoint 257 seconds prev to current time
	// 		header.TimeStamp = uint64(time.Now().Add(-(helper.CheckpointBufferTime + time.Second)).Unix())
	// 		t.Log("Sending checkpoint with timestamp", "Timestamp", header.TimeStamp, "Current", time.Now().UTC().Unix())
	// 		// send old checkpoint
	// 		SentValidCheckpoint(header, ck, stakingKeeper, ctx, contractCallerObj, t)

	// 		header, err = GenRandCheckpointHeader(0, 10)
	// 		header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// 		// create new checkpoint with current time
	// 		header.TimeStamp = uint64(time.Now().Unix())
	// 		accs := stakingKeeper.GetAllDividendAccounts(ctx)
	// 		root, err := types.GetAccountRootHash(accs)

	// 		header.AccountRootHash = hmTypes.BytesToHeimdallHash(root)

	// 		msgCheckpoint := types.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.AccountRootHash, header.TimeStamp)
	// 		// send new checkpoint which should replace old one
	// 		got := suite.handler(ctx, msgCheckpoint)
	// 		require.True(t, got.IsOK(), "expected send-checkpoint to be  ok, got %v", got)
	// 	})

	// suite.Run("beforeTimeout", func() {
	// 	ctx, stakingKeeper, ck := CreateTestInput(t, false)
	// 	// generate proposer for validator set
	// 	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	// 	stakingKeeper.IncrementAccum(ctx, 1)
	// 	header, err := GenRandCheckpointHeader(0, 10)
	// 	require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	// 	// add current proposer to header
	// 	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// 	// make sure proposer has min ether
	// 	contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
	// 	// add current proposer to header
	// 	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// 	// send old checkpoint
	// 	SentValidCheckpoint(header, ck, stakingKeeper, ctx, contractCallerObj, t)
	// 	accs := stakingKeeper.GetAllDividendAccounts(ctx)
	// 	root, err := types.GetAccountRootHash(accs)

	// 	header.AccountRootHash = hmTypes.BytesToHeimdallHash(root)

	// 	// create checkpoint msg
	// 	msgCheckpoint := types.NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.AccountRootHash, uint64(time.Now().Unix()))

	// 	// send checkpoint to handler
	// 	got := suite.handler(ctx, msgCheckpoint)
	// 	require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got)
	// })
	// })

}

func (suite *HandlerTestSuite) SentValidCheckpoint(header hmTypes.CheckpointBlockHeader) {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	// keeper := app.CheckpointKeeper
	topupKeeper := app.TopupKeeper

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)
	accountRoot := hmTypes.BytesToHeimdallHash(accRootHash)

	borChainId := "1234"
	// create checkpoint msg
	msgCheckpoint := types.NewMsgCheckpointBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		accountRoot,
		borChainId,
	)

	t.Log("Checkpoint msg created", msgCheckpoint)

	// send checkpoint to handler
	got := suite.handler(ctx, msgCheckpoint)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)
	// storedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	// t.Log("Header added to buffer", storedHeader.String())
	// require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)
}
