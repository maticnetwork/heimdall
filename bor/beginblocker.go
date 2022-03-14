package bor

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/consensus/bor"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/bor/client/rest"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {

	// TODO - Enter height
	if ctx.BlockHeight() == 1 {

		j, ok := rest.SPAN_OVERRIDES[helper.GenesisDoc.ChainID]
		if !ok {
			// TODO - Log error
		}

		var spans []*bor.ResponseWithHeight
		if err := json.Unmarshal(j, &spans); err != nil {
			return
		}

		for _, span := range spans {
			var heimdallSpan hmTypes.Span
			if err := json.Unmarshal(span.Result, &heimdallSpan); err != nil {
				continue
			}

			if err := k.AddNewRawSpan(ctx, heimdallSpan); err != nil {
				k.Logger(ctx).Error("Error AddNewRawSpan", "error", err)
			}
			k.UpdateLastSpan(ctx, heimdallSpan.ID)
		}
	}
}
