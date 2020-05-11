package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// SideTxHandler defines the core of side-tx execution of an application
type SideTxHandler func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx

// PostTxHandler defines the core of the state transition function of an application after side-tx execution
type PostTxHandler func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result
