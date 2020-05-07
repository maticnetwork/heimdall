package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
)

// PostDeliverTxHandler runs after deliver tx handler
func (app *HeimdallApp) PostDeliverTxHandler(ctx sdk.Context, tx sdk.Tx, result sdk.Result) {
	height := ctx.BlockHeader().Height
	if height <= 2 {
		return
	}

	if result.IsOK() {
		anySideMsg := false
		for _, msg := range tx.GetMsgs() {
			if _, ok := msg.(types.SideTxMsg); ok {
				anySideMsg = true
				break
			}
		}

		// save tx bytes if any tx-msg is side-tx msg
		if anySideMsg && ctx.TxBytes() != nil {
			app.SidechannelKeeper.SetTx(ctx, ctx.BlockHeader().Height, ctx.TxBytes())
		}
	}
}

// BeginSideBlocker runs before side block
func (app *HeimdallApp) BeginSideBlocker(ctx sdk.Context, req abci.RequestBeginSideBlock) (res abci.ResponseBeginSideBlock) {
	height := ctx.BlockHeader().Height
	if height <= 2 {
		return
	}

	targetHeight := height - 2 // sidechannel takes 2 blocks to process

	// get logger
	logger := app.Logger()

	logger.Debug("[sidechannel] Processing side block", "height", height, "targetHeight", targetHeight)

	// get all validators
	// begin-block stores validators and end-block removes validators for each height - check sidechannel module.go
	validators := app.SidechannelKeeper.GetValidators(ctx, height)
	if len(validators) == 0 {
		return
	}

	// calculate power
	var totalPower int64
	for _, v := range validators {
		totalPower = totalPower + v.Power
	}

	// get empty events
	events := sdk.EmptyEvents()

	for _, sideTxResult := range req.SideTxResults {
		txHash := sideTxResult.TxHash
		// get tx from the store
		tx := app.SidechannelKeeper.GetTx(ctx, targetHeight, txHash)
		if tx != nil {
			// remove tx to avoid duplicate execution
			app.SidechannelKeeper.RemoveTx(ctx, targetHeight, txHash)

			usedValidator := make(map[int]bool)

			// signed power
			signedPower := make(map[abci.SideTxResultType]int64)
			signedPower[abci.SideTxResultType_Yes] = 0
			signedPower[abci.SideTxResultType_Skip] = 0
			signedPower[abci.SideTxResultType_No] = 0

			for _, sigObj := range sideTxResult.Sigs {
				// get validator by sig address
				if i := getValidatorIndexByAddress(sigObj.Address, validators); i != -1 {
					// check if validator already voted on tx
					if _, ok := usedValidator[i]; !ok {
						signedPower[sigObj.Result] = signedPower[sigObj.Result] + validators[i].Power
						usedValidator[i] = true
					}
				}
			}

			var result sdk.Result

			// check vote majority
			if signedPower[abci.SideTxResultType_Yes] >= (totalPower*2/3 + 1) {
				// approved
				logger.Debug("[sidechannel] Approved side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `yes`
				result = app.runTx(ctx, tx, abci.SideTxResultType_Yes)
			} else if signedPower[abci.SideTxResultType_No] >= (totalPower*2/3 + 1) {
				// rejected
				logger.Debug("[sidechannel] Rejected side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `no`
				result = app.runTx(ctx, tx, abci.SideTxResultType_No)
			} else {
				// skipped
				logger.Debug("[sidechannel] Skipped side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `skip`
				result = app.runTx(ctx, tx, abci.SideTxResultType_Skip)
			}

			// add events
			events = events.AppendEvents(result.Events)
		}
	}

	// remove all pending txs before exiting
	txs := app.SidechannelKeeper.GetTxs(ctx, targetHeight)
	for _, tx := range txs {
		app.SidechannelKeeper.RemoveTx(ctx, targetHeight, tx.Hash())

		// skipped
		logger.Debug("[sidechannel] Skipped side-tx", "txHash", hex.EncodeToString(tx.Hash()))

		// execute tx with `skip`
		result := app.runTx(ctx, tx, abci.SideTxResultType_Skip)

		// add events
		events = events.AppendEvents(result.Events)
	}

	// set event to response
	res.Events = events.ToABCIEvents()

	return res
}

// DeliverSideTxHandler runs for each side tx
func (app *HeimdallApp) DeliverSideTxHandler(ctx sdk.Context, tx sdk.Tx, req abci.RequestDeliverSideTx) (res abci.ResponseDeliverSideTx) {
	var code uint32
	var codespace string

	result := abci.SideTxResultType_Skip
	data := make([]byte, 0)

	for _, msg := range tx.GetMsgs() {
		sideMsg, isSideTxMsg := msg.(types.SideTxMsg)

		// match message route
		msgRoute := msg.Route()
		handlers := app.sideRouter.GetRoute(msgRoute)
		if handlers != nil && handlers.SideTxHandler != nil && isSideTxMsg {
			// Create a new context based off of the existing context with a cache wrapped multi-store (for state-less execution)
			runMsgCtx, _ := app.cacheTxContext(ctx, req.Tx)
			// execute side-tx handler
			msgResult := handlers.SideTxHandler(runMsgCtx, msg)

			// stop execution and return on first failed message
			if msgResult.Code != uint32(sdk.CodeOK) {
				// set data to empty if error
				data = make([]byte, 0)

				code = msgResult.Code
				codespace = msgResult.Codespace
				// skip side-tx if result is error
				result = abci.SideTxResultType_Skip
				break
			}

			// Each message result's Data must be length prefixed in order to separate
			// each result.
			data = append(data, msgResult.Data...)
			result = msgResult.Result

			// msg result is empty, get side sign bytes and append into data
			if len(msgResult.Data) == 0 {
				data = append(data, sideMsg.GetSideSignBytes()...)
			}
		}
	}

	return abci.ResponseDeliverSideTx{
		Code:      uint32(code),
		Codespace: string(codespace),
		Data:      data,
		Result:    result,
	}
}

//
// Internal functions
//

func (app *HeimdallApp) runTx(ctx sdk.Context, txBytes []byte, sideTxResult abci.SideTxResultType) (result sdk.Result) {
	// get decoder
	decoder := authTypes.DefaultTxDecoder(app.cdc)
	tx, err := decoder(txBytes)
	if err != nil {
		return
	}

	// recover if runMsgs fails
	defer func() {
		if r := recover(); r != nil {
			log := fmt.Sprintf("recovered: %v\nstack:\n%v", r, string(debug.Stack()))
			result = sdk.ErrInternal(log).Result()
		}
	}()

	// get context with tx bytes
	ctx = ctx.WithTxBytes(txBytes)
	// Create a new context based off of the existing context with a cache wrapped
	// multi-store in case message processing fails.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)
	result = app.runMsgs(runMsgCtx, tx.GetMsgs(), sideTxResult)
	// only update state if all messages pass
	if result.IsOK() {
		msCache.Write()
	}

	return
}

/// runMsgs iterates through all the messages and executes them.
func (app *HeimdallApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
	data := make([]byte, 0, len(msgs))

	var (
		code      sdk.CodeType
		codespace sdk.CodespaceType
	)

	// get empty events
	events := sdk.EmptyEvents()

	for _, msg := range msgs {
		_, isSideTxMsg := msg.(types.SideTxMsg)

		// match message route
		msgRoute := msg.Route()
		handler := app.sideRouter.GetRoute(msgRoute)
		if handler != nil && handler.PostTxHandler != nil && isSideTxMsg {
			msgResult := handler.PostTxHandler(ctx, msg, sideTxResult)

			// Each message result's Data must be length prefixed in order to separate
			// each result.
			data = append(data, msgResult.Data...)

			// msg events
			events = events.AppendEvents(msgResult.Events)

			// stop execution and return on first failed message
			if !msgResult.IsOK() {
				code = msgResult.Code
				codespace = msgResult.Codespace
				break
			}
		}
	}

	return sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Events:    events,
	}
}

// cacheTxContext returns a new context based off of the provided context with
// a cache wrapped multi-store.
func (app *HeimdallApp) cacheTxContext(ctx sdk.Context, txBytes []byte) (sdk.Context, sdk.CacheMultiStore) {
	ms := ctx.MultiStore()
	msCache := ms.CacheMultiStore()

	return ctx.WithMultiStore(msCache), msCache
}

//
// utils
//

func getValidatorIndexByAddress(address []byte, validators []abci.Validator) int {
	for i, v := range validators {
		if bytes.Equal(address, v.Address) {
			return i
		}
	}

	return -1
}
