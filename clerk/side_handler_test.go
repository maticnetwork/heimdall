package clerk_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Create test suite
//

// SideHandlerTestSuite integrate test suite context object
type SideHandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	sideHandler    hmTypes.SideTxHandler
	postHandler    hmTypes.PostTxHandler
	contractCaller mocks.IContractCaller
	chainID        string
	r              *rand.Rand
}

func (suite *SideHandlerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = clerk.NewSideTxHandler(suite.app.ClerkKeeper, &suite.contractCaller)
	suite.postHandler = clerk.NewPostTxHandler(suite.app.ClerkKeeper, &suite.contractCaller)

	// fetch chain id
	suite.chainID = suite.app.ChainKeeper.GetParams(suite.ctx).ChainParams.BorChainID

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SideHandlerTestSuite))
}

// Test cases
func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, uint32(sdk.CodeUnknownRequest), result.Code)
	require.Equal(t, abci.SideTxResultType_Skip, result.Result)
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgEventRecord() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r
	chainParams := app.ChainKeeper.GetParams(suite.ctx)

	_, _, addr1 := sdkAuth.KeyTestPubAddr()

	id := r.Uint64()

	t.Run("Success", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		msg := types.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			txHash,
			logIndex,
			blockNumber,
			id,
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			make([]byte, 0),
			suite.chainID,
		)

		// mock external calls
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		event := &statesender.StatesenderStateSynced{
			Id:              new(big.Int).SetUint64(msg.ID),
			ContractAddress: msg.ContractAddress.EthAddress(),
			Data:            msg.Data,
		}
		suite.contractCaller.On("DecodeStateSyncedEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.Nil(t, storedEventRecord)
		require.Error(t, err)
	})

	t.Run("NoReceipt", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(200)
		blockNumber := uint64(51)
		txHash := hmTypes.HexToHeimdallHash("no receipt hash")

		msg := types.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			txHash,
			logIndex,
			blockNumber,
			id,
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			make([]byte, 0),
			suite.chainID,
		)

		// mock external calls -- no receipt
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)
		suite.contractCaller.On("DecodeStateSyncedEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), nil, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	})

	t.Run("NoLog", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(100)
		blockNumber := uint64(510)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("no log hash")

		msg := types.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			txHash,
			logIndex,
			blockNumber,
			id,
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			make([]byte, 0),
			suite.chainID,
		)

		// mock external calls -- no receipt
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeStateSyncedEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	})

	t.Run("EventDataExceed", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}
		id := uint64(111)
		logIndex := uint64(1)
		blockNumber := uint64(1000)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		const letterBytes = "abcdefABCDEF"
		b := make([]byte, helper.LegacyMaxStateSyncSize+3)
		for i := range b {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		}

		// data created after trimming
		msg := types.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			txHash,
			logIndex,
			blockNumber,
			id,
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.HexToHexBytes(""),
			suite.chainID,
		)

		// mock external calls
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		event := &statesender.StatesenderStateSynced{
			Id:              new(big.Int).SetUint64(msg.ID),
			ContractAddress: msg.ContractAddress.EthAddress(),
			Data:            b,
		}
		suite.contractCaller.On("DecodeStateSyncedEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should pass")

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.Nil(t, storedEventRecord)
		require.Error(t, err)
	})

}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result := suite.postHandler(ctx, nil, abci.SideTxResultType_Yes)
	require.False(t, result.IsOK(), "Post handler should fail")
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgEventRecord() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, addr1 := sdkAuth.KeyTestPubAddr()

	id := r.Uint64()
	logIndex := r.Uint64()
	blockNumber := r.Uint64()
	txHash := hmTypes.HexToHeimdallHash("no log hash")

	msg := types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		txHash,
		logIndex,
		blockNumber,
		id,
		hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
		make([]byte, 0),
		suite.chainID,
	)

	t.Run("NoResult", func(t *testing.T) {
		result := suite.postHandler(ctx, msg, abci.SideTxResultType_No)
		require.False(t, result.IsOK(), "Post handler should fail")
		require.Equal(t, common.CodeSideTxValidationFailed, result.Code)
		require.Equal(t, 0, len(result.Events), "No error should be emitted for failed post-tx")

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.Nil(t, storedEventRecord)
		require.Error(t, err)
	})

	t.Run("YesResult", func(t *testing.T) {
		result := suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "Post handler should succeed")
		require.Greater(t, len(result.Events), 0, "Events should be emitted for successful post-tx")

		// sequence id
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		// check sequence
		hasSequence := app.ClerkKeeper.HasRecordSequence(ctx, sequence.String())
		require.True(t, hasSequence, "Sequence should be stored correctly")

		// there should be no stored event record
		storedEventRecord, err := app.ClerkKeeper.GetEventRecord(ctx, id)
		require.NotNil(t, storedEventRecord)
		require.NoError(t, err)
	})

	t.Run("Replay", func(t *testing.T) {
		id := r.Uint64()
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmTypes.HexToHeimdallHash("Replay hash")
		_, _, addr2 := sdkAuth.KeyTestPubAddr()

		msg := types.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			txHash,
			logIndex,
			blockNumber,
			id,
			hmTypes.BytesToHeimdallAddress(addr2.Bytes()),
			make([]byte, 0),
			suite.chainID,
		)

		result := suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "Post handler should succeed")

		result = suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.False(t, result.IsOK(), "Post handler should prevent replay attack")
		require.Equal(t, common.CodeOldTx, result.Code)
	})
}
