package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/types"
)

var testSideTxHandler = func(_ sdk.Context, _ sdk.Msg) abci.ResponseDeliverSideTx {
	return abci.ResponseDeliverSideTx{}
}

var testPostTxHandler = func(_ sdk.Context, _ sdk.Msg, _ abci.SideTxResultType) sdk.Result {
	return sdk.Result{}
}

func TestSideRouter(t *testing.T) {
	rtr := types.NewSideRouter()
	handler := &types.SideHandlers{
		SideTxHandler: testSideTxHandler,
		PostTxHandler: testPostTxHandler,
	}

	// require panic on invalid route
	require.Panics(t, func() {
		rtr.AddRoute("*", handler)
	})

	rtr.AddRoute("testRoute", handler)
	h := rtr.GetRoute("testRoute")
	require.NotNil(t, h)

	// require panic on duplicate route
	require.Panics(t, func() {
		rtr.AddRoute("testRoute", handler)
	})
}
