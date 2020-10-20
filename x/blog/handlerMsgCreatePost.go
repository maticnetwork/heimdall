package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/blog/keeper"
	"github.com/maticnetwork/heimdall/x/blog/types"
)

func handleMsgCreatePost(ctx sdk.Context, k keeper.Keeper, post *types.MsgPost) (*sdk.Result, error) {
	k.CreatePost(ctx, *post)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
