package topup_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/x/topup/test_helper"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	ethCommon "github.com/maticnetwork/bor/common"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/x/topup"
	"github.com/maticnetwork/heimdall/x/topup/types"
	abci "github.com/tendermint/tendermint/proto/tendermint/types"
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
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)

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
	require.Equal(t, sdkerrors.ErrUnknownRequest.ABCICode(), result.Code)
	require.Equal(t, abci.SideTxResultType_SKIP, result.Result)
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgTopup() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	chainParams := app.ChainKeeper.GetParams(suite.ctx)

	_, _, addr1 := testdata.KeyTestPubAddr()
	generated_address1, _ := sdk.AccAddressFromHex(addr1.String())
	_, _, addr2 := testdata.KeyTestPubAddr()

	t.Run("Success", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// sequence id
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		commonAddr := common.BytesToAddress(addr1.Bytes())
		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: commonAddr,
			Fee:  coins.AmountOf(hmTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		stakingInfoAddr, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		commonAddrdd := common.BytesToAddress(stakingInfoAddr)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", commonAddrdd, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)

		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_YES, result.Result, "Result should be `yes`")

		// there should be no stored event record
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.False(t, ok)
	})

	t.Run("NoReceipt", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, nil, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)

		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrWaitForConfirmation, result.Code)
	})

	t.Run("NoLog", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, txReceipt, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, sdkerrors.ErrTxDecode.ABCICode(), result.Code)
	})

	t.Run("BlockMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber + 1),
		}
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr1.Bytes()),
			Fee:  coins.AmountOf(hmTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		//require.Equal(t, common.ErrInvalidMsg.ABCICode(), result.Code)
	})

	t.Run("UserMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(addr2.Bytes()),
			Fee:  coins.AmountOf(hmTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), result.Code)
	})

	t.Run("FeeMismatch", func(t *testing.T) {
		suite.contractCaller = mocks.IContractCaller{}

		logIndex := uint64(10)
		blockNumber := uint64(599)
		txReceipt := &ethTypes.Receipt{
			BlockNumber: new(big.Int).SetUint64(blockNumber),
		}
		txHash := hmCommonTypes.HexToHeimdallHash("success hash")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
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
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result, err := suite.postHandler(ctx, nil, abci.SideTxResultType_YES)
	require.Nil(t, result)
	require.Error(t, err, "Post handler should fail")
	// require.Equal(t, sdkerrors.ErrUnknownRequest, err)
	// require.Equal(t, sdkerrors.ErrInvalidRequest.ABCICode(), result.Code)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgTopup() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, addr1 := testdata.KeyTestPubAddr()
	generated_address1, _ := sdk.AccAddressFromHex(addr1.String())
	_, _, addr2 := testdata.KeyTestPubAddr()
	generated_address2, _ := sdk.AccAddressFromHex(addr2.String())
	_, _, addr3 := testdata.KeyTestPubAddr()
	generated_address3, _ := sdk.AccAddressFromHex(addr3.String())

	logIndex := r.Uint64()
	blockNumber := r.Uint64()
	txHash := hmCommonTypes.HexToHeimdallHash("hash")

	t.Run("NoResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			addr1,
			addr1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result, err := suite.postHandler(ctx, &msg, abci.SideTxResultType_NO)
		require.Error(t, err)
		require.Nil(t, result)
		// require.False(t, result.IsOK(), "Post handler should fail")
		require.Equal(t, hmCommon.ErrSideTxValidation, err)

		// require.Equal(t, 0, len(result.Events), "No error should be emitted for failed post-tx")

		require.Nil(t, result)

		// there should be no stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.False(t, ok)

		// account coins should be empty
		acc1 := app.AccountKeeper.GetAccount(ctx, generated_address1)
		require.Nil(t, acc1)
	})

	t.Run("YesResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result, err := suite.postHandler(ctx, &msg, abci.SideTxResultType_YES)
		require.Nil(t, err)
		require.NotNil(t, result)

		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// account coins should be empty
		acc1 := app.AccountKeeper.GetAccount(ctx, addr1)

		require.NotNil(t, acc1)

		// require.False(t, acc1.GetCoins().Empty())
		require.False(t, app.BankKeeper.GetAllBalances(ctx, generated_address1).Empty())

		// require.True(t, acc1.GetCoins().IsEqual(coins)) // for same proposer
		require.True(t, app.BankKeeper.GetAllBalances(ctx, generated_address1).IsEqual(coins))
	})

	t.Run("WithProposer", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmCommonTypes.HexToHeimdallHash("hash with proposer")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			addr2, // different proposer
			addr3,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result, err := suite.postHandler(ctx, &msg, abci.SideTxResultType_YES)

		// require.True(t, result.IsOK(), "Post handler should succeed")
		require.NoError(t, err)
		require.NotNil(t, result)

		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// account coins should not be empty
		acc2 := app.AccountKeeper.GetAccount(ctx, generated_address2)
		require.NotNil(t, acc2)

		// require.False(t, acc2.GetCoins().Empty())
		require.False(t, app.BankKeeper.GetAllBalances(ctx, generated_address2).Empty())

		// account coins should be empty
		// acc3 := app.AccountKeeper.GetAccount(ctx, generated_address3)
		// require.NotNil(t, acc3)

		// require.False(t, acc3.GetCoins().Empty())
		// require.False(t, app.BankKeeper.GetAllBalances(ctx, generated_address3).Empty())

		// check coins = acc1.coins + acc2.coins
		// require.True(t, coins.IsEqual(acc3.GetCoins().Add(acc2.GetCoins()))) // for same proposer
		require.True(t, coins.IsEqual(app.BankKeeper.GetAllBalances(ctx, generated_address3).Add(app.BankKeeper.GetBalance(ctx, generated_address2, hmTypes.FeeToken)))) // for same proposer
	})

	t.Run("Replay", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmCommonTypes.HexToHeimdallHash("hash replay")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generated_address1,
			generated_address1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// check if incoming tx is older
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		result, err := suite.postHandler(ctx, &msg, abci.SideTxResultType_YES)

		// require.True(t, result.IsOK(), "Post handler should succeed")
		require.NoError(t, err)
		require.NotNil(t, result)

		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := app.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		result, err = suite.postHandler(ctx, &msg, abci.SideTxResultType_YES)
		require.Error(t, err)
		require.Nil(t, result)
		// require.False(t, result.IsOK(), "Post handler should fail while replaying same tx")

	})
}
