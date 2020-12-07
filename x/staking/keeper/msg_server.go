package keeper

import (
	"context"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

type msgServer struct {
	Keeper
	contractCaller helper.IContractCaller
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.MsgServer {
	return &msgServer{Keeper: keeper, contractCaller: contractCaller}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ValidatorJoin(goCtx context.Context, msg *types.MsgValidatorJoin) (*types.MsgValidatorJoinResponse, error) {
	// k.Logger(ctx).Info("Handling new validator join", "msg", msg)

	// ctx := sdk.UnwrapSDKContext(goCtx)
	// params := k.GetParams(ctx)
	// chainParams := params.

	return &types.MsgValidatorJoinResponse{}, nil

}

func (k msgServer) StakeUpdate(goCtx context.Context, msg *types.MsgStakeUpdate) (*types.MsgStakeUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgStakeUpdateResponse{}, nil
}

func (k msgServer) SignerUpdate(goCtx context.Context, msg *types.MsgSignerUpdate) (*types.MsgSignerUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgSignerUpdateResponse{}, nil

}

func (k msgServer) ValidatorExit(goCtx context.Context, msg *types.MsgValidatorExit) (*types.MsgValidatorExitResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgValidatorExitResponse{}, nil
}
