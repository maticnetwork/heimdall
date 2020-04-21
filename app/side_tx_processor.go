package app

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// PostDeliverTxHandler runs after deliver tx handler
func (app *HeimdallApp) PostDeliverTxHandler(ctx sdk.Context, tx sdk.Tx, result sdk.Result) {
	if result.IsOK() {
		app.SidechannelKeeper.SetTx(ctx, ctx.BlockHeader().Height, ctx.TxBytes())
	}
}

// BeginSideBlocker runs before side block
func (app *HeimdallApp) BeginSideBlocker(ctx sdk.Context, req abci.RequestBeginSideBlock) (res abci.ResponseBeginSideBlock) {
	height := ctx.BlockHeader().Height
	if height <= 2 {
		return
	}

	targetHeight := height - 2 // sidechannel takes 2 blocks to process

	fmt.Println("[sidechannel] Processing side block",
		"height", height,
		"targetHeight", targetHeight,
	)

	// get all txs
	txs := app.SidechannelKeeper.GetTxs(ctx, targetHeight)
	if len(txs) == 0 {
		return
	}

	// remove all txs before exiting
	defer func() {
		for _, tx := range txs {
			app.SidechannelKeeper.RemoveTx(ctx, targetHeight, tx.Hash())
		}
	}()

	// get all validators
	validators := app.SidechannelKeeper.GetValidators(ctx, targetHeight)
	if len(validators) == 0 {
		return
	}

	// calculate power
	var totalPower int64
	for _, v := range validators {
		totalPower = totalPower + v.Power
	}

	for _, r := range req.SideTxResults {
		fmt.Println("                              ", "txHash", hex.EncodeToString(r.TxHash), "totalPower", totalPower)
		for _, s := range r.Sigs {
			fmt.Println("                                 ", "result", s.Result, "sig", hex.EncodeToString(s.Sig))
		}
	}

	return res
}

// DeliverSideTxHandler runs for each side tx
func (app *HeimdallApp) DeliverSideTxHandler(ctx sdk.Context, tx sdk.Tx, req abci.RequestDeliverSideTx) (res abci.ResponseDeliverSideTx) {
	fmt.Println("[sidechannel] DeliverSideTx",
		"tx", hex.EncodeToString(req.Tx),
		"msgRoute", tx.GetMsgs()[0].Route(),
		"msgType", tx.GetMsgs()[0].Type(),
	)

	var result sdk.Result

	return abci.ResponseDeliverSideTx{
		Code:      uint32(result.Code),
		Codespace: string(result.Codespace),
		Data:      result.Data,

		// yes vote
		Result: abci.SideTxResultType_Yes,
	}
}
