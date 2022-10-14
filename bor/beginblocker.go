package bor

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/consensus/bor"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/bor/client/rest"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {

	if ctx.BlockHeight() == int64(helper.SpanOverrideBlockHeight) {
		k.Logger(ctx).Info("overriding span BeginBlocker", "height", ctx.BlockHeight())
		j, ok := rest.SPAN_OVERRIDES[helper.GenesisDoc.ChainID]
		if !ok {
			k.Logger(ctx).Info("No Override span found")
			return
		}

		var spans []*bor.ResponseWithHeight
		if err := json.Unmarshal(j, &spans); err != nil {
			k.Logger(ctx).Error("Error Unmarshal spans", "error", err)
			panic(err)
		}

		for _, span := range spans {
			k.Logger(ctx).Info("overriding span", "height", span.Height, "span", span)
			var heimdallSpan hmTypes.Span
			if err := json.Unmarshal(span.Result, &heimdallSpan); err != nil {
				k.Logger(ctx).Error("Error Unmarshal heimdallSpan", "error", err)
				panic(err)
			}

			if err := k.AddNewRawSpan(ctx, heimdallSpan); err != nil {
				k.Logger(ctx).Error("Error AddNewRawSpan", "error", err)
				panic(err)
			}
			k.UpdateLastSpan(ctx, heimdallSpan.ID)
		}
	}
}
