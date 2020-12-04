package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ApplyAndReturnValidatorSetUpdates applys and returns validator set for genutil module
func (k Keeper) ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	return nil, nil
}
