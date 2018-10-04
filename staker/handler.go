package staker
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
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
	fmt.Printf("validator is %v",validator)
	keeper.SetValidatorSet(ctx)
	return sdk.Result{}
}