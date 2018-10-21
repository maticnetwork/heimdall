package staker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	contract "github.com/maticnetwork/heimdall/contracts"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateMaticValidator:
			return handleMsgCreateMaticValidator(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in staker module ").Result()
		}
	}
}
func handleMsgCreateMaticValidator(ctx sdk.Context, validator MsgCreateMaticValidator, keeper Keeper) sdk.Result {
	fmt.Println("entered handler")
	fmt.Printf("validator is %v", validator)
	//keeper.SetValidatorSet(ctx)
	return sdk.Result{}
}

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {

	// validator := k.GetValidatorInfo(ctx, _address)
	fmt.Printf("prev validators are %v \n", k.GetAllValidators(ctx))

	validatorSet := contract.GetValidators()
	fmt.Printf("all validators are %v \n", k.GetAllValidators(ctx))
	k.FlushValidatorSet(ctx)
	k.SetValidatorSet(ctx, validatorSet)

	return validatorSet
}
