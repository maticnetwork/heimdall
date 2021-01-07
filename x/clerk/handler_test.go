package clerk_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/clerk"
	"github.com/maticnetwork/heimdall/x/clerk/types"
)

//
// Test suite
//

type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         client.Context
	chainID        string
	handler        sdk.Handler
	contractCaller helper.IContractCaller
	r              *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
	// TODO - Check this
	// suite.contractCaller = helper.IContractCaller{}
	suite.handler = clerk.NewHandler(suite.app.ClerkKeeper, suite.contractCaller)

	// fetch chain id
	suite.chainID = suite.app.ChainKeeper.GetParams(suite.ctx).ChainParams.BorChainID

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

//
// Test cases
//

func (suite *HandlerTestSuite) TestHandleMsgEventRecord() {
	t, app, ctx, chainID, r := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

	// keys and addresses
	_, _, addr1 := testdata.KeyTestPubAddr()

	id := r.Uint64()
	logIndex := r.Uint64()
	blockNumber := r.Uint64()

	// successful message
	msg := types.NewMsgEventRecord(
		addr1,
		hmCommon.HexToHeimdallHash("123"),
		logIndex,
		blockNumber,
		id,
		addr1,
		make([]byte, 0),
		chainID,
	)

	t.Run("Success", func(t *testing.T) {
		result, err := suite.handler(ctx, msg)
		require.Error(t, err)
		// require.True(t, result.IsOK(), "expected msg record to be ok, got %v", result)

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.Nil(t, storedEventRecord)
		require.Error(t, err)
	})

	t.Run("ExistingRecord", func(t *testing.T) {
		// store event record in keeper
		app.ClerkKeeper.SetEventRecord(ctx,
			types.NewEventRecord(
				msg.TxHash,
				msg.LogIndex,
				msg.Id,
				msg.ContractAddress,
				msg.Data,
				msg.ChainId,
				time.Now(),
			),
		)

		result, err := suite.handler(ctx, msg)
		require.Error(t, err)
		// require.False(t, result.IsOK(), "should fail due to existent event record but succeeded")
		// require.Equal(t, types.CodeEventRecordAlreadySynced, result.Code)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgEventRecordSequence() {
	t, app, ctx, chainID, r := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

	_, _, addr1 := testdata.KeyTestPubAddr()

	msg := types.NewMsgEventRecord(
		addr1,
		hmCommon.HexToHeimdallHash("123"),
		r.Uint64(),
		r.Uint64(),
		r.Uint64(),
		addr1,
		make([]byte, 0),
		chainID,
	)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))
	app.ClerkKeeper.SetRecordSequence(ctx, sequence.String())

	result, err := suite.handler(ctx, msg)
	require.Error(t, err)
	// require.False(t, result.IsOK(), "should fail due to existent sequence but succeeded")
	// require.Equal(t, common.CodeOldTx, result.Code)
}

func (suite *HandlerTestSuite) TestHandleMsgEventRecordChainID() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, addr1 := testdata.KeyTestPubAddr()

	id := r.Uint64()

	// wrong chain id
	msg := types.NewMsgEventRecord(
		addr1,
		hmCommon.HexToHeimdallHash("123"),
		r.Uint64(),
		r.Uint64(),
		id,
		addr1,
		make([]byte, 0),
		"random chain id",
	)
	result, err := suite.handler(ctx, msg)
	require.Error(t, err)
	// require.False(t, result.IsOK(), "error invalid bor chain id %v", result.Code)
	// require.Equal(t, common.CodeInvalidBorChainID, result.Code)

	// there should be no stored event record
	storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
	require.Nil(t, storedEventRecord)
	require.Error(t, err)
}
