package clerk

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler creates new handler for handling messages for checkpoint module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgEventRecord:
			return handleMsgEventRecord(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in clerk module").Result()
		}
	}
}

func handleMsgEventRecord(ctx sdk.Context, msg types.MsgEventRecord, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	// check if event record exists
	if exists := k.HasEventRecord(ctx, msg.ID); exists {
		return types.ErrEventRecordAlreadySynced(k.Codespace()).Result()
	}

	// get confirmed tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash())
	if receipt == nil || err != nil {
		return common.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// get event log for topup
	eventLog, err := contractCaller.DecodeStateSyncedEvent(helper.GetStateSenderAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return common.ErrInvalidMsg(k.Codespace(), "Unable to fetch log for txHash").Result()
	}

	// check if message and event log matches
	if eventLog.Id.Uint64() != msg.ID {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "stateIdFromTx", eventLog.Id)
		return common.ErrInvalidMsg(k.Codespace(), "ID in message doesn't match with id in log. msgId %v stateIdFromTx %v", msg.ID, eventLog.Id).Result()
	}

	// create event record
	record := types.NewEventRecord(
		msg.TxHash,
		msg.LogIndex,
		eventLog.Id.Uint64(),
		hmTypes.BytesToHeimdallAddress(eventLog.ContractAddress.Bytes()),
		eventLog.Data,
		msg.ChainID,
	)

	// save event into state
	if err := k.SetEventRecord(ctx, record); err != nil {
		k.Logger(ctx).Error("Unable to update event record", "error", err, "id", msg.ID)
		return types.ErrEventUpdate(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecord,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyRecordID, strconv.FormatUint(msg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeyRecordContract, eventLog.ContractAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecordTxHash, msg.TxHash.String()),
			sdk.NewAttribute(types.AttributeKeyRecordTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
