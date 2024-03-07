package clerk_test

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	querier        sdk.Querier
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier testing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = clerk.NewQuerier(suite.app.ClerkKeeper, &suite.contractCaller)
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

func (suite *QuerierTestSuite) TestHandleQueryRecord() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryRecord}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecord)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	_, sdkErr := querier(ctx, path, req)
	require.Error(t, sdkErr, "failed to parse params")

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryRecordParams(2)),
	}
	_, sdkErr = querier(ctx, path, req)
	require.Error(t, sdkErr, "could not get state record")

	hAddr := hmTypes.BytesToHeimdallAddress([]byte("some-address"))
	hHash := hmTypes.BytesToHeimdallHash([]byte("some-address"))
	testRecord1 := types.NewEventRecord(hHash, 1, 1, hAddr, make([]byte, 0), "1", time.Now())

	// SetEventRecord
	ck := app.ClerkKeeper
	err := ck.SetEventRecord(ctx, testRecord1)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryRecordParams(1)),
	}
	record, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, record)
}

func (suite *QuerierTestSuite) TestHandleQueryRecordList() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryRecordList}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordList)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	_, sdkErr := querier(ctx, path, req)
	require.Error(t, sdkErr, "failed to parse params")

	hAddr := hmTypes.BytesToHeimdallAddress([]byte("some-address"))
	hHash := hmTypes.BytesToHeimdallHash([]byte("some-address"))
	testRecord1 := types.NewEventRecord(hHash, 1, 1, hAddr, make([]byte, 0), "1", time.Now())

	// SetEventRecord
	ck := app.ClerkKeeper
	err := ck.SetEventRecord(ctx, testRecord1)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(hmTypes.NewQueryPaginationParams(1, 1)),
	}
	record, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, record)
}

func (suite *QuerierTestSuite) TestHandleQueryRecordSequence() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryRecordSequence}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordSequence)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	_, err := querier(ctx, path, req)
	require.Error(t, err, "failed to parse params")

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	chainParams := app.ChainKeeper.GetParams(ctx)
	txreceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(10),
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(nil, errors.New("err confirmed txn receipt"))

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryRecordSequenceParams("123", logIndex)),
	}
	_, err = querier(ctx, path, req)
	require.NotNil(t, err, "failed to parse params")

	index = simulation.RandIntBetween(r1, 0, 100)
	logIndex = uint64(index)
	txHash = hmTypes.HexToHeimdallHash("1234")

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryRecordSequenceParams("1234", logIndex)),
	}
	resp, err := querier(ctx, path, req)
	require.Nil(t, err)
	require.Nil(t, resp)

	testSeq := "1000010"
	ck := app.ClerkKeeper

	ck.SetRecordSequence(ctx, testSeq)

	logIndex = uint64(10)
	txHash = hmTypes.HexToHeimdallHash("12345")

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txreceipt, nil)

	req = abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryRecordSequenceParams("12345", logIndex)),
	}
	resp, err = querier(ctx, path, req)
	require.Nil(t, err)
	require.NotNil(t, resp)
}
