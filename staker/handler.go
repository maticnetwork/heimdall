package staker
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"

	"encoding/hex"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
	fmt.Printf("validator is %v",validator)
	//keeper.SetValidatorSet(ctx)
	return sdk.Result{}
}

func EndBlocker(ctx sdk.Context,k Keeper) (validators []abci.Validator){
	var _pubkey secp256k1.PubKeySecp256k1
	_pub, _ := hex.DecodeString("041FE1CDE7D9D8C9182AC967EC8362262216FF8A10061F0DE0F1472F9E45F965D0909DE527E18C7BFB9FCD42335E60FB6E18367A4DC37F1A7FC3265C7241597973")
	copy(_pubkey[:], _pub[:])
	_address,_:= hex.DecodeString("F6CEBE8030E5F7F7ED4ADAD55040EE9FDA382EF1")

	val1 := abci.Validator{
		Address:_address,
		Power:int64(1),
		PubKey: tmtypes.TM2PB.PubKey(_pubkey),

	}
	validatorSet:=[]abci.Validator{val1}
	fmt.Printf("all validators are %v \n",k.GetAllValidators(ctx))
	k.SetValidatorSet(ctx,validatorSet)
	fmt.Printf("all validators are %v \n",k.GetAllValidators(ctx))
	k.FlushValidatorSet(ctx)
	fmt.Printf("all validators are %v \n",k.GetAllValidators(ctx))
	return validatorSet
}