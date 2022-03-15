package clerk

import (
	"bytes"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "topup" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgEventRecord:
			return SideHandleMsgEventRecord(ctx, k, msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgEventRecord:
			return PostHandleMsgEventRecord(ctx, k, msg, sideTxResult)
		default:
			return sdk.ErrUnknownRequest("Unknown msg type").Result()
		}
	}
}

func SideHandleMsgEventRecord(ctx sdk.Context, k Keeper, msg types.MsgEventRecord, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ Validating External call for clerk msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get confirmed tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if receipt == nil || err != nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	// get event log for topup
	eventLog, err := contractCaller.DecodeStateSyncedEvent(chainParams.StateSenderAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64())
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check if message and event log matches
	if eventLog.Id.Uint64() != msg.ID {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "stateIdFromTx", eventLog.Id)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if !bytes.Equal(eventLog.ContractAddress.Bytes(), msg.ContractAddress.Bytes()) {
		k.Logger(ctx).Error(
			"ContractAddress from event does not match with Msg ContractAddress",
			"EventContractAddress", eventLog.ContractAddress.String(),
			"MsgContractAddress", msg.ContractAddress.String(),
		)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if !bytes.Equal(eventLog.Data, msg.Data) {
		if ctx.BlockHeight() > helper.SpanOverrideBlockHeight {
			if !(len(eventLog.Data) > helper.MaxStateSyncSize && bytes.Equal(msg.Data, hmTypes.HexToHexBytes(""))) {
				k.Logger(ctx).Error(
					"Data from event does not match with Msg Data",
					"EventData", hmTypes.BytesToHexBytes(eventLog.Data),
					"MsgData", hmTypes.BytesToHexBytes(msg.Data),
				)
				return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
			}
		} else {
			if !(len(eventLog.Data) > helper.LegacyMaxStateSyncSize && bytes.Equal(msg.Data, hmTypes.HexToHexBytes(""))) {
				k.Logger(ctx).Error(
					"Data from event does not match with Msg Data",
					"EventData", hmTypes.BytesToHexBytes(eventLog.Data),
					"MsgData", hmTypes.BytesToHexBytes(msg.Data),
				)
				return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
			}
		}
	}

	result.Result = abci.SideTxResultType_Yes
	return
}

func PostHandleMsgEventRecord(ctx sdk.Context, k Keeper, msg types.MsgEventRecord, sideTxResult abci.SideTxResultType) sdk.Result {

	// Skip handler if clerk is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new clerk since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay
	if k.HasEventRecord(ctx, msg.ID) {
		k.Logger(ctx).Debug("Skipping new clerk record as it's already processed")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting clerk state", "sideTxResult", sideTxResult)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// create event record
	record := types.NewEventRecord(
		msg.TxHash,
		msg.LogIndex,
		msg.ID,
		msg.ContractAddress,
		msg.Data,
		msg.ChainID,
		ctx.BlockTime(),
	)

	// save event into state
	if err := k.SetEventRecord(ctx, record); err != nil {
		k.Logger(ctx).Error("Unable to update event record", "error", err, "id", msg.ID)
		return types.ErrEventUpdate(k.Codespace()).Result()
	}

	// save record sequence
	k.SetRecordSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecord,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(types.AttributeKeyRecordTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()), // result
			sdk.NewAttribute(types.AttributeKeyRecordID, strconv.FormatUint(msg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeyRecordContract, msg.ContractAddress.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
