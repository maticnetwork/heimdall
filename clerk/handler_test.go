package clerk_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/clerk/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/tendermint/tendermint/crypto/secp256k1"

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
	suite.app, suite.ctx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = clerk.NewHandler(suite.app.ClerkKeeper, &suite.contractCaller)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgEventRecord() {
	t, _, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	id := uint64(1)
	blockNumber := uint64(1)
	chainId := "15001"
	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()
	msgEventRecord := types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		txHash,
		logIndex,
		blockNumber,
		id,
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		make([]byte, 0),
		chainId,
	)
	result := suite.handler(ctx, msgEventRecord)
	require.True(t, result.IsOK(), "expected msg record to be ok, got %v", result)

	chainId = "1"
	id = uint64(2)
	index = simulation.RandIntBetween(r1, 0, 100)
	logIndex = uint64(index)
	msgEventRecord = types.NewMsgEventRecord(
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		txHash,
		logIndex,
		blockNumber,
		id,
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		make([]byte, 0),
		chainId,
	)
	result = suite.handler(ctx, msgEventRecord)
	require.False(t, result.IsOK(), "error invalid bor chain id %v", result)
}
