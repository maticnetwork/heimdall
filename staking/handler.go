package staking

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmcmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k hmcmn.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return handleMsgValidatorJoin(ctx, msg, k)
		case MsgValidatorExit:
			return handleMsgValidatorExit(ctx, msg, k)
		case MsgSignerUpdate:
			return handleMsgSignerUpdate(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func handleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k hmcmn.Keeper) sdk.Result {
	// fetch validator from mainchain
	validator, err := helper.GetValidatorInfo(msg.ValidatorAddress)
	if err != nil {
		return hmcmn.ErrNoValidator(k.Codespace).Result()
	}

	// validate if start epoch is after current tip
	ackCount := k.GetACKCount(ctx)
	if int(validator.StartEpoch) < ackCount {
		// TODO add log
		return hmcmn.ErrOldValidator(k.Codespace).Result()
	}

	pubKey := helper.BytesToPubkey(msg.ValidatorPubKey)
	// check if the address of signer matches address from pubkey
	if !bytes.Equal(pubKey.Address().Bytes(), validator.Signer.Bytes()) {
		// TODO add log
		return hmcmn.ErrValSignerMismatch(k.Codespace).Result()
	}

	// add pubkey generated to validator
	validator.PubKey = pubKey

	// add validator to store
	k.AddValidator(ctx, validator)

	// validator set changed
	k.SetValidatorSetChangedFlag(ctx, true)

	return sdk.Result{}
}

func handleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k hmcmn.Keeper) sdk.Result {
	// pull val from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddress)
	if err != nil {
		hmcmn.StakingLogger.Error("Fetching of validator from store failed", "error", err, "validatorAddress", msg.ValidatorAddress)
		return hmcmn.ErrNoValidator(k.Codespace).Result()
	}

	// pull val from mainchain
	newValidator, err := helper.GetValidatorInfo(msg.ValidatorAddress)
	if err != nil {
		hmcmn.StakingLogger.Error("Unable to fetch validator from stakemanager", "error", err, "currentValidatorAddress", msg.ValidatorAddress)
	}

	pubKey := helper.BytesToPubkey(msg.NewValidatorPubKey)
	// check if the address of signer matches address from pubkey
	if !bytes.Equal(pubKey.Address().Bytes(), newValidator.Signer.Bytes()) {
		// TODO add log
		return hmcmn.ErrValSignerMismatch(k.Codespace).Result()
	}

	// check for already updated
	if !bytes.Equal(newValidator.Signer.Bytes(), validator.Signer.Bytes()) {
		hmcmn.StakingLogger.Error("No signer update on stakemanager found or signer already updated", "error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		return hmcmn.ErrValidatorAlreadySynced(k.Codespace).Result()
	}

	// update
	err = k.UpdateSigner(ctx, newValidator.Signer, pubKey, msg.ValidatorAddress)
	if err != nil {
		hmcmn.StakingLogger.Error("Unable to update signer", "Error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		panic(err)
	}

	// TODO: Not sure how to communicate signer changes to TM

	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k hmcmn.Keeper) sdk.Result {
	// fetch validator from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddr)
	if err != nil {
		hmcmn.StakingLogger.Error("Fetching of validator from store failed", "Error", err, "ValidatorAddress", msg.ValidatorAddr)
		return hmcmn.ErrNoValidator(k.Codespace).Result()
	}

	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		hmcmn.StakingLogger.Error("Validator already unbonded")
		return hmcmn.ErrValUnbonded(k.Codespace).Result()
	}

	// allow only validators to exit from validator set
	if !validator.IsCurrentValidator(k.GetACKCount(ctx)) {
		hmcmn.StakingLogger.Error("Validator is not in validator set, exit not possible")
		return hmcmn.ErrValIsNotCurrentVal(k.Codespace).Result()
	}

	// means exit has been processed but validator in unbonding period
	if validator.Power != int64(0) {
		hmcmn.StakingLogger.Error("Validator already unbonded")
		return hmcmn.ErrValUnbonded(k.Codespace).Result()
	}

	// Add deactivation time for validator
	k.AddDeactivationEpoch(ctx, msg.ValidatorAddr, validator)

	return sdk.Result{}
}
