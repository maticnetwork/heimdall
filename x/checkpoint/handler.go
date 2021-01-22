package checkpoint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/checkpoint/keeper"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k, contractCaller)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCheckpoint:
			res, err := msgServer.Checkpoint(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCheckpointAck:
			res, err := msgServer.CheckpointAck(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgCheckpointNoAck:
			res, err := msgServer.CheckpointNoAck(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
