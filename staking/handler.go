package staking

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return handleMsgValidatorJoin(ctx, msg, k)
		case MsgValidatorExit:
			return handleMsgValidatorExit(ctx, msg, k)
		case MsgValidatorUpdate:
			return handleMsgValidatorUpdate(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}
func handleMsgValidatorUpdate(ctx sdk.Context, msg MsgValidatorUpdate, k Keeper) sdk.Result {

	// verify from mainchain
	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k Keeper) sdk.Result {
	// fetch validator from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddr)
	// check if its post endEpoch

	// verify deactivation from ACK count
	return sdk.Result{}
}

func handleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k Keeper) sdk.Result {
	// fetch validator from mainchain
	validator, err := helper.GetValidatorInfo(msg.ValidatorAddr)
	if err != nil {
		return ErrNoValidator(k.codespace).Result()
	}

	// validate if start epoch is after current tip
	ACKs := k.checkpointKeeper.GetACKCount(ctx)
	if int(validator.StartEpoch) < ACKs {
		// TODO add log
		return ErrOldValidator(k.codespace).Result()
	}

	// create crypto.pubkey from pubkey(string)
	pubkey, err := helper.StringToPubkey(msg.Pubkey)
	if err != nil {
		StakingLogger.Error("Invalid Pubkey", "Error", err, "PubkeyString", msg.Pubkey)
		return ErrValSignerMismatch(k.codespace).Result()
	}

	// check if the address of signer matches address from pubkey
	if !bytes.Equal(pubkey.Address().Bytes(), validator.Signer.Bytes()) {
		// TODO add log
		return ErrValSignerMismatch(k.codespace).Result()
	}

	// add pubkey generated to validator
	validator.Pubkey = pubkey

	// add validator to store
	k.AddValidator(ctx, validator)

	return sdk.Result{}
}
