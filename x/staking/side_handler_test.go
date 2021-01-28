package staking_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/maticnetwork/bor/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"
	"github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

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
	suite.sideHandler = staking.NewSideTxHandler(suite.app.StakingKeeper, &suite.contractCaller)
	suite.postHandler = staking.NewPostTxHandler(suite.app.StakingKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SideHandlerTestSuite))
}

// Test Cases

func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, result.Code, uint32(6))
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgValidatorJoin() {
	t, initApp, ctx, r := suite.T(), suite.app, suite.ctx, suite.r
	txHash := hmCommonTypes.HexToHeimdallHash("123")
	index := r.Uint64()
	logIndex := index
	validatorId := hmTypes.NewValidatorID(uint64(1))
	amount, _ := big.NewInt(0).SetString("1000000000000000000", 10)

	// keys and addresses
	privateKey := secp256k1.GenPrivKey()
	pubkey := hmCommonTypes.NewPubKey(privateKey.PubKey().Bytes())
	address := sdk.AccAddress(privateKey.PubKey().Address().Bytes())
	commonAddr := common.BytesToAddress(address.Bytes())

	chainParams := initApp.ChainKeeper.GetParams(ctx)
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(3)

	// todo: invalid public key
	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			address,
			validatorId.Uint64(),
			uint64(1),
			sdk.NewInt(2000000000000000000),
			pubkey,
			txHash,
			logIndex,
			0,
			1,
		)

		stakingInfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          commonAddr,
			ValidatorId:     new(big.Int).SetUint64(validatorId.Uint64()),
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		addr, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorJoinEvent", addr, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

		result := suite.sideHandler(ctx, &msgValJoin)
		fmt.Printf("result %+v\n", result)
		//require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
		//require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")
	})

	suite.Run("No receipt", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId.Uint64(),
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakingInfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          commonAddr,
			ValidatorId:     new(big.Int).SetUint64(validatorId.Uint64()),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    hmCommonTypes.NewPubKey(pubkey.Bytes())[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)

		addr, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorJoinEvent", addr, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

		result := suite.sideHandler(ctx, &msgValJoin)
		fmt.Printf("result %+v\n", result)

		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should Fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		//require.Equal(t, uint32(common.CodeWaitFrConfirmation), result.Code)
	})
	//
	suite.Run("No EventLog", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId.Uint64(),
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorJoinEvent", stakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(nil, nil)

		result := suite.sideHandler(ctx, &msgValJoin)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should Fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
		//require.Equal(t, uint32(common.CodeErrDecodeEvent), result.Code)
	})
	// todo: panic invalid public key
	//suite.Run("Invalid Signer pubkey", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		validatorId.Uint64(),
	//		uint64(1),
	//		sdk.NewInt(1000000000000000000),
	//		hmCommonTypes.NewPubKey([]byte{234}),
	//		txHash,
	//		logIndex,
	//		blockNumber.Uint64(),
	//		nonce.Uint64(),
	//	)
	//
	//	stakingInfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     new(big.Int).SetUint64(validatorId.Uint64()),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    hmCommonTypes.NewPubKey(pubkey.Bytes())[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//	stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
	//	require.NoError(t, err)
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", stakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	require.NotEqual(t, uint32(0), result.Code, "Side tx handler should Fail")
	//	require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
	//
	suite.Run("Invalid Signer address", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId.Uint64(),
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakingInfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          hmCommonTypes.ZeroHeimdallAddress.EthAddress(),
			ValidatorId:     new(big.Int).SetUint64(validatorId.Uint64()),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    hmCommonTypes.PubKey(pubkey.Bytes())[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorJoinEvent", stakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

		result := suite.sideHandler(ctx, &msgValJoin)
		fmt.Printf("result %+v\n", result)

		//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})
	//
	//suite.Run("Invalid Validator Id", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		uint64(10),
	//		uint64(1),
	//		sdk.NewInt(1000000000000000000),
	//		pubkey,
	//		txHash,
	//		logIndex,
	//		blockNumber.Uint64(),
	//		nonce.Uint64(),
	//	)
	//
	//	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     big.NewInt(1),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    pubkey.Bytes()[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
	//	//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
	//
	//suite.Run("Invalid Activation Epoch", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		validatorId,
	//		uint64(10),
	//		sdk.NewInt(1000000000000000000),
	//		pubkey,
	//		txHash,
	//		logIndex,
	//		blockNumber.Uint64(),
	//		nonce.Uint64(),
	//	)
	//
	//	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     new(big.Int).SetUint64(validatorId),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    pubkey.Bytes()[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//	//suite.contractCaller.On("DecodeValidatorJoinEvent",stakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
	//	//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
	//
	//suite.Run("Invalid Amount", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		validatorId,
	//		uint64(1),
	//		sdk.NewInt(100000000000000000),
	//		pubkey,
	//		txHash,
	//		logIndex,
	//		blockNumber.Uint64(),
	//		nonce.Uint64(),
	//	)
	//
	//	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     new(big.Int).SetUint64(validatorId),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    pubkey.Bytes()[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
	//	//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
	//
	//suite.Run("Invalid Block Number", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		validatorId,
	//		uint64(1),
	//		sdk.NewInt(1000000000000000000),
	//		pubkey,
	//		txHash,
	//		logIndex,
	//		uint64(20),
	//		nonce.Uint64(),
	//	)
	//
	//	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     new(big.Int).SetUint64(validatorId),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    pubkey.Bytes()[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
	//	//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
	//
	//suite.Run("Invalid nonce", func() {
	//	suite.contractCaller = mocks.IContractCaller{}
	//	txReceipt := &ethTypes.Receipt{
	//		BlockNumber: blockNumber,
	//	}
	//
	//	msgValJoin := types.NewMsgValidatorJoin(
	//		address.Bytes(),
	//		validatorId,
	//		uint64(1),
	//		sdk.NewInt(1000000000000000000),
	//		pubkey,
	//		txHash,
	//		logIndex,
	//		blockNumber.Uint64(),
	//		uint64(9),
	//	)
	//
	//	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//		Signer:          commonAddr,
	//		ValidatorId:     new(big.Int).SetUint64(validatorId),
	//		Nonce:           nonce,
	//		ActivationEpoch: big.NewInt(1),
	//		Amount:          amount,
	//		Total:           big.NewInt(10),
	//		SignerPubkey:    pubkey.Bytes()[1:],
	//	}
	//
	//	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
	//
	//	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)
	//
	//	result := suite.sideHandler(ctx, &msgValJoin)
	//	fmt.Printf("result %+v\n", result)
	//
	//	//require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
	//	//require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
	//	//require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	//})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgValidatorExit() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")
	chainParams := initApp.ChainKeeper.GetParams(ctx)
	logIndex := uint64(0)
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(9)

	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}
		validators[0].EndEpoch = 10

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)
		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, tmprototypes.SideTxResultType_YES, result.Result, "Result should be `yes`")
	})

	suite.Run("No Receipt", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)

		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}
		validators[0].EndEpoch = 10

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("No Eventlog", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		validators[0].EndEpoch = 10

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(nil, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)

		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid BlockNumber", func() {
		suite.contractCaller = mocks.IContractCaller{}
		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}
		validators[0].EndEpoch = 10
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			uint64(5),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)

		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid validatorId", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}
		validators[0].EndEpoch = 10

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(66),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid DeactivationEpoch", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			uint64(1000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid Nonce", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txReceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
		stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
			User:              hmCommonTypes.HexToHeimdallAddress(validators[0].Signer).EthAddress(),
			ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
			Nonce:             nonce,
			DeactivationEpoch: big.NewInt(10),
			Amount:            amount,
		}
		validators[0].EndEpoch = 10
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorExitEvent", stakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			uint64(6),
		)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgSignerUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := suite.app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	oldSigner := oldValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	chainParams := initApp.ChainKeeper.GetParams(ctx)
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(5)

	// gen msg
	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")

	suite.Run("Success", func() {
		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID),
			hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeSignerUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, tmprototypes.SideTxResultType_YES, result.Result, "Result should be `yes`")
	})

	suite.Run("No Eventlog", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID), hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeSignerUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(nil, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})

	suite.Run("Invalid BlockNumber", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(
			sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID),
			hmCommonTypes.PubKey(newSigner[0].PubKey),
			msgTxHash,
			0,
			uint64(9),
			nonce.Uint64(),
		)

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeSignerUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})

	suite.Run("Invalid validator", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(6), hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeSignerUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})

	suite.Run("Invalid signer pubkey", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID), hmCommonTypes.NewPubKey([]byte{123}), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeSignerUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})

	suite.Run("Invalid new signer address", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(sdk.AccAddress(hmCommonTypes.ZeroHeimdallAddress.String()), uint64(oldSigner.ID), hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.ZeroHeimdallAddress.EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeSignerUpdateEvent",
			stakingInfoAddress,
			txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})

	suite.Run("Invalid nonce", func() {
		suite.contractCaller = mocks.IContractCaller{}

		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID), hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), uint64(12))

		txReceipt := &ethTypes.Receipt{BlockNumber: blockNumber}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
			ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
			Nonce:        nonce,
			OldSigner:    hmCommonTypes.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
			NewSigner:    hmCommonTypes.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
			SignerPubkey: hmCommonTypes.PubKey(newSigner[0].PubKey).Bytes()[1:],
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeSignerUpdateEvent",
			stakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")

	})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgStakeUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper

	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	chainParams := initApp.ChainKeeper.GetParams(ctx)

	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(1)

	suite.Run("Success", func() {
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(),
			chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent",
			stakingInfoAddress,
			txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)

		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, tmprototypes.SideTxResultType_YES, result.Result, "Result should be `yes`")
	})

	suite.Run("No Receipt", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent",
			stakingInfoAddress,
			txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("No Eventlog", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}

		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent",
			stakingInfoAddress,
			txReceipt, uint64(0)).Return(nil, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid BlockNumber", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			uint64(15),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent",
			stakingInfoAddress,
			txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid ValidatorID", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			uint64(13),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}
		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", hmCommonTypes.HexToHeimdallAddress(
			chainParams.ChainParams.StakingInfoAddress).EthAddress(), txReceipt, uint64(0),
		).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid Amount", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(200000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}
		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)
		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})

	suite.Run("Invalid Nonce", func() {
		suite.contractCaller = mocks.IContractCaller{}
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			uint64(9))

		txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
		suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

		stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
			ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
			NewAmount:   new(big.Int).SetInt64(2000000000000000000),
			Nonce:       nonce,
		}

		stakingInfoAddress, err := sdk.AccAddressFromHex(chainParams.ChainParams.StakingInfoAddress)
		require.NoError(t, err)

		suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", stakingInfoAddress, txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

		result := suite.sideHandler(ctx, &msg)
		require.NotEqual(t, uint32(0), result.Code, "Side tx handler should fail")
		require.Equal(t, tmprototypes.SideTxResultType_SKIP, result.Result, "Result should skip")
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result, err := suite.postHandler(ctx, nil, tmprototypes.SideTxResultType_YES)
	require.Nil(t, result)
	require.Error(t, err)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgValidatorJoin() {
	t, initApp, ctx, r := suite.T(), suite.app, suite.ctx, suite.r
	txHash := hmCommonTypes.HexToHeimdallHash("123")
	logIndex := r.Uint64()
	validatorId := uint64(1)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmCommonTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	blockNumber := big.NewInt(10)
	nonce := big.NewInt(3)

	suite.Run("No Result", func() {

		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		_, err := suite.postHandler(ctx, &msgValJoin, tmprototypes.SideTxResultType_NO)
		require.Error(t, err)

		_, ok := initApp.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
		require.False(t, ok, "Should not add validator")
	})

	suite.Run("Success", func() {
		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result, err := suite.postHandler(ctx, &msgValJoin, tmprototypes.SideTxResultType_YES)
		require.NotNil(t, result)
		require.NoError(t, err, "expected validator join to be ok, got %v", result)

		actualResult, ok := initApp.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
		require.True(t, ok, "Should add validator")
		require.NotNil(t, actualResult, "got %v", actualResult)
	})

	suite.Run("Replay", func() {
		blockNumber := big.NewInt(11)

		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result, err := suite.postHandler(ctx, &msgValJoin, tmprototypes.SideTxResultType_YES)
		require.NoError(t, err, "expected validator join to be ok, got %v", result)
		require.NotNil(t, result)

		actualResult, ok := initApp.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
		require.True(t, ok, "Should add validator")
		require.NotNil(t, actualResult, "got %v", actualResult)

		result, err = suite.postHandler(ctx, &msgValJoin, tmprototypes.SideTxResultType_YES)
		require.Nil(t, result)
		require.Error(t, err)
	})

	suite.Run("Invalid Power", func() {
		msgValJoin := types.NewMsgValidatorJoin(
			address.Bytes(),
			validatorId,
			uint64(1),
			sdk.NewInt(1),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result, err := suite.postHandler(ctx, &msgValJoin, tmprototypes.SideTxResultType_YES)
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgSignerUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	oldSigner := oldValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(5)

	// gen msg
	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")

	suite.Run("No Success", func() {
		msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(oldSigner.ID), hmCommonTypes.PubKey(newSigner[0].PubKey), msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())
		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_NO)
		require.Error(t, err)
		require.Nil(t, result)
	})

	suite.Run("Success", func() {
		msg := types.NewMsgSignerUpdate(
			sdk.AccAddress(newSigner[0].Signer),
			uint64(oldSigner.ID),
			hmCommonTypes.NewPubKeyFromHex(newSigner[0].PubKey),
			msgTxHash, 0, blockNumber.Uint64(), nonce.Uint64())

		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_YES)
		require.NotNil(t, result)
		require.NoError(t, err, "Post handler should succeed")

		newValidators := keeper.GetCurrentValidators(ctx)
		require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")

		setUpdates := helper.GetUpdatedValidators(oldValSet, keeper.GetAllValidators(ctx), 5)
		err = oldValSet.UpdateWithChangeSet(setUpdates)
		require.NoError(t, err)
		_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

		ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
		require.True(t, ok, "new signer should be found, got %v", ok)
		require.Equal(t, ValFrmID.GetSigner().String(), newSigner[0].Signer, "New Signer should be mapped to old validator ID")
		require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)

		removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.GetSigner())
		require.Empty(t, err, "deleted validator should be found, got %v", err)
		require.Equal(t, removedVal.VotingPower, int64(0), "removed validator VotingPower should be zero")
	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgValidatorExit() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(9)

	suite.Run("No Success", func() {
		validators[0].EndEpoch = 10

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_NO)
		require.Error(t, err)
		require.Nil(t, result)
	})

	suite.Run("Success", func() {
		validators[0].EndEpoch = 10

		msg := types.NewMsgValidatorExit(
			sdk.AccAddress(validators[0].Signer),
			uint64(validators[0].ID),
			validators[0].EndEpoch,
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_YES)
		require.NoError(t, err)
		require.NotNil(t, result)

		currentVals := keeper.GetCurrentValidators(ctx)
		require.Equal(t, 4, len(currentVals), "No of current validators should exist before epoch passes")

		initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 20)
		currentVals = keeper.GetCurrentValidators(ctx)
		require.Equal(t, 3, len(currentVals), "No of current validators should reduce after epoch passes")
	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgStakeUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper

	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	msgTxHash := hmCommonTypes.HexToHeimdallHash("123")
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(1)
	newAmount := new(big.Int).SetInt64(2000000000000000000)

	suite.Run("No result", func() {
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_NO)
		require.Error(t, err)
		require.Nil(t, result)

		updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.GetSigner())
		require.Empty(t, err, "unable to fetch validator info %v-", err)

		actualPower, err := helper.GetPowerFromAmount(newAmount)
		require.NoError(t, err)
		require.NotEqual(t, actualPower.Int64(), updatedVal.VotingPower, "Validator VotingPower should be updated to %v", newAmount.Uint64())
	})

	suite.Run("Success", func() {
		msg := types.NewMsgStakeUpdate(
			sdk.AccAddress(oldVal.Signer),
			oldVal.ID.Uint64(),
			sdk.NewInt(2000000000000000000),
			msgTxHash,
			0,
			blockNumber.Uint64(),
			nonce.Uint64())

		result, err := suite.postHandler(ctx, &msg, tmprototypes.SideTxResultType_YES)
		require.NotNil(t, result)
		require.NoError(t, err, "expected validator stake update to be ok, got %v", result)

		updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.GetSigner())
		require.Empty(t, err, "unable to fetch validator info %v-", err)

		actualPower, err := helper.GetPowerFromAmount(new(big.Int).SetInt64(2000000000000000000))
		require.NoError(t, err)
		require.Equal(t, actualPower.Int64(), updatedVal.VotingPower, "Validator VotingPower should be updated to %v", newAmount.Uint64())
	})
}
