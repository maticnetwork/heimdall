package topup_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite integrate test suite context object
type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	querier        sdk.Querier
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier tesing
func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = topup.NewHandler(suite.app.TopupKeeper, &suite.contractCaller)
}

// TestHandlerTestSuite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgTopup() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)

	pAddress := hmTypes.HexToHeimdallAddress("123")
	validatorId := uint64(simulation.RandIntBetween(r1, 0, 100))
	blockNumber := big.NewInt(10)
	chainParams := app.ChainKeeper.GetParams(ctx)
	txreceipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
	}

	msgTopup := types.NewMsgTopup(pAddress, validatorId, pAddress, sdk.NewInt(1000000000000000000), txHash, logIndex, blockNumber.Uint64())

	stakinginfoTopUpFee := &stakinginfo.StakinginfoTopUpFee{
		ValidatorId: new(big.Int).SetUint64(validatorId),
		Signer:      pAddress.EthAddress(),
		Fee:         big.NewInt(100000000000000000),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), mock.Anything, msgTopup.LogIndex).Return(stakinginfoTopUpFee, nil)
	result := suite.handler(ctx, msgTopup)
	require.True(t, result.IsOK(), "expected topup to be done, got %v", result)
}

func (suite *HandlerTestSuite) TestHandleMsgWithdrawFee() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	amount, _ := big.NewInt(0).SetString("0", 10)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	txHash := hmTypes.HexToHeimdallHash("123")
	logIndex := simulation.RandIntBetween(r1, 0, 100)

	validatorId := uint64(simulation.RandIntBetween(r1, 0, 100))

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	validatorAddress := pubkey.Address()

	chainParams := app.ChainKeeper.GetParams(ctx)
	blockNumber := big.NewInt(10)
	txreceipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
	}
	signer := hmTypes.BytesToHeimdallAddress(validatorAddress.Bytes())
	msgTopup := types.NewMsgTopup(signer, validatorId, signer, sdk.NewInt(1000000000000000000), txHash, uint64(logIndex), blockNumber.Uint64())

	stakinginfoTopUpFee := &stakinginfo.StakinginfoTopUpFee{
		ValidatorId: new(big.Int).SetUint64(validatorId),
		Signer:      validatorAddress,
		Fee:         big.NewInt(100000000000000000),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	suite.contractCaller.On("DecodeValidatorTopupFeesEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), mock.Anything, msgTopup.LogIndex).Return(stakinginfoTopUpFee, nil)
	topupResult := suite.handler(ctx, msgTopup)

	require.True(t, topupResult.IsOK(), "expected topup to be done, got %v", topupResult)

	// start Withdraw fees
	startBlock := uint64(simulation.RandIntBetween(r1, 1, 100))

	power := simulation.RandIntBetween(r1, 1, 100)

	timeAlive := uint64(10)

	newVal := hmTypes.Validator{
		ID:               hmTypes.NewValidatorID(validatorId),
		StartEpoch:       startBlock,
		EndEpoch:         startBlock + timeAlive,
		VotingPower:      int64(power),
		Signer:           hmTypes.HexToHeimdallAddress(pubkey.Address().String()),
		PubKey:           pubkey,
		ProposerPriority: 0,
	}
	app.StakingKeeper.AddValidator(ctx, newVal)

	msgWithdrawFee := types.NewMsgWithdrawFee(
		hmTypes.BytesToHeimdallAddress(validatorAddress.Bytes()),
		sdk.NewIntFromBigInt(amount),
	)
	withdrawResult := suite.handler(ctx, msgWithdrawFee)
	require.True(t, withdrawResult.IsOK(), "expected withdraw tobe done, got %v", withdrawResult)
}
