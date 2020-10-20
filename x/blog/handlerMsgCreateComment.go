package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/blog/keeper"
	"github.com/maticnetwork/heimdall/x/blog/types"
)

func handleMsgCreateComment(ctx sdk.Context, k keeper.Keeper, comment *types.MsgComment) (*sdk.Result, error) {
	k.CreateComment(ctx, *comment)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
