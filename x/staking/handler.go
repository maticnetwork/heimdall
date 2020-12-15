package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/maticnetwork/heimdall/helper"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	"github.com/maticnetwork/heimdall/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	types1 "github.com/tendermint/tendermint/proto/tendermint/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k, contractCaller)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgValidatorJoin:
			res, err := msgServer.ValidatorJoin(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgStakeUpdate:
			res, err := msgServer.StakeUpdate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSignerUpdate:
			res, err := msgServer.SignerUpdate(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgValidatorExit:
			res, err := msgServer.ValidatorExit(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		// this line is used by starport scaffolding # 1
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// NewSideTxHandler returns a side handler for "staking" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgValidatorJoin:
			return abci.ResponseDeliverSideTx{}
			// return SideHandleMsgValidatorJoin(ctx, msg, k, contractCaller)
		case *types.MsgValidatorExit:
			return abci.ResponseDeliverSideTx{}
			// return SideHandleMsgValidatorExit(ctx, msg, k, contractCaller)
		case *types.MsgSignerUpdate:
			return abci.ResponseDeliverSideTx{}
			// return SideHandleMsgSignerUpdate(ctx, msg, k, contractCaller)
		case *types.MsgStakeUpdate:
			return abci.ResponseDeliverSideTx{}
			// return SideHandleMsgStakeUpdate(ctx, msg, k, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(6), // TODO this should be defined like sdk.ErrUnknownRequest
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult types1.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgValidatorJoin:
			return nil
			// return PostHandleMsgValidatorJoin(ctx, k, msg, sideTxResult)
		case *types.MsgValidatorExit:
			return nil
			// return PostHandleMsgValidatorExit(ctx, k, msg, sideTxResult)
		case *types.MsgSignerUpdate:
			return nil
			// return PostHandleMsgSignerUpdate(ctx, k, msg, sideTxResult)
		case *types.MsgStakeUpdate:
			return nil
			// return PostHandleMsgStakeUpdate(ctx, k, msg, sideTxResult)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
