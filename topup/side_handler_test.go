package topup_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	r              *rand.Rand
}

func (suite *SideHandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = topup.NewSideTxHandler(suite.app.TopupKeeper, &suite.contractCaller)
	suite.postHandler = topup.NewPostTxHandler(suite.app.TopupKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SideHandlerTestSuite))
}

//
// Test cases
//

func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, uint32(sdk.CodeUnknownRequest), result.Code)
	require.Equal(t, abci.SideTxResultType_Skip, result.Result)
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgTopup() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	chainParams := app.ChainKeeper.GetParams(suite.ctx)

	_, _, addr1 := sdkAuth.KeyTestPubAddr()
	_, _, addr2 := sdkAuth.KeyTestPubAddr()

	t.Run("Success", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// sequence id
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr1.Bytes()),
			Fee:  coins.AmountOf(authTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")

		// there should be no stored event record
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.False(t, ok)
	})

	t.Run("NoReceipt", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), nil, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeWaitFrConfirmation), result.Code)
	})

	t.Run("NoLog", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeErrDecodeEvent), result.Code)
	})

	t.Run("BlockMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber + 1),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr1.Bytes()),
			Fee:  coins.AmountOf(authTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	t.Run("UserMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr2.Bytes()),
			Fee:  coins.AmountOf(authTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	t.Run("FeeMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr1.Bytes()),
			Fee:  big.NewInt(1), // different fee
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress.EthAddress(), txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, msg)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result := suite.postHandler(ctx, nil, abci.SideTxResultType_Yes)
	require.False(t, result.IsOK(), "Post handler should fail")
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgTopup() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, addr1 := sdkAuth.KeyTestPubAddr()
	_, _, addr2 := sdkAuth.KeyTestPubAddr()
	_, _, addr3 := sdkAuth.KeyTestPubAddr()

	logIndex := r.Uint64()
	blockNumber := r.Uint64()
	txHash := hmTypes.HexToHeimdallHash("hash")

	t.Run("NoResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result := suite.postHandler(ctx, msg, abci.SideTxResultType_No)
		require.False(t, result.IsOK(), "Post handler should fail")
		require.Equal(t, common.CodeSideTxValidationFailed, result.Code)
		require.Equal(t, 0, len(result.Events), "No error should be emitted for failed post-tx")

		// there should be no stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.False(t, ok)

		// account coins should be empty
		acc1 := app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
		require.Nil(t, acc1)
	})

	t.Run("YesResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result := suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "Post handler should succeed")
		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// account coins should be empty
		acc1 := app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
		require.NotNil(t, acc1)
		require.False(t, acc1.GetCoins().Empty())
		require.True(t, acc1.GetCoins().IsEqual(coins)) // for same proposer
	})

	t.Run("WithProposer", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmTypes.HexToHeimdallHash("hash with proposer")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr2.Bytes()), // different proposer
			hmTypes.BytesToHeimdallAddress(addr3.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result := suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "Post handler should succeed")
		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// account coins should not be empty
		acc2 := app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
		require.NotNil(t, acc2)
		require.False(t, acc2.GetCoins().Empty())

		// account coins should be empty
		acc3 := app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr3))
		require.NotNil(t, acc3)
		require.False(t, acc3.GetCoins().Empty())

		// check coins = acc1.coins + acc2.coins
		require.True(t, coins.IsEqual(acc3.GetCoins().Add(acc2.GetCoins()))) // for same proposer
	})

	t.Run("Replay", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmTypes.HexToHeimdallHash("hash replay")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr1.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result := suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "Post handler should succeed")
		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		result = suite.postHandler(ctx, msg, abci.SideTxResultType_Yes)
		require.False(t, result.IsOK(), "Post handler should fail while replaying same tx")
	})
}
