package staking_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/staking/types"
	cmn "github.com/maticnetwork/heimdall/test"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	querier        sdk.Querier
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier tesing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = staking.NewQuerier(suite.app.StakingKeeper, &suite.contractCaller)
}

// TestQuerierTestSuite
func TestQuerierTestSuite(t *testing.T) {
	suite.Run(t, new(QuerierTestSuite))
}

// TestInvalidQuery checks request query
func (suite *QuerierTestSuite) TestInvalidQuery() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, bz)

	bz, err = querier(ctx, []string{types.QuerierRoute}, req)
	require.Error(t, err)
	require.Nil(t, bz)
}

func (suite *QuerierTestSuite) TestHandleQueryCurrentValidatorSet() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	path := []string{types.QueryCurrentValidatorSet}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentValidatorSet)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TesthandleQuerySigner() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	validators := keeper.GetAllValidators(ctx)
	path := []string{types.QuerySigner}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigner)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQuerySignerParams(validators[0].Signer.Bytes())),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TesthandleQueryValidator() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetAllValidators(ctx)

	path := []string{types.QueryValidator}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidator)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryValidatorParams(validators[0].ID)),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryValidatorStatus() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetAllValidators(ctx)

	path := []string{types.QueryValidatorStatus}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorStatus)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQuerySignerParams(validators[0].Signer.Bytes())),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryProposer() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	path := []string{types.QueryProposer}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProposer)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryProposerParams(uint64(2))),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryCurrentProposer() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	path := []string{types.QueryCurrentProposer}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentProposer)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryDividendAccount() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	var divAcc hmTypes.DividendAccount
	path := []string{types.QueryDividendAccount}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccount)
	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.NewDividendAccountID(uint64(1)),
		big.NewInt(0).String(),
		big.NewInt(0).String(),
	)
	app.StakingKeeper.AddDividendAccount(ctx, dividendAccount)
	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryDividendAccountParams(dividendAccount.ID)),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
	json.Unmarshal(res, &divAcc)
	require.Equal(t, dividendAccount, divAcc)
}

func (suite *QuerierTestSuite) TestHandleDividendAccountRoot() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	keeper := app.StakingKeeper
	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.NewDividendAccountID(uint64(1)),
		big.NewInt(0).String(),
		big.NewInt(0).String(),
	)
	keeper.AddDividendAccount(ctx, dividendAccount)

	path := []string{types.QueryDividendAccountRoot}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccountRoot)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryAccountProof() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	var accountRoot [32]byte

	path := []string{types.QueryAccountProof}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountProof)
	stakingInfo := &stakinginfo.Stakinginfo{}

	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.NewDividendAccountID(uint64(1)),
		big.NewInt(0).String(),
		big.NewInt(0).String(),
	)
	app.StakingKeeper.AddDividendAccount(ctx, dividendAccount)
	dividendAccounts := app.StakingKeeper.GetAllDividendAccounts(ctx)

	accRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)
	copy(accountRoot[:], accRoot)
	suite.contractCaller.On("GetStakingInfoInstance", mock.Anything).Return(stakingInfo, nil)

	suite.contractCaller.On("CurrentAccountStateRoot", stakingInfo).Return(accountRoot, nil)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(dividendAccount),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryVerifyAccountProof() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.NewDividendAccountID(uint64(1)),
		big.NewInt(0).String(),
		big.NewInt(0).String(),
	)
	app.StakingKeeper.AddDividendAccount(ctx, dividendAccount)

	path := []string{types.QueryVerifyAccountProof}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVerifyAccountProof)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(dividendAccount),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
	require.Equal(t, "true", string(res))
}

func (suite *QuerierTestSuite) TestHandleQueryStakingSequence() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	chainParams := app.ChainKeeper.GetParams(ctx)

	txHash := hmTypes.HexToHeimdallHash("123")

	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}

	logIndex := uint64(simulation.RandIntBetween(r1, 0, 100))

	msg := types.NewQueryStakingSequenceParams(txHash.String(), logIndex)

	sequence := new(big.Int).Mul(txreceipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	app.StakingKeeper.SetStakingSequence(ctx, sequence.String())

	suite.contractCaller.On("GetConfirmedTxReceipt", mock.Anything, txHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	path := []string{types.QueryStakingSequence}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryStakingSequence)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(msg),
	}

	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
	require.Equal(t, sequence.String(), string(res))
}
