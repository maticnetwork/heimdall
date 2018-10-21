package staker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	contract "github.com/maticnetwork/heimdall/contracts"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {

	// validator := k.GetValidatorInfo(ctx, _address)
	fmt.Printf("prev validators are %v \n", k.GetAllValidators(ctx))

	validatorSet := contract.GetValidators()
	fmt.Printf("all validators are %v \n", k.GetAllValidators(ctx))
	k.FlushValidatorSet(ctx)
	k.SetValidatorSet(ctx, validatorSet)

	return validatorSet
}
