package app

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/types"
)

// PostDeliverTxHandler runs after deliver tx handler
func (app *HeimdallApp) PostDeliverTxHandler(ctx sdk.Context, tx sdk.Tx, result *sdk.Result) {
	height := ctx.BlockHeader().Height
	if height <= 2 {
		return
	}

	anySideMsg := false
	for _, msg := range tx.GetMsgs() {
		if _, ok := IsSideMsg(msg); ok {
			anySideMsg = true
			break
		}
	}

	// save tx bytes if any tx-msg is side-tx msg
	if anySideMsg && ctx.TxBytes() != nil {
		app.SidechannelKeeper.SetTx(ctx, uint64(height), ctx.TxBytes())
	}
}

// BeginSideBlocker runs before side block
func (app *HeimdallApp) BeginSideBlocker(ctx sdk.Context, req abci.RequestBeginSideBlock) (res abci.ResponseBeginSideBlock) {
	height := ctx.BlockHeader().Height
	if height <= 2 {
		return
	}

	targetHeight := uint64(height - 2) // sidechannel takes 2 blocks to process

	// get logger
	logger := app.Logger()

	logger.Debug("[sidechannel] Processing side block", "height", height, "targetHeight", targetHeight)

	// get all validators
	// begin-block stores validators and end-block removes validators for each height - check sidechannel module.go
	validators := app.SidechannelKeeper.GetValidators(ctx, uint64(height))
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
			signedPower := make(map[tmprototypes.SideTxResultType]int64)
			signedPower[tmprototypes.SideTxResultType_YES] = 0
			signedPower[tmprototypes.SideTxResultType_SKIP] = 0
			signedPower[tmprototypes.SideTxResultType_NO] = 0

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

			var result *sdk.Result
			var rerr error

			// check vote majority
			if signedPower[tmprototypes.SideTxResultType_YES] >= (totalPower*2/3 + 1) {
				// approved
				logger.Debug("[sidechannel] Approved side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `yes`
				result, rerr = app.runTx(ctx, tx, tmprototypes.SideTxResultType_YES)
			} else if signedPower[tmprototypes.SideTxResultType_NO] >= (totalPower*2/3 + 1) {
				// rejected
				logger.Debug("[sidechannel] Rejected side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `no`
				result, rerr = app.runTx(ctx, tx, tmprototypes.SideTxResultType_NO)
			} else {
				// skipped
				logger.Debug("[sidechannel] Skipped side-tx", "txHash", hex.EncodeToString(tx.Hash()))

				// execute tx with `skip`
				result, rerr = app.runTx(ctx, tx, tmprototypes.SideTxResultType_SKIP)
			}

			if rerr != nil {
				logger.Error("[sidechannel] Error while processing side-tx in begin side block",
					"txHash", hex.EncodeToString(tx.Hash()),
					"totalPower", totalPower,
					"yesVotes", signedPower[tmprototypes.SideTxResultType_YES],
					"noVotes", signedPower[tmprototypes.SideTxResultType_NO],
					"skipVotes", signedPower[tmprototypes.SideTxResultType_SKIP],
					"err", rerr,
				)
			} else {
				// add events
				events = events.AppendEvents(result.GetEvents())
			}
		}
	}

	// remove all pending txs before exiting
	txs := app.SidechannelKeeper.GetTxs(ctx, targetHeight)
	for _, tx := range txs {
		app.SidechannelKeeper.RemoveTx(ctx, targetHeight, tx.Hash())

		// skipped
		logger.Debug("[sidechannel] Skipped side-tx", "txHash", hex.EncodeToString(tx.Hash()))

		// execute tx with `skip`
		result, serr := app.runTx(ctx, tx, tmprototypes.SideTxResultType_SKIP)
		if serr != nil {
			logger.Error("[sidechannel] Error while processing skipped side-tx in beginside block", "txHash", hex.EncodeToString(tx.Hash()))
		} else {
			// add events
			events = events.AppendEvents(result.GetEvents())
		}
	}

	// set event to response
	res.Events = events.ToABCIEvents()

	return res
}

