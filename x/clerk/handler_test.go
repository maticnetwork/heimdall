package clerk_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	// "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
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
	// cliCtx         client.Context
	chainID        string
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = clerk.NewHandler(suite.app.ClerkKeeper, &suite.contractCaller)

	// TODO - Check this
	// fetch chain id
	// suite.chainID = suite.app.ChainKeeper.GetParams(suite.ctx).ChainParams.BorChainID
	suite.chainID = "testchainid"

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
		result, err := suite.handler(ctx, &msg)
		require.Nil(t, err)
		require.NotNil(t, result, "expected msg record to be ok, got %v", result)

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.Nil(t, storedEventRecord)
		require.Error(t, err)
	})

	t.Run("ExistingRecord", func(t *testing.T) {
		addr, _ := sdk.AccAddressFromHex(msg.ContractAddress)
		// store event record in keeper
		err := app.ClerkKeeper.SetEventRecord(ctx,
			types.NewEventRecord(
				hmCommon.HexToHeimdallHash(msg.TxHash),
				msg.LogIndex,
				msg.Id,
				addr,
				msg.Data,
				msg.ChainId,
				time.Now(),
			),
		)
		require.Nil(t, err)

		result, err := suite.handler(ctx, &msg)
		require.Error(t, err)
		require.Nil(t, result, "should fail due to existent event record but succeeded")
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

	result, err := suite.handler(ctx, &msg)
	require.Error(t, err)
	require.Nil(t, result, "should fail due to existent sequence but succeeded")
	// require.False(t, result.IsOK(), "should fail due to existent sequence but succeeded")
	// require.Equal(t, common.CodeOldTx, result.Code)
}

// TODO - Check this
// func (suite *HandlerTestSuite) TestHandleMsgEventRecordChainID() {
// 	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

// 	_, _, addr1 := testdata.KeyTestPubAddr()

// 	id := r.Uint64()

// 	// wrong chain id
// 	msg := types.NewMsgEventRecord(
// 		addr1,
// 		hmCommon.HexToHeimdallHash("123"),
// 		r.Uint64(),
// 		r.Uint64(),
// 		id,
// 		addr1,
// 		make([]byte, 0),
// 		"testchainid",
// 	)
// 	_, err := suite.handler(ctx, &msg)
// 	require.Error(t, err)
// 	// require.False(t, result.IsOK(), "error invalid bor chain id %v", result.Code)
// 	// require.Equal(t, common.CodeInvalidBorChainID, result.Code)

// 	// there should be no stored event record
// 	storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
// 	require.Nil(t, storedEventRecord)
// 	require.Error(t, err)
// }
