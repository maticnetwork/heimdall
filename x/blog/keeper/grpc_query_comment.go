package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/maticnetwork/heimdall/x/blog/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllComment(c context.Context, req *types.QueryAllCommentRequest) (*types.QueryAllCommentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var comments []*types.MsgComment
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	commentStore := prefix.NewStore(store, types.KeyPrefix(types.CommentKey))

	pageRes, err := query.Paginate(commentStore, req.Pagination, func(key []byte, value []byte) error {
		var comment types.MsgComment
		if err := k.cdc.UnmarshalBinaryBare(value, &comment); err != nil {
			return err
		}

		comments = append(comments, &comment)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCommentResponse{Comment: comments, Pagination: pageRes}, nil
}
