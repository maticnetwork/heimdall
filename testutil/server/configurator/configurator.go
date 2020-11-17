// Package configurator a Fixture defines an implementation of server.Fixture which uses an in-memory test data store
// with no backing ABCI app. It implements the module.Configurator interface and is designed
// to be used in server integration tests for each module independent of a larger app.
//
// Ex:
//
package configurator

import (
	"context"
	"testing"

	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/matiknetwork/heimdall/testutil/server"
)

// Fixture is an implementation of server.Fixture which uses an in-memory test data store
// with no backing ABCI app. It implements the module.Configurator interface and is designed
// to be used in server integration tests for each module independent of a larger app.
type Fixture struct {
	queryRouter *testRouter
	msgRouter   *testRouter
	keys        []sdk.StoreKey
	t           *testing.T
	signers     []sdk.AccAddress
	ctx         context.Context
}

// NewFixture returns a new Fixture instance.
func NewFixture(t *testing.T, storeKeys []sdk.StoreKey, signers []sdk.AccAddress) Fixture {
	return Fixture{
		queryRouter: newTestRouter(false),
		msgRouter:   newTestRouter(true),
		keys:        storeKeys,
		t:           t,
		signers:     signers,
	}
}

var _ module.Configurator = Fixture{}

// MsgServer implements the Configurator.MsgServer method
func (c Fixture) MsgServer() gogogrpc.Server {
	return c.msgRouter
}

// QueryServer implements the Configurator.QueryServer method
func (c Fixture) QueryServer() gogogrpc.Server {
	return c.queryRouter
}

var _ server.FixtureFactory = Fixture{}

// Setup implements the FixtureFactory.Setup method
func (c Fixture) Setup() server.Fixture {
	db := dbm.NewMemDB()

	ms := store.NewCommitMultiStore(db)
	for _, key := range c.keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	}
	err := ms.LoadLatestVersion()
	require.NoError(c.t, err)

	c.ctx = sdk.WrapSDKContext(sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger()))

	return c
}

var _ server.Fixture = Fixture{}

// TxConn implements the Fixture.TxConn method
func (c Fixture) TxConn() grpc.ClientConnInterface {
	return c.msgRouter
}

// QueryConn implements the Fixture.QueryConn method
func (c Fixture) QueryConn() grpc.ClientConnInterface {
	return c.queryRouter
}

// Signers implements the Fixture.Signers method
func (c Fixture) Signers() []sdk.AccAddress {
	return c.signers
}

// Context implements the Fixture.Context method
func (c Fixture) Context() context.Context {
	return c.ctx
}

// Teardown implements the Fixture.Teardown method
func (c Fixture) Teardown() {}
