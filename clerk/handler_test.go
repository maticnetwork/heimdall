package clerk_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Test suite
//

type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	chainID        string
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = clerk.NewHandler(suite.app.ClerkKeeper, &suite.contractCaller)

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
	_, _, addr1 := sdkAuth.KeyTestPubAddr()

	id := r.Uint64()
	logIndex := r.Uint64()
	blockNumber := r.Uint64()

	// successful message
	msg := types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		hmTypes.HexToHeimdallHash("123"),
		logIndex,
		blockNumber,
		id,
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		make([]byte, 0),
		chainID,
	)

	t.Run("Success", func(t *testing.T) {
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "expected msg record to be ok, got %v", result)

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
				msg.ID,
				msg.ContractAddress,
				msg.Data,
				msg.ChainID,
				time.Now(),
			),
		)

		result := suite.handler(ctx, msg)
		require.False(t, result.IsOK(), "should fail due to existent event record but succeeded")
		require.Equal(t, types.CodeEventRecordAlreadySynced, result.Code)
	})

	t.Run("EventSizeExceed", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		const letterBytes = "abcdefABCDEF"
		b := make([]byte, helper.LegacyMaxStateSyncSize+3)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		msg.Data = b

		err := msg.ValidateBasic()
		require.Error(t, err)
	})

}

func (suite *HandlerTestSuite) TestHandleMsgEventRecordSequence() {
	t, app, ctx, chainID, r := suite.T(), suite.app, suite.ctx, suite.chainID, suite.r

	_, _, addr1 := sdkAuth.KeyTestPubAddr()

	msg := types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		hmTypes.HexToHeimdallHash("123"),
		r.Uint64(),
		r.Uint64(),
		r.Uint64(),
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		make([]byte, 0),
		chainID,
	)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))
	app.ClerkKeeper.SetRecordSequence(ctx, sequence.String())

	result := suite.handler(ctx, msg)
	require.False(t, result.IsOK(), "should fail due to existent sequence but succeeded")
	require.Equal(t, common.CodeOldTx, result.Code)
}

func (suite *HandlerTestSuite) TestHandleMsgEventRecordChainID() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, addr1 := sdkAuth.KeyTestPubAddr()

	id := r.Uint64()

	// wrong chain id
	msg := types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		hmTypes.HexToHeimdallHash("123"),
		r.Uint64(),
		r.Uint64(),
		id,
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		make([]byte, 0),
		"random chain id",
	)
	result := suite.handler(ctx, msg)
	require.False(t, result.IsOK(), "error invalid bor chain id %v", result.Code)
	require.Equal(t, common.CodeInvalidBorChainID, result.Code)

	// there should be no stored event record
	storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
	require.Nil(t, storedEventRecord)
	require.Error(t, err)
}
