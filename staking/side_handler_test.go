package staking_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
	suite.sideHandler = staking.NewSideTxHandler(suite.app.StakingKeeper, &suite.contractCaller)
	suite.postHandler = staking.NewPostTxHandler(suite.app.StakingKeeper, &suite.contractCaller)

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

func (suite *SideHandlerTestSuite) TestSideHandleMsgValidatorJoin() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	validatorId := uint64(1)
	amount, _ := big.NewInt(0).SetString("1000000000000000000", 10)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	chainParams := app.ChainKeeper.GetParams(ctx)
	blockNumber := big.NewInt(10)
	nonce := big.NewInt(3)

	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")
	})

	suite.Run("No receipt", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeWaitFrConfirmation), result.Code)
	})

	suite.Run("No EventLog", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(nil, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeErrDecodeEvent), result.Code)
	})

	suite.Run("Invalid Signer pubkey", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			hmTypes.NewPubKey([]byte{234}),
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid Signer address", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          hmTypes.ZeroHeimdallAddress.EthAddress(),
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid Validator Id", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			uint64(10),
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     big.NewInt(1),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid Activation Epoch", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(10),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid Amount", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(100000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid Block Number", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			uint64(20),
			nonce.Uint64(),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})

	suite.Run("Invalid nonce", func() {
		suite.contractCaller = mocks.IContractCaller{}
		txreceipt := &ethTypes.Receipt{
			BlockNumber: blockNumber,
		}

		msgValJoin := types.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(address.Bytes()),
			validatorId,
			uint64(1),
			sdk.NewInt(1000000000000000000),
			pubkey,
			txHash,
			logIndex,
			blockNumber.Uint64(),
			uint64(9),
		)

		stakinginfoStaked := &stakinginfo.StakinginfoStaked{
			Signer:          address,
			ValidatorId:     new(big.Int).SetUint64(validatorId),
			Nonce:           nonce,
			ActivationEpoch: big.NewInt(1),
			Amount:          amount,
			Total:           big.NewInt(10),
			SignerPubkey:    pubkey.Bytes()[1:],
		}

		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

		suite.contractCaller.On("DecodeValidatorJoinEvent", chainParams.ChainParams.StakingInfoAddress.EthAddress(), txreceipt, msgValJoin.LogIndex).Return(stakinginfoStaked, nil)

		result := suite.sideHandler(ctx, msgValJoin)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidMsg), result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgSingerUpdate() {
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgValidatorExit() {
}

func (suite *SideHandlerTestSuite) TestHandleMsgStakeUpdate() {
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// post tx handler
	result := suite.postHandler(ctx, nil, abci.SideTxResultType_Yes)
	require.False(t, result.IsOK(), "Post handler should fail")
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func (suite *SideHandlerTestSuite) TestpostHandleMsgValidatorJoin() {
}

func (suite *SideHandlerTestSuite) TestpostHandleMsgSingerUpdate() {
}

func (suite *SideHandlerTestSuite) TestpostHandleMsgValidatorExit() {
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgStakeUpdate() {
}
