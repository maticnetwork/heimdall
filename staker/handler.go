package staker

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"encoding/hex"

	"github.com/tendermint/tendermint/privval"

	tmtypes "github.com/tendermint/tendermint/types"
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
	privVal := privval.LoadFilePV("/Users/vc/.tendermint/config/priv_validator.json")
	fmt.Printf("the _pub is %v and address %v", privVal.GetPubKey(), hex.EncodeToString(privVal.Address))
	_address, _ := hex.DecodeString(hex.EncodeToString(privVal.Address))
	// validator := k.GetValidatorInfo(ctx, _address)
	val1 := abci.Validator{
		Address: _address,
		Power:   int64(1),
		PubKey:  tmtypes.TM2PB.PubKey(privVal.GetPubKey()),
	}
	validatorSet := []abci.Validator{val1}

	fmt.Printf("all validators are %v \n", k.GetAllValidators(ctx))
	k.SetValidatorSet(ctx, validatorSet)
	// k.FlushValidatorSet(ctx)
	return validatorSet
}
