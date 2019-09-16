package clerk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
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
	resTags := sdk.NewTags(
	// clerkTypes.RecordID, []byte(strconv.FormatUint(msg.ID, 10)),
	)

	return sdk.Result{Tags: resTags}
}
