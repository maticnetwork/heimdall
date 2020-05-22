package topup_test

import (
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Tests

func (suite *KeeperTestSuite) TestTopupSequenceSet() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	topupSequence := strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))

	app.TopupKeeper.SetTopupSequence(ctx, topupSequence)

	actualResult := app.TopupKeeper.HasTopupSequence(ctx, topupSequence)
	require.Equal(t, true, actualResult)
}

// tests setter/getters for Dividend account
func (suite *KeeperTestSuite) TestDividendAccount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	dividendAccount := types.DividendAccount{
		User:      hmTypes.BytesToHeimdallAddress([]byte("some-address")),
		FeeAmount: big.NewInt(0).String(),
	}
	app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	ok := app.TopupKeeper.CheckIfDividendAccountExists(ctx, dividendAccount.User)
	require.Equal(t, ok, true)
}

func (suite *KeeperTestSuite) TestAddFeeToDividendAccount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	address := hmTypes.HexToHeimdallAddress("234452")
	amount, _ := big.NewInt(0).SetString("0", 10)
	app.TopupKeeper.AddFeeToDividendAccount(ctx, address, amount)
	dividentAccount, _ := app.TopupKeeper.GetDividendAccountByAddress(ctx, address)
	actualResult, ok := big.NewInt(0).SetString(dividentAccount.FeeAmount, 10)
	require.Equal(t, ok, true)
	require.Equal(t, amount, actualResult)
}

func (suite *KeeperTestSuite) TestDividendAccountTree() {
	t := suite.T()

	divAccounts := make([]hmTypes.DividendAccount, 5)
	for i := 0; i < len(divAccounts); i++ {
		divAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.HexToHeimdallAddress("1234"),
			big.NewInt(0).String(),
		)
	}

	accountRoot, err := checkpointTypes.GetAccountRootHash(divAccounts)
	require.NotNil(t, accountRoot)
	require.NoError(t, err)

	accountProof, _, err := checkpointTypes.GetAccountProof(divAccounts, hmTypes.HexToHeimdallAddress("1234"))
	require.NotNil(t, accountProof)
	require.NoError(t, err)

	leafHash, err := divAccounts[0].CalculateHash()
	require.NotNil(t, leafHash)
	require.NoError(t, err)
}
