package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/blog/types"
)

func (k Keeper) CreateComment(ctx sdk.Context, comment types.MsgComment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	b := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(types.KeyPrefix(types.CommentKey), b)
}

func (k Keeper) GetAllComment(ctx sdk.Context) (msgs []types.MsgComment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.CommentKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgComment
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
