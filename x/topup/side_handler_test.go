package topup_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/x/topup/test_helper"

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

var (
	SuccessCode = uint32(0)
)

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

func createAccount() (cryptotypes.PrivKey, cryptotypes.PubKey, sdk.AccAddress) {
	newPrivKey := secp256k1.GenPrivKey()
	newPubKey := newPrivKey.PubKey()
	newSigner := sdk.AccAddress(newPrivKey.PubKey().Address().Bytes())
	return newPrivKey, newPubKey, newSigner
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgTopup() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	chainParams := initApp.ChainKeeper.GetParams(suite.ctx)

	_, _, generatedAddress1 := createAccount()
	_, _, addr2 := createAccount()

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
			generatedAddress1,
			generatedAddress1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// sequence id
		bn := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(bn, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		commonAddr := common.BytesToAddress(generatedAddress1)
		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: commonAddr,
			Fee:  coins.AmountOf(hmTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stateSenderAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", stateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)

		require.Equal(t, SuccessCode, result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_YES, result.Result, "Result should be `yes`")

		// there should be no stored event record
		ok := initApp.TopupKeeper.HasTopupSequence(ctx, sequence.String())
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
			generatedAddress1,
			generatedAddress1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StateSenderAddress, nil, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		fmt.Printf("fee result %+v\n", result)
		require.NotEqual(t, SuccessCode, result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrWaitForConfirmation.ABCICode(), result.Code)
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
			generatedAddress1,
			generatedAddress1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		StateSenderAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", StateSenderAddress, txReceipt, logIndex).Return(nil, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, SuccessCode, result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrDecodeEvent.ABCICode(), result.Code)
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
			generatedAddress1,
			generatedAddress1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(generatedAddress1.Bytes()),
			Fee:  coins.AmountOf(hmTypes.FeeToken).BigInt(),
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		stateSenderAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", stateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		fmt.Printf("Result %+v\n", result)

		fmt.Printf("result %+v\n", result)
		require.NotEqual(t, SuccessCode, result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrInvalidMsg.ABCICode(), result.Code)
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
			generatedAddress1,
			generatedAddress1,
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
		stateSenderAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", stateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, SuccessCode, result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrInvalidMsg.ABCICode(), result.Code)
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
			generatedAddress1,
			generatedAddress1,
			coins.AmountOf(hmTypes.FeeToken),
			txHash,
			logIndex,
			blockNumber,
		)

		// mock external call
		event := &stakinginfo.StakinginfoTopUpFee{
			User: ethCommon.BytesToAddress(generatedAddress1.Bytes()),
			Fee:  big.NewInt(1), // different fee
		}
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		stateSenderAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StateSenderAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorTopupFeesEvent", stateSenderAddress, txReceipt, logIndex).Return(event, nil)

		// execute handler
		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, SuccessCode, result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		require.Equal(t, hmCommon.ErrInvalidMsg.ABCICode(), result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result, err := suite.postHandler(ctx, nil, abci.SideTxResultType_YES)
	require.Nil(t, result)
	require.Error(t, err, "Post handler should fail")
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgTopup() {
	t, initApp, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	_, _, generatedAddress1 := createAccount()
	_, _, generatedAddress2 := createAccount()
	_, _, generatedAddress3 := createAccount()

	logIndex := r.Uint64()
	blockNumber := r.Uint64()
	txHash := hmCommonTypes.HexToHeimdallHash("hash")

	t.Run("NoResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generatedAddress1,
			generatedAddress1,
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
		require.Equal(t, hmCommon.ErrSideTxValidation, err)
		// there should be no stored sequence
		ok := initApp.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.False(t, ok)

		// account coins should be empty
		acc1 := initApp.AccountKeeper.GetAccount(ctx, generatedAddress1)
		require.Nil(t, acc1)
	})

	t.Run("YesResult", func(t *testing.T) {
		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generatedAddress1,
			generatedAddress1,
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
		ok := initApp.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// account coins should be empty
		acc1 := initApp.AccountKeeper.GetAccount(ctx, generatedAddress1)
		require.NotNil(t, acc1)
		require.False(t, initApp.BankKeeper.GetAllBalances(ctx, generatedAddress1).Empty())
		require.True(t, initApp.BankKeeper.GetAllBalances(ctx, generatedAddress1).IsEqual(coins))
	})

	t.Run("WithProposer", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmCommonTypes.HexToHeimdallHash("hash with proposer")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generatedAddress2, // different proposer
			generatedAddress3,
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
		ok := initApp.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		// todo: we need to check this test cases with get account
		// account coins should be empty
		require.True(t, initApp.BankKeeper.GetBalance(ctx, generatedAddress3, types.FeeToken).IsZero())
		require.False(t, initApp.BankKeeper.GetBalance(ctx, generatedAddress2, types.FeeToken).IsZero())

		require.True(t, coins.IsEqual(initApp.BankKeeper.GetAllBalances(ctx, generatedAddress3).Add(initApp.BankKeeper.GetBalance(ctx, generatedAddress2, hmTypes.FeeToken)))) // for same proposer
	})

	t.Run("Replay", func(t *testing.T) {
		logIndex := r.Uint64()
		blockNumber := r.Uint64()
		txHash := hmCommonTypes.HexToHeimdallHash("hash replay")

		// set coins
		coins := simulation.RandomFeeCoins()

		// topup msg
		msg := types.NewMsgTopup(
			generatedAddress1,
			generatedAddress1,
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

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Events), 0, "Appropriate error should be emitted for successful post-tx")

		// there should be stored sequence
		ok := initApp.TopupKeeper.HasTopupSequence(ctx, sequence.String())
		require.True(t, ok)

		result, err = suite.postHandler(ctx, &msg, abci.SideTxResultType_YES)
		require.Error(t, err)
		require.Nil(t, result)
	})
}
