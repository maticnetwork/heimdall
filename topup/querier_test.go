package topup_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/maticnetwork/heimdall/app"
	chainTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	querier        sdk.Querier
	contractCaller mocks.IContractCaller
	chainParams    chainTypes.Params
}

// SetupTest setup all necessary things for querier testing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = topup.NewQuerier(suite.app.TopupKeeper, &suite.contractCaller)
	suite.chainParams = suite.app.ChainKeeper.GetParams(suite.ctx)
}

// TestQuerierTestSuite
func TestQuerierTestSuite(t *testing.T) {
	t.Parallel()
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

// TestQuerySequence queries sequence data
func (suite *QuerierTestSuite) TestQuerySequence() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	chainParams := app.ChainKeeper.GetParams(ctx)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	txHash := hmTypes.HexToHeimdallHash("123")
	logIndex := uint64(simulation.RandIntBetween(r1, 0, 100))
	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}

	// set topup sequence
	sequence := new(big.Int).Mul(txreceipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))
	app.TopupKeeper.SetTopupSequence(ctx, sequence.String())

	// mock external calls
	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	path := []string{types.QuerySequence}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySequence)
	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQuerySequenceParams(txHash.String(), logIndex)),
	}

	// fetch sequence
	res, err := querier(ctx, path, req)
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
	require.Equal(t, sequence.String(), string(res))
}

func (suite *QuerierTestSuite) TestHandleQueryDividendAccount() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryDividendAccount}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccount)
	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.BytesToHeimdallAddress([]byte("some-address")),
		big.NewInt(0).String(),
	)
	err := app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryDividendAccountParams(dividendAccount.User)),
	}
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var divAcc hmTypes.DividendAccount
	err = jsoniter.ConfigFastest.Unmarshal(res, &divAcc)
	require.NoError(t, err)
	require.Equal(t, dividendAccount, divAcc)
}

func (suite *QuerierTestSuite) TestHandleDividendAccountRoot() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.BytesToHeimdallAddress([]byte("some-address")),
		big.NewInt(0).String(),
	)
	err := app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	path := []string{types.QueryDividendAccountRoot}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccountRoot)
	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryAccountProof() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	var accountRoot [32]byte

	path := []string{types.QueryAccountProof}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountProof)
	stakingInfo := &stakinginfo.Stakinginfo{}

	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.BytesToHeimdallAddress([]byte("some-address")),
		big.NewInt(0).String(),
	)
	err := app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	dividendAccounts := app.TopupKeeper.GetAllDividendAccounts(ctx)

	accRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)
	copy(accountRoot[:], accRoot)

	// mock contracts
	suite.contractCaller.On("GetStakingInfoInstance", mock.Anything).Return(stakingInfo, nil)
	suite.contractCaller.On("CurrentAccountStateRoot", stakingInfo).Return(accountRoot, nil)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(dividendAccount),
	}
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func (suite *QuerierTestSuite) TestHandleQueryVerifyAccountProof() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	dividendAccount := hmTypes.NewDividendAccount(
		hmTypes.BytesToHeimdallAddress([]byte("some-address")),
		big.NewInt(0).String(),
	)
	err := app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

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
