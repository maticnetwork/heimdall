package clerk

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

// NewHandler creates new handler for handling messages for checkpoint module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case clerkTypes.MsgEventRecord:
			return handleMsgEventRecord(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in clerk module").Result()
		}
	}
}

func handleMsgEventRecord(ctx sdk.Context, msg clerkTypes.MsgEventRecord, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	// check if event record exists
	if exists := k.HasEventRecord(ctx, msg.ID); exists {
		return clerkTypes.ErrEventRecordAlreadySynced(k.Codespace()).Result()
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
		return clerkTypes.ErrEventRecordInvalid(k.Codespace()).Result()
	}

	// create event record
	record := clerkTypes.NewEventRecord(
		msg.TxHash,
		msg.LogIndex,
		msg.ID,
		types.BytesToHeimdallAddress(parsedLog.ContractAddress.Bytes()),
		parsedLog.Data,
	)

	// save event into state
	if err := k.SetEventRecord(ctx, record); err != nil {
		k.Logger(ctx).Error("Unable to update event record", "error", err, "id", msg.ID)
		return clerkTypes.ErrEventUpdate(k.Codespace()).Result()
	}

	// increment state sync event count
	k.IncrementStateSyncEventCount(ctx)

	resTags := sdk.NewTags(
		clerkTypes.RecordID, []byte(strconv.FormatUint(msg.ID, 10)),
		clerkTypes.RecordContract, []byte(parsedLog.ContractAddress.String()),
		clerkTypes.RecordTxHash, []byte(msg.TxHash.String()),
		clerkTypes.RecordTxLogIndex, []byte(strconv.FormatUint(msg.LogIndex, 10)),
	)

	return sdk.Result{Tags: resTags}
}
