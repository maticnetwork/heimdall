package staking_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/types/simulation"

	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"
	"github.com/stretchr/testify/mock"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/topup"

	"github.com/maticnetwork/heimdall/contracts/stakinginfo"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
	checkPointSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/x/staking"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"

	"github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         client.Context
	handler        sdk.Handler
	topUpHandler   sdk.Handler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = staking.NewHandler(suite.app.StakingKeeper, &suite.contractCaller)
	suite.topUpHandler = topup.NewHandler(suite.app.TopupKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorJoin() {

	t, initApp, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	// keys and addresses
	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmCommon.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	// loading the validators
	stakingSim.LoadValidatorSet(4, t, initApp.StakingKeeper, ctx, false, 10)

	validatorId := r.Uint64()
	logIndex := r.Uint64()
	activationEpoch := r.Uint64()
	blockNumber := r.Uint64()
	nonce := r.Uint64()
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	txHash := hmCommon.HexToHeimdallHash("123")
	chainParams := initApp.ChainKeeper.GetParams(ctx)

	txReceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	msgValJoin := types.NewMsgValidatorJoin(
		address.Bytes(),
		validatorId,
		activationEpoch,
		sdk.NewInt(1000000000000000000),
		pubkey,
		txHash,
		logIndex,
		blockNumber,
		nonce,
	)
	//
	stakingInfoStaked := &stakinginfo.StakinginfoStaked{
		Signer:          address,
		ValidatorId:     new(big.Int).SetUint64(validatorId),
		ActivationEpoch: big.NewInt(1),
		Amount:          amount,
		Total:           big.NewInt(10),
		SignerPubkey:    pubkey.Bytes()[1:],
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

	result, err := suite.handler(ctx, &msgValJoin)
	require.NotNil(t, result, "expected validator join to be ok, got %v", result)
	require.NoError(t, err)

	actualResult, ok := initApp.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
	require.Equal(t, false, ok, "Should not add validator")
	require.NotNil(t, actualResult, "got %v", actualResult)
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := suite.app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	// vals := oldValSet.(*Validators)
	oldSigner := oldValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	t.Log("To be Updated ===>", "Validator", newSigner[0].String())
	chainParams := initApp.ChainKeeper.GetParams(ctx)

	// gen msg
	msgTxHash := hmCommon.HexToHeimdallHash("123")
	msg := types.NewMsgSignerUpdate(sdk.AccAddress(newSigner[0].Signer), uint64(newSigner[0].ID), hmCommon.NewPubKeyFromHex(newSigner[0].PubKey), msgTxHash, 0, 0, 1)

	txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
		ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
		OldSigner:    hmCommon.HexToHeimdallAddress(oldSigner.Signer).EthAddress(),
		NewSigner:    hmCommon.HexToHeimdallAddress(newSigner[0].Signer).EthAddress(),
		SignerPubkey: hmCommon.PubKey(newSigner[0].PubKey).Bytes()[1:],
	}
	suite.contractCaller.On("DecodeSignerUpdateEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, uint64(0)).Return(signerUpdateEvent, nil)

	result, err := suite.handler(ctx, &msg)
	require.NotNil(t, result, "expected validator update to be ok, got %v", result)
	require.NoError(t, err)
	newValidators := keeper.GetCurrentValidators(ctx)
	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")

	setUpdates := helper.GetUpdatedValidators(oldValSet, keeper.GetAllValidators(ctx), 5)
	err = oldValSet.UpdateWithChangeSet(setUpdates)
	require.NoError(t, err)
	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
	require.True(t, ok, "signer should be found, got %v", ok)
	require.NotEqual(t, oldSigner.Signer, newSigner[0].Signer, "Should not update state")
	require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)

	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.GetSigner())
	require.Empty(t, err)
	require.NotEqual(t, removedVal.VotingPower, int64(0), "should not update state")
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorExit() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, suite.app.StakingKeeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	fmt.Printf("Validatros %+v\n", validators)

	msgTxHash := hmCommon.HexToHeimdallHash("123")
	chainParams := intiApp.ChainKeeper.GetParams(ctx)
	logIndex := uint64(0)

	txReceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	stakingInfoUnStakeInit := &stakinginfo.StakinginfoUnstakeInit{
		User:              hmCommon.HexToHeimdallAddress(validators[0].GetSigner().String()).EthAddress(),
		ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
		DeactivationEpoch: big.NewInt(10),
		Amount:            amount,
	}
	validators[0].EndEpoch = 10

	suite.contractCaller.On("DecodeValidatorExitEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, logIndex).Return(stakingInfoUnStakeInit, nil)

	msg := types.NewMsgValidatorExit(validators[0].GetSigner(), uint64(validators[0].ID), validators[0].EndEpoch, msgTxHash, 0, 0, 1)

	result, err := suite.handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, result, "expected validator exit to be ok, got %v", result)

	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].GetSigner())
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].GetSigner(), err)
	require.NotEqual(t, updatedValInfo.EndEpoch, validators[0].EndEpoch, "should not update deactivation epoch")

	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")

	result, err = suite.handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, result, "should not fail, as state is not updated for validatorExit")
}

