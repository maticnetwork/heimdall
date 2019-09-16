package clerk

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

// NewHandler creates new handler for handling messages for checkpoint module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case clerkTypes.MsgStateRecord:
			return handleMsgStateRecord(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in clerk module").Result()
		}
	}
}

func handleMsgStateRecord(ctx sdk.Context, msg clerkTypes.MsgStateRecord, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash.EthHash()); !confirmed {
		return common.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// check if event record exists
	if exists := k.HasEventRecord(ctx, msg.ID); exists {
		return types.ErrEventRecordAlreadySynced(k.Codespace()).Result()
	}

	// create event record
	record := clerkTypes.EventRecord{
		ID:       msg.ID,
		Contract: msg.Contract,
		Data:     msg.Data,
	}

	// save event into state
	if err := k.SetEventRecord(ctx, record); err != nil {
		k.Logger(ctx).Error("Unable to update event record", "error", err, "id", msg.ID)
		return types.ErrEventUpdate(k.Codespace()).Result()
	}

	resTags := sdk.NewTags(
		clerkTypes.RecordID, []byte(strconv.FormatUint(msg.ID, 10)),
		clerkTypes.RecordContract, []byte(msg.Contract.String()),
		clerkTypes.RecordTxHash, []byte(msg.TxHash.String()),
	)

	return sdk.Result{Tags: resTags}
}
