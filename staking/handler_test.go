package staking_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	errs "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/maticnetwork/heimdall/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app    *app.HeimdallApp
	ctx    sdk.Context
	cliCtx context.CLIContext

	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = staking.NewHandler(suite.app.StakingKeeper, &suite.contractCaller)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorJoin() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	validatorId := uint64(1)
	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	chainParams := app.ChainKeeper.GetParams(ctx)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	msgValJoin := stakingTypes.NewMsgValidatorJoin(
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		validatorId,
		pubkey,
		txHash,
		logIndex,
	)

	stakinginfoStaked := &stakinginfo.StakinginfoStaked{
		Signer:          address,
		ValidatorId:     new(big.Int).SetUint64(validatorId),
		ActivationEpoch: big.NewInt(1),
		Amount:          amount,
		Total:           big.NewInt(10),
		SignerPubkey:    pubkey.Bytes()[1:],
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", mock.Anything, txHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

	result := suite.handler(ctx, msgValJoin)
	require.True(t, result.IsOK(), "expected validator join to be ok, got %v", result)

	actualResult, ok := suite.app.StakingKeeper.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorId))
	require.Equal(t, ok, true)
	require.NotNil(t, actualResult)
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorUpdate() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := suite.app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)

	// vals := oldValSet.(*Validators)
	oldSigner := oldValSet.Validators[0]
	newSigner := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
	newSigner[0].ID = oldSigner.ID
	newSigner[0].VotingPower = oldSigner.VotingPower
	t.Log("To be Updated ===>", "Validator", newSigner[0].String())
	chainParams := app.ChainKeeper.GetParams(ctx)

	// gen msg
	msgTxHash := types.HexToHeimdallHash("123")
	msg := stakingTypes.NewMsgSignerUpdate(newSigner[0].Signer, uint64(newSigner[0].ID), newSigner[0].PubKey, msgTxHash, 0)

	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", mock.Anything, msgTxHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	signerUpdateEvent := &stakinginfo.StakinginfoSignerChange{
		ValidatorId:  new(big.Int).SetUint64(oldSigner.ID.Uint64()),
		OldSigner:    oldSigner.Signer.EthAddress(),
		NewSigner:    newSigner[0].Signer.EthAddress(),
		SignerPubkey: newSigner[0].PubKey.Bytes()[1:],
	}
	suite.contractCaller.On("DecodeSignerUpdateEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, uint64(0)).Return(signerUpdateEvent, nil)

	result := suite.handler(ctx, msg)

	require.True(t, result.IsOK(), "expected validator update to be ok, got %v", result)
	newValidators := keeper.GetCurrentValidators(ctx)
	require.Equal(t, len(oldValSet.Validators), len(newValidators), "Number of current validators should be equal")

	setUpdates := helper.GetUpdatedValidators(&oldValSet, keeper.GetAllValidators(ctx), 5)
	oldValSet.UpdateWithChangeSet(setUpdates)
	_ = keeper.UpdateValidatorSetInStore(ctx, oldValSet)

	ValFrmID, ok := keeper.GetValidatorFromValID(ctx, oldSigner.ID)
	require.True(t, ok, "new signer should be found, got %v", ok)
	require.Equal(t, ValFrmID.Signer.Bytes(), newSigner[0].Signer.Bytes(), "New Signer should be mapped to old validator ID")
	require.Equal(t, ValFrmID.VotingPower, oldSigner.VotingPower, "VotingPower of new signer %v should be equal to old signer %v", ValFrmID.VotingPower, oldSigner.VotingPower)

	removedVal, err := keeper.GetValidatorInfo(ctx, oldSigner.Signer.Bytes())
	require.Empty(t, err, "deleted validator should be found, got %v", err)
	require.Equal(t, removedVal.VotingPower, int64(0), "removed validator VotingPower should be zero")
	t.Log("Deleted validator ===>", "Validator", removedVal.String())
}

func (suite *HandlerTestSuite) TestHandleMsgValidatorExit() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	// pass 0 as time alive to generate non de-activated validators
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	validators := keeper.GetCurrentValidators(ctx)
	msgTxHash := types.HexToHeimdallHash("123")
	chainParams := app.ChainKeeper.GetParams(ctx)
	logIndex := uint64(0)

	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", mock.Anything, msgTxHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	amount, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	stakinginfoUnstakeInit := &stakinginfo.StakinginfoUnstakeInit{
		User:              validators[0].Signer.EthAddress(),
		ValidatorId:       big.NewInt(0).SetUint64(validators[0].ID.Uint64()),
		DeactivationEpoch: big.NewInt(10),
		Amount:            amount,
	}

	suite.contractCaller.On("DecodeValidatorExitEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, logIndex).Return(stakinginfoUnstakeInit, nil)

	validators[0].EndEpoch = 10
	msg := stakingTypes.NewMsgValidatorExit(validators[0].Signer, uint64(validators[0].ID), msgTxHash, 0)

	got := suite.handler(ctx, msg)

	require.True(t, got.IsOK(), "expected validator exit to be ok, got %v", got)

	updatedValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	require.Empty(t, err, "Unable to get validator info from val address,ValAddr:%v Error:%v ", validators[0].Signer.String(), err)
	require.Equal(t, updatedValInfo.EndEpoch, validators[0].EndEpoch, "deactivation epoch should be set correctly")

	_, found := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.True(t, found, "Validator should be present even after deactivation")

	got = suite.handler(ctx, msg)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	currentVals := keeper.GetCurrentValidators(ctx)
	require.Equal(t, 4, len(currentVals), "No of current validators should exist before epoch passes")

	app.CheckpointKeeper.UpdateACKCountWithValue(ctx, 20)
	currentVals = keeper.GetCurrentValidators(ctx)
	require.Equal(t, 3, len(currentVals), "No of current validators should reduce after epoch passes")
}

func (suite *HandlerTestSuite) TestHandleMsgStakeUpdate() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper

	// pass 0 as time alive to generate non de-activated validators
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 0)
	oldValSet := keeper.GetValidatorSet(ctx)
	oldVal := oldValSet.Validators[0]

	t.Log("To be Updated ===>", "Validator", oldVal.String())

	chainParams := app.ChainKeeper.GetParams(ctx)

	msgTxHash := types.HexToHeimdallHash("123")
	msg := stakingTypes.NewMsgStakeUpdate(oldVal.Signer, oldVal.ID.Uint64(), msgTxHash, 0)

	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	suite.contractCaller.On("GetConfirmedTxReceipt", mock.Anything, msgTxHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	stakinginfoStakeUpdate := &stakinginfo.StakinginfoStakeUpdate{
		ValidatorId: new(big.Int).SetUint64(oldVal.ID.Uint64()),
		NewAmount:   new(big.Int).SetInt64(2000000000000000000),
	}

	suite.contractCaller.On("DecodeValidatorStakeUpdateEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, uint64(0)).Return(stakinginfoStakeUpdate, nil)

	got := suite.handler(ctx, msg)
	require.True(t, got.IsOK(), "expected validator stake update to be ok, got %v", got)
	updatedVal, err := keeper.GetValidatorInfo(ctx, oldVal.Signer.Bytes())
	require.Empty(t, err, "unable to fetch validator info %v-", err)
	require.Equal(t, stakinginfoStakeUpdate.NewAmount.Int64(), updatedVal.VotingPower, "Validator VotingPower should be updated to %v", stakinginfoStakeUpdate.NewAmount.Uint64())

}
