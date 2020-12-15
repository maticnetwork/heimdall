package keeper

import (
	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/x/sidechannel/keeper"
	"github.com/maticnetwork/heimdall/x/sidechannel/types"
)

// NewQuerier creates new querier
func NewQuerier(_ keeper.Keeper, _ *codec.LegacyAmino) sdk.Querier {
	return func(_ sdk.Context, path []string, _ abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		// this line is used by starport scaffolding # 1
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}
