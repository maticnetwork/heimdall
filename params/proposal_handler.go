package params

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/maticnetwork/heimdall/gov/types"
	"github.com/maticnetwork/heimdall/params/types"
)

// NewParamChangeProposalHandler new param changes proposal handler
func NewParamChangeProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case types.ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k Keeper, p types.ParameterChangeProposal) sdk.Error {
	for _, c := range p.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return types.ErrUnknownSubspace(k.codespace, c.Subspace)
		}

		var err error
		if len(c.Subkey) == 0 {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, value: %s", c.Key, c.Value),
			)

			err = ss.Update(ctx, []byte(c.Key), []byte(c.Value))
		} else {
			k.Logger(ctx).Info(
				fmt.Sprintf("setting new parameter; key: %s, subkey: %s, value: %s", c.Key, c.Subspace, c.Value),
			)
			err = ss.UpdateWithSubkey(ctx, []byte(c.Key), []byte(c.Subkey), []byte(c.Value))
		}

		if err != nil {
			return types.ErrSettingParameter(k.codespace, c.Key, c.Subkey, c.Value, err.Error())
		}
	}

	return nil
}
