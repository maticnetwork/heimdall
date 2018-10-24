package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
	"log"
)

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {
	// validator := k.GetValidatorInfo(ctx, _address)
	log.Print("Current Validators Are : %v \n", k.GetAllValidators(ctx))

	validatorSet := helper.GetValidators()
	k.FlushValidatorSet(ctx)
	k.SetValidatorSet(ctx, validatorSet)
	log.Print("New Validators Are : %v \n", k.GetAllValidators(ctx))

	return validatorSet
}
