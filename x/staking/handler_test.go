package staking_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"
	"github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         client.Context
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = staking.NewHandler(suite.app.StakingKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorJoin() {

	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	// keys and addresses
	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmCommon.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	// loading the validators
	stakingSim.LoadValidatorSet(4, t, app.StakingKeeper, ctx, false, 10)
	//validators := app.StakingKeeper.GetAllValidators(ctx)
	//validator := validators[0]

	validatorId := r.Uint64()
	logIndex := r.Uint64()
	activationEpoch := r.Uint64()
	blockNumber := r.Uint64()
	nonce := r.Uint64()
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	txHash := hmCommon.HexToHeimdallHash("123")
	chainParams := app.ChainKeeper.GetParams(ctx)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	msgValJoin := types.NewMsgValidatorJoin(
		address.Bytes(),
		validatorId,
		activationEpoch,
		sdk.NewInt(amount.Int64()),
		pubkey,
		txHash,
		logIndex,
		blockNumber,
		nonce,
	)
	//
	//stakinginfoStaked := &stakinginfo.StakinginfoStaked{
	//	Signer:          address,
	//	ValidatorId:     new(big.Int).SetUint64(validatorId),
	//	ActivationEpoch: big.NewInt(1),
	//	Amount:          amount,
	//	Total:           big.NewInt(10),
	//	SignerPubkey:    pubkey.Bytes()[1:],
	//}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	//suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.Bytes(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

	result, err := suite.handler(ctx, &msgValJoin)
	fmt.Printf("Result is %v\n", result)
	fmt.Printf("Error ir %v\n", err)
	//require.True(t, result.IsOK(), "expected validator join to be ok, got %v", result)

	actualResult, ok := app.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
	require.Equal(t, false, ok, "Should not add validator")
	require.NotNil(t, actualResult, "got %v", actualResult)
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorExit() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, suite.app.StakingKeeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)
	fmt.Printf("Validatros %+v\n", validators)

	msgTxHash := hmCommon.HexToHeimdallHash("123")
	chainParams := app.ChainKeeper.GetParams(ctx)
	//logIndex := uint64(0)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", msgTxHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	//amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	//stakinginfoUnstakeInit := &stakinginfo.StakinginfoUnstakeInit{
	//	User:              hmCommon.HexToHeimdallAddress(validators[0].GetSigner().String()).EthAddress(),
	//	ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
	//	DeactivationEpoch: big.NewInt(10),
	//	Amount:            amount,
	//}
	validators[0].EndEpoch = 10

	//suite.contractCaller.On("DecodeValidatorExitEvent", chainParams.ChainParams.StakingInfoAddress, txreceipt, logIndex).Return(stakinginfoUnstakeInit, nil)

	msg := types.NewMsgValidatorExit(validators[0].GetSigner(), uint64(validators[0].ID), validators[0].EndEpoch, msgTxHash, 0, 0, 1)

	got, err := suite.handler(ctx, &msg)
	fmt.Printf("got %+v\n", got)
	fmt.Printf("error %+v\n", err)

	//require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)

	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].GetSigner())
	// updatedValInfo.EndEpoch = 10
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].GetSigner(), err)
	require.NotEqual(t, updatedValInfo.EndEpoch, validators[0].EndEpoch, "should not update deactivation epoch")

	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")

	got, err = suite.handler(ctx, &msg)
	fmt.Printf("got %+v\n", got)
	fmt.Printf("error %+v\n", err)
	//require.NotNil(t, err)
	//require.True(t, got.IsOK(), "should not fail, as state is not updated for validatorExit")
}