// DeliverSideTxHandler runs for each side tx
func (app *HeimdallApp) DeliverSideTxHandler(ctx sdk.Context, tx sdk.Tx, req abci.RequestDeliverSideTx) (res abci.ResponseDeliverSideTx) {
	var code uint32
	var codespace string

	result := tmprototypes.SideTxResultType_SKIP
	data := make([]byte, 0)

	for _, msg := range tx.GetMsgs() {
		sideMsg, isSideTxMsg := IsSideMsg(msg)

		// match message route
		msgRoute := msg.Route()
		handlers := app.sideRouter.GetRoute(msgRoute)
		if handlers != nil && handlers.SideTxHandler != nil && isSideTxMsg {
			// Create a new context based off of the existing context with a cache wrapped multi-store (for state-less execution)
			runMsgCtx, _ := app.cacheTxContext(ctx, req.Tx)
			// execute side-tx handler
			msgResult := handlers.SideTxHandler(runMsgCtx, msg)

			// stop execution and return on first failed message
			if msgResult.Code != abci.CodeTypeOK {
				// set data to empty if error
				data = make([]byte, 0)

				code = msgResult.Code
				codespace = msgResult.Codespace
				// skip side-tx if result is error
				result = tmprototypes.SideTxResultType_SKIP
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

func (app *HeimdallApp) runTx(ctx sdk.Context, txBytes []byte, sideTxResult tmprototypes.SideTxResultType) (result *sdk.Result, err error) {
	// decode tx
	tx, err := app.txDecoder(txBytes)
	if err != nil {
		return
	}

	// recover if runMsgs fails
	defer func() {
		if r := recover(); r != nil {
			log := fmt.Sprintf("recovered: %v\nstack:\n%v", r, string(debug.Stack()))

			// set error and result
			result = nil
			err = errors.New(log)
		}
	}()

	// get context with tx bytes
	ctx = ctx.WithTxBytes(txBytes)
	// Create a new context based off of the existing context with a cache wrapped
	// multi-store in case message processing fails.
	runMsgCtx, msCache := app.cacheTxContext(ctx, txBytes)
	result, err = app.runMsgs(runMsgCtx, tx.GetMsgs(), sideTxResult)
	// only update state if all messages pass
	if err == nil {
		msCache.Write()
	}

	return
}

/// runMsgs iterates through all the messages and executes them.
func (app *HeimdallApp) runMsgs(ctx sdk.Context, msgs []sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	// get empty events
	msgLogs := make(sdk.ABCIMessageLogs, 0, len(msgs))
	events := sdk.EmptyEvents()
	txMsgData := &sdk.TxMsgData{
		Data: make([]*sdk.MsgData, 0, len(msgs)),
	}

	for i, msg := range msgs {
		_, isSideTxMsg := IsSideMsg(msg)

		// match message route
		msgRoute := msg.Route()
		handler := app.sideRouter.GetRoute(msgRoute)
		if handler != nil && handler.PostTxHandler != nil && isSideTxMsg {
			msgResult, perr := handler.PostTxHandler(ctx, msg, sideTxResult)
			// stop execution and return on first failed message
			if perr != nil {
				return nil, sdkerrors.Wrapf(perr, "failed to execute message with post tx handler; message index: %d", i)
			}

			// append data
			txMsgData.Data = append(txMsgData.Data, &sdk.MsgData{
				MsgType: msg.Type(),
				Data:    msgResult.Data,
			})

			// msg events
			events = events.AppendEvents(msgResult.GetEvents())

			// add msg logs
			msgLogs = append(msgLogs, sdk.NewABCIMessageLog(uint32(i), msgResult.Log, events))
		}
	}

	data, err := proto.Marshal(txMsgData)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal tx msg data")
	}

	return &sdk.Result{
		Log:    msgLogs.String(),
		Data:   data,
		Events: events.ToABCIEvents(),
	}, nil
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

// IsSideMsg is side msg
func IsSideMsg(msg sdk.Msg) (types.SideTxMsg, bool) {
	if svcMsg, ok := msg.(sdk.ServiceMsg); ok {
		if m, ok := svcMsg.Request.(types.SideTxMsg); ok {
			return m, true
		}
	} else if m, ok := msg.(types.SideTxMsg); ok {
		return m, true
	}

	return nil, false
}

func getValidatorIndexByAddress(address []byte, validators []*abci.Validator) int {
	for i, v := range validators {
		if bytes.Equal(address, v.Address) {
			return i
		}
	}

	return -1
}
