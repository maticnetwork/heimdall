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
	// pull val from store
	validator, err := k.GetValidatorInfo(ctx, msg.CurrentValAddress)
	if err != nil {
		StakingLogger.Error("Fetching of validator from store failed", "Error", err, "ValidatorAddress", msg.CurrentValAddress)
		return ErrNoValidator(k.codespace).Result()
	}

	// pull val from mainchain
	newValidator, err := helper.GetValidatorInfo(msg.CurrentValAddress)
	if err != nil {
		StakingLogger.Error("Unable to fetch validator from stakemanager", "Error", err, "CurrentValidatorAddress", msg.CurrentValAddress)
	}

	pubkey, err := helper.StringToPubkey(msg.NewValPubkey)
	if err != nil {
		StakingLogger.Error("Invalid Pubkey", "Error", err, "PubkeyString", msg.NewValPubkey)
		return ErrValSignerMismatch(k.codespace).Result()
	}

	if !bytes.Equal(newValidator.Signer.Bytes(), validator.Signer.Bytes()) {
		StakingLogger.Error("No signer update on stakemanager found or signer already updated", "Error", err, "CurrentSigner", validator.Signer.String(), "SignerFromMsg", pubkey.Address().String())
		return ErrSignerAlreadySynced(k.codespace).Result()
	}

	// update
	err = k.UpdateSigner(ctx, newValidator.Signer, pubkey, msg.CurrentValAddress)
	if err != nil {
		StakingLogger.Error("Unable to update signer", "Error", err, "CurrentSigner", validator.Signer.String(), "SignerFromMsg", pubkey.Address().String())
		panic(err)
	}

	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k Keeper) sdk.Result {
	// fetch validator from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddr)
	if err != nil {
		StakingLogger.Error("Fetching of validator from store failed", "Error", err, "ValidatorAddress", msg.ValidatorAddr)
		return ErrNoValidator(k.codespace).Result()
	}

	// check if its validator exits in validator set
	if validator.IsCurrentValidator(k.checkpointKeeper.GetACKCount(ctx)) {
		StakingLogger.Error("Validator is locked in till deactivation period , exit denied")
		return ErrValIsCurrentVal(k.codespace).Result()
	}

	// check if validator is bonded
	if validator.Power == int64(0) {
		StakingLogger.Error("Validator already unbonded")
		return ErrValUnbonded(k.codespace).Result()
	}

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