func (suite *HandlerTestSuite) TestHandlerMsgValidatorJoin() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r

	// keys and addresses
	_, testPubKey, addr1 := testdata.KeyTestPubAddr()
	// loading the validators
	stakingSim.LoadValidatorSet(4, t, app.StakingKeeper, ctx, false, 10)
	validators := app.StakingKeeper.GetAllValidators(ctx)
	validator := validators[0]

	validatorId := r.Uint64()
	logIndex := r.Uint64()
	activationEpoch := r.Uint64()
	blockNumber := r.Uint64()
	nonce := r.Uint64()
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	txHash := hmCommon.HexToHeimdallHash("123")
	pubKey := hmCommon.NewPubKeyFromHex(validator.PubKey)

	t.Parallel()
	t.Run("Success : Already Joined", func(t *testing.T) {
		msg := types.NewMsgValidatorJoin(validator.GetSigner(), validator.ID.Uint64(), activationEpoch, sdk.NewInt(amount.Int64()), pubKey, txHash, logIndex, blockNumber, nonce)
		_, err := suite.handler(ctx, &msg)
		require.NotNil(t, err)
	})

	t.Run("Success : Validator Joined", func(t *testing.T) {
		msg := types.NewMsgValidatorJoin(addr1, validatorId, activationEpoch, sdk.NewInt(amount.Int64()), testPubKey.Bytes(), txHash, logIndex, blockNumber, nonce)
		_, err := suite.handler(ctx, &msg)
		require.NotNil(t, err)
		actualResult, ok := app.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
		require.Equal(t, false, ok, "Should not add validator")
		require.NotNil(t, actualResult, "got %v", actualResult)
	})
}

func (suite *HandlerTestSuite) TestHandlerStakeUpdate() {
	t, app, ctx, r := suite.T(), suite.app, suite.ctx, suite.r
	// loading the validators
	stakingSim.LoadValidatorSet(4, t, app.StakingKeeper, ctx, false, 10)
	validators := app.StakingKeeper.GetAllValidators(ctx)
	validator := validators[0]

	logIndex := r.Uint64()
	blockNumber := r.Uint64()
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	txHash := hmCommon.HexToHeimdallHash("123")

	msg := types.NewMsgStakeUpdate(validator.GetSigner(), validator.ID.Uint64(), sdk.NewInt(amount.Int64()), txHash, logIndex, blockNumber, uint64(validator.Nonce+uint64(1)))

	t.Run("Success", func(t *testing.T) {
		result, err := suite.handler(ctx, &msg)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgSignerUpdate() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	// loading the validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	oldValSet := keeper.GetValidatorSet(ctx)

	// vals := oldValSet.(*Validators)
	oldSigner := oldValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	t.Log("To be Updated ===>", "Validator", newSigner[0].String())

	// gen msg
	msgTxHash := hmCommon.HexToHeimdallHash("123")
	msg := types.NewMsgSignerUpdate(
		newSigner[0].GetSigner(),
		uint64(newSigner[0].ID),
		hmCommon.NewPubKeyFromHex(newSigner[0].PubKey),
		msgTxHash, 0, 0, 1)

	result, err := suite.handler(ctx, &msg)

	require.NoError(t, err)
	require.NotNil(t, result)
	newValidators := keeper.GetCurrentValidators(ctx)
	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")

	setUpdates := helper.GetUpdatedValidators(oldValSet, keeper.GetAllValidators(ctx), 5)
	oldValSet.UpdateWithChangeSet(setUpdates)
	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
	require.True(t, ok, "signer should be found, got %v", ok)
	require.NotEqual(t, oldSigner.GetSigner(), newSigner[0].GetSigner(), "Should not update state")
	require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)

	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.GetSigner())
	require.Empty(t, err)
	require.NotEqual(t, removedVal.VotingPower, int64(0), "should not update state")
}

func (suite *HandlerTestSuite) TestHandlerValidatorExit() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	stakingSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := hmCommon.HexToHeimdallHash("123")

	validators[0].EndEpoch = 10
	msg := types.NewMsgValidatorExit(
		validators[0].GetSigner(), uint64(validators[0].ID),
		validators[0].EndEpoch, msgTxHash, 0, 0, 1)

	got, err := suite.handler(ctx, &msg)

	require.NoError(t, err)
	require.NotNil(t, got)
	//require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)

	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].GetSigner())
	// updatedValInfo.EndEpoch = 10
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].GetSigner(), err)
	require.NotEqual(t, updatedValInfo.EndEpoch, validators[0].EndEpoch, "should not update deactivation epoch")

	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")

	got, err = suite.handler(ctx, &msg)
	require.NoError(t, err)
	require.NotNil(t, got)
	//require.True(t, got.IsOK(), "should not fail, as state is not updated for validatorExit")
}
