package clerk

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler creates new handler for handling messages for checkpoint module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
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
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if receipt == nil || err != nil {
		return common.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	var parsedLog *statesender.StatesenderStateSynced
	for _, log := range receipt.Logs {
		if uint64(log.Index) == msg.LogIndex && len(log.Topics) == 3 {
			p, err := contractCaller.EncodeStateSyncedEvent(log)
			if err != nil {
				break
			}

			if p != nil && msg.ID == p.Id.Uint64() {
				parsedLog = p
			}
		}
	}

	if parsedLog == nil {
		return types.ErrEventRecordInvalid(k.Codespace()).Result()
	}

	// create event record
	record := types.NewEventRecord(
		msg.TxHash,
		msg.LogIndex,
		msg.ID,
		hmTypes.BytesToHeimdallAddress(parsedLog.ContractAddress.Bytes()),
		parsedLog.Data,
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
			sdk.NewAttribute(types.AttributeKeyRecordContract, parsedLog.ContractAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecordTxHash, msg.TxHash.String()),
			sdk.NewAttribute(types.AttributeKeyRecordTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
