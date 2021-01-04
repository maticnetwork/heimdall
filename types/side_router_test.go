package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/types"
)

var testSideTxHandler = func(_ sdk.Context, _ sdk.Msg) abci.ResponseDeliverSideTx {
	return abci.ResponseDeliverSideTx{}
}

var testPostTxHandler = func(_ sdk.Context, _ sdk.Msg, _ tmprototypes.SideTxResultType) (*sdk.Result, error) {
	return &sdk.Result{}, nil
}

func TestSideRouter(t *testing.T) {
	rtr := types.NewSideRouter()
	handler := &types.SideHandlers{
		SideTxHandler: testSideTxHandler,
		PostTxHandler: testPostTxHandler,
	}

	rtr.AddRoute("testRoute", handler)
	h := rtr.GetRoute("testRoute")
	require.NotNil(t, h)

	// require panic on duplicate route
	require.Panics(t, func() {
		rtr.AddRoute("testRoute", handler)
	})
}