func (suite *HandlerTestSuite) TestHandleMsgStakeUpdate() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper

	// pass 0 as time alive to generate non de-activated validators
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	t.Log("To be Updated ===>", "Validator", oldVal.String())
	chainParams := initApp.ChainKeeper.GetParams(ctx)

	msgTxHash := hmCommon.HexToHeimdallHash("123")
	msg := types.NewMsgStakeUpdate(oldVal.GetSigner(), oldVal.ID.Uint64(), sdk.NewInt(2000000000000000000), msgTxHash, 0, 0, 1)

	txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	stakingInfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
		ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
		NewAmount:   new(big.Int).SetInt64(2000000000000000000),
	}

	suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, uint64(0)).Return(stakingInfoStakeUpdate, nil)

	result, err := suite.handler(ctx, &msg)
	require.NotNil(t, result, "expected validator stake update to be ok, got %v", result)
	require.NoError(t, err)
	updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.GetSigner())
	require.Empty(t, err, "unable to fetch validator info %v-", err)
	require.NotEqual(t, stakingInfoStakeUpdate.NewAmount.Int64(), updatedVal.VotingPower, "Validator VotingPower should not be updated to %v", stakingInfoStakeUpdate.NewAmount.Uint64())
}

func (suite *HandlerTestSuite) TestExitedValidatorJoiningAgain() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	accounts := simulation.RandomAccounts(r1, 1)
	pubKey := hmCommon.NewPubKey(accounts[0].PubKey.Bytes())

	signerAddress, err := sdk.AccAddressFromHex(pubKey.Address().Hex())
	require.NoError(t, err)

	commonAddr := common.BytesToAddress(signerAddress.Bytes())

	txHash := hmCommon.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

	initApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 20)

	validatorId := hmTypes.NewValidatorID(uint64(1))
	validator := hmTypes.NewValidator(
		validatorId,
		10,
		15,
		1,
		int64(0), // power
		pubKey,
		signerAddress,
	)

	err = initApp.StakingKeeper.AddValidator(ctx, *validator)
	if err != nil {
		t.Error("Error while adding validator to store", err)
	}
	require.NoError(t, err)

	isCurrentValidator := validator.IsCurrentValidator(14)
	require.False(t, isCurrentValidator)

	totalValidators := initApp.StakingKeeper.GetAllValidators(ctx)
	require.Equal(t, totalValidators[0].Signer, signerAddress.String())

	chainParams := initApp.ChainKeeper.GetParams(ctx)

	txReceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}
	msgValJoin := types.NewMsgValidatorJoin(
		signerAddress,
		validatorId.Uint64(),
		uint64(1),
		sdk.NewInt(100000),
		pubKey,
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
		SignerPubkey:    hmCommon.PubKey(pubKey.Bytes())[1:],
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

	result, err := suite.handler(ctx, &msgValJoin)
	require.Error(t, err)
	require.Nil(t, result)
}

func (suite *HandlerTestSuite) TestTopupSuccessBeforeValidatorJoin() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	privateKey := secp256k1.GenPrivKey()
	pubKey := hmCommon.NewPubKey(privateKey.PubKey().Bytes())
	signer := sdk.AccAddress(privateKey.PubKey().Address().Bytes())

	commonAddr := common.BytesToAddress(signer.Bytes())

	txHash := hmCommon.HexToHeimdallHash("123")
	logIndex := uint64(2)
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

	validatorId := hmTypes.NewValidatorID(uint64(1))

	chainParams := initApp.ChainKeeper.GetParams(ctx)

	msgTopUp := topupTypes.NewMsgTopup(signer.Bytes(), signer.Bytes(), sdk.NewInt(2000000000000000000), txHash, logIndex, uint64(2))

	stakingInfoTopUpFee := &stakinginfo.StakinginfoTopUpFee{
		User: commonAddr,
		Fee:  big.NewInt(100000000000000000),
	}

	txReceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	stakingInfoStaked := &stakinginfo.StakinginfoStaked{
		Signer:          commonAddr,
		ValidatorId:     new(big.Int).SetUint64(validatorId.Uint64()),
		ActivationEpoch: big.NewInt(1),
		Amount:          amount,
		Total:           big.NewInt(10),
		SignerPubkey:    pubKey.Bytes()[1:],
	}

	msgValJoin := types.NewMsgValidatorJoin(
		signer.Bytes(),
		validatorId.Uint64(),
		uint64(1),
		sdk.NewInt(2000000000000000000),
		pubKey,
		txHash,
		logIndex,
		0,
		1,
	)

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash, chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress, txReceipt, msgValJoin.LogIndex).Return(stakingInfoStaked, nil)

	suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StakingInfoAddress, mock.Anything, msgTopUp.LogIndex).Return(stakingInfoTopUpFee, nil)

	topUpResult, err := suite.topUpHandler(ctx, &msgTopUp)
	require.NoError(t, err)
	require.NotNil(t, topUpResult, "expected topup to be done, got %v", topUpResult)

	result, err := suite.handler(ctx, &msgValJoin)
	require.NoError(t, err)
	require.NotNil(t, result, "expected validator stake update to be ok, got %v", result)

}
