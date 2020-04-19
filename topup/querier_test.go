package topup_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	suite.querier = topup.NewQuerier(suite.app.TopupKeeper, &suite.contractCaller)
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

// TestQuerySequence queries sequence data
func (suite *QuerierTestSuite) TestQuerySequence() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")

	logIndex := uint64(simulation.RandIntBetween(r1, 0, 100))

	chainParams := app.ChainKeeper.GetParams(ctx)
	txreceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}

	// set topup sequence
	sequence := new(big.Int).Mul(txreceipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))
	app.TopupKeeper.SetTopupSequence(ctx, sequence.String())

	//NEEDHELP: Here time.Now().UTC() does not maches with time.Now().UTC() in querySequence
	// for now, have changed querySequence time to time.Unix(0, 0)
	suite.contractCaller.On("GetConfirmedTxReceipt", time.Unix(0, 0), txHash.EthHash(), chainParams.TxConfirmationTime).Return(txreceipt, nil)

	path := []string{types.QuerySequence}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySequence)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQuerySequenceParams(txHash.String(), logIndex)),
	}
	res, err := querier(ctx, path, req)
	// check no error found
	require.NoError(t, err)

	// check response is not nil
	require.NotNil(t, res)
	require.Equal(t, sequence.String(), string(res))
}
