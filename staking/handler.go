package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {
	// validator := k.GetValidatorInfo(ctx, _address)
	fmt.Printf("prev validators are %v \n", k.GetAllValidators(ctx))

	validatorSet := helper.GetValidators()
	fmt.Printf("all validators are %v \n", k.GetAllValidators(ctx))
	k.FlushValidatorSet(ctx)
	k.SetValidatorSet(ctx, validatorSet)

	return validatorSet
}
