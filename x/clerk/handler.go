package clerk

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/clerk/keeper"
	"github.com/maticnetwork/heimdall/x/clerk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/maticnetwork/heimdall/helper"
)

// NewHandler returns a handler for "clerk" type messages.
func NewHandler(k keeper.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k, contractCaller)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgEventRecordRequest:
			res, err := msgServer.MsgEventRecord(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
