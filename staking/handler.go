package staking

import (
	"bytes"
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
func handleMsgValidatorUpdate(context sdk.Context, update MsgValidatorUpdate, keeper Keeper) sdk.Result {
	// verify from mainchain
	return sdk.Result{}
}
func handleMsgValidatorExit(context sdk.Context, exit MsgValidatorExit, keeper Keeper) sdk.Result {
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
	var pubkeyBytes secp256k1.PubKeySecp256k1
	_pubkey, err := hex.DecodeString(msg.Pubkey)
	if err != nil {
		StakingLogger.Error("Decoding of pubkey(string) to pubkey failed", "Error", err, "PubkeyString", msg.Pubkey)
	}
	copy(pubkeyBytes[:], _pubkey)

	// check if the address of signer matches address from pubkey
	if !bytes.Equal(pubkeyBytes.Address().Bytes(), validator.Signer.Bytes()) {
		// TODO add log
		return ErrValSignerMismatch(k.codespace).Result()
	}

	// add pubkey generated to validator
	validator.Pubkey = pubkeyBytes

	// add validator to store
	k.AddValidator(ctx, validator)

	return sdk.Result{}
}
