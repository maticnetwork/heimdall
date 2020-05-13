package checkpoint_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	errs "github.com/maticnetwork/heimdall/common"
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
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		ID:        hmTypes.NewDividendAccountID(1),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)
	keeper.FlushCheckpointBuffer(ctx)

	// check valid checkpoint
	// generate proposer for validator set
	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)
	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// keeper.SetCheckpointBuffer(ctx, header)
	// require.Empty(t, err, "Unable to create random header block, Error:%v", err)

	// make sure proposer has min ether
	suite.contractCaller.On("GetBalance", stakingKeeper.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)

	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)
	// storedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	// t.Log("Header added to buffer", storedHeader.String())
	// require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointWithInvalidProposer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		ID:        hmTypes.NewDividendAccountID(1),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)
	keeper.FlushCheckpointBuffer(ctx)

	// check invalid proposer
	// generate proposer for validator set
	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)

	// add wrong proposer to header
	header.Proposer = hmTypes.HexToHeimdallAddress("1234")

	got := suite.SendCheckpoint(header)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAfterBufferTimeOut() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		ID:        hmTypes.NewDividendAccountID(1),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)
	keeper.FlushCheckpointBuffer(ctx)

	// generate proposer for validator set
	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)
	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)
	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// make sure proposer has min ether
	suite.contractCaller.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)

	// create checkpoint `checkpointBufferTime` seconds prev to current time
	checkpointBufferTime := params.CheckpointBufferTime

	header.TimeStamp = uint64(time.Now().Add(-(checkpointBufferTime + time.Second)).Unix())
	t.Log("Sending checkpoint with timestamp", "Timestamp", header.TimeStamp, "Current", time.Now().UTC().Unix())
	// send old checkpoint
	res := suite.SendCheckpoint(header)
	require.True(t, res.IsOK(), "expected send-checkpoint to be  ok, got %v", res)

	// lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	// if err == nil {
	// 	start = start + lastCheckpoint.EndBlock + 1
	// }
	// header, err = suite.GenRandCheckpointHeader(0, maxSize, params.MaxCheckpointLength)
	// header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer
	// create new checkpoint with current time
	header.TimeStamp = uint64(time.Now().Unix())

	// send new checkpoint which should replace old one
	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-checkpoint to be  ok, got %v", got)
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointExistInBuffer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		ID:        hmTypes.NewDividendAccountID(1),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)
	keeper.FlushCheckpointBuffer(ctx)

	cmn.LoadValidatorSet(4, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)
	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)
	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// make sure proposer has min ether
	suite.contractCaller.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)

	// send old checkpoint
	res := suite.SendCheckpoint(header)
	require.True(t, res.IsOK(), "expected send-checkpoint to be  ok, got %v", res)

	// TODO: check why not adding to buffer
	// send checkpoint to handler
	// got := suite.SendCheckpoint(header)
	// require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))

}

// GenRandCheckpointHeader create random header block
func (suite *HandlerTestSuite) GenRandCheckpointHeader(start uint64, headerSize uint64, maxCheckpointLenght uint64) (headerBlock hmTypes.CheckpointBlockHeader, err error) {
	app, ctx := suite.app, suite.ctx

	topupKeeper := app.TopupKeeper
	end := start + headerSize
	borChainID := "1234"

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	rootHash := hmTypes.BytesToHeimdallHash(accRootHash)
	proposer := hmTypes.HeimdallAddress{}
	headerBlock = hmTypes.CreateBlock(
		start,
		end,
		rootHash,
		proposer,
		borChainID,
		uint64(time.Now().UTC().Unix()))

	return headerBlock, nil
}

func (suite *HandlerTestSuite) SendCheckpoint(header hmTypes.CheckpointBlockHeader) (res sdk.Result) {
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
	return suite.handler(ctx, msgCheckpoint)
}
