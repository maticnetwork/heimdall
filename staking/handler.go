package staking

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func NewHandler(k hmCommon.Keeper) sdk.Handler {
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

func handleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k hmCommon.Keeper) sdk.Result {
	//fetch validator from mainchain
	validator, err := helper.GetValidatorInfo(msg.ValidatorAddress)
	if err != nil || bytes.Equal(validator.Address.Bytes(), helper.ZeroAddress.Bytes()) {
		hmCommon.StakingLogger.Error(
			"Unable to fetch validator from rootchain",
			"error", err,
			"msgValidator", msg.ValidatorAddress.String(),
			"mainchainValidator", validator.Address.String(),
		)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}
	hmCommon.StakingLogger.Debug("Fetched validator from rootchain successfully", "validator", validator.String())

	// check validator address in message corresponds
	if !bytes.Equal(msg.ValidatorAddress.Bytes(), validator.Address.Bytes()) || msg.StartEpoch != validator.StartEpoch {
		hmCommon.StakingLogger.Error(
			"Validator address or startEpoch doesn't match",
			"msgValidator", msg.ValidatorAddress.String(),
			"mainchainValidator", validator.Address.String(),
			"msgStartEpoch", msg.StartEpoch,
			"mainchainStartEpoch", validator.StartEpoch,
		)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// Check if validator has been validator before
	if _, ok := k.GetSignerFromValidator(ctx, msg.ValidatorAddress); ok {
		hmCommon.StakingLogger.Error("Validator has been validator before, cannot join with same address", "presentValidator", msg.ValidatorAddress.String())
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace).Result()
	}

	// create new validator
	newValidator := hmTypes.Validator{
		Address:    msg.ValidatorAddress,
		StartEpoch: msg.StartEpoch,
		EndEpoch:   msg.EndEpoch,
		Power:      msg.GetPower(),
		PubKey:     pubkey,
		Signer:     signer,
	}

	// add validator to store
	hmCommon.StakingLogger.Debug("Adding new validator to state", "validator", newValidator.String())
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
		return hmCommon.ErrValidatorSave(k.Codespace).Result()
	}

	return sdk.Result{}
}

func handleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k hmCommon.Keeper) sdk.Result {
	var validator hmTypes.Validator

	// pull val from store
	if ok := k.GetValidatorFromValAddr(ctx, msg.ValidatorAddress, &validator); !ok {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	pubKey := msg.NewSignerPubKey
	newSigner := pubKey.Address()

	// check if updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		hmCommon.StakingLogger.Debug("Updating new signer", "signer", newSigner.String(), "validatorAddress", msg.ValidatorAddress.String())

		// update
		if err := k.UpdateSigner(ctx, newSigner, pubKey, validator.Signer); err != nil {
			hmCommon.StakingLogger.Error("Unable to update signer", "error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
			return hmCommon.ErrSignerUpdateError(k.Codespace).Result()
		}
	}

	// set new power
	validator.Power = msg.NewAmount // TODO change new amount to validator power
	// update validator to store with new power
	k.AddValidator(ctx, validator)

	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k hmCommon.Keeper) sdk.Result {
	var validator hmTypes.Validator

	// fetch validator from store
	if ok := k.GetValidatorFromValAddr(ctx, msg.ValidatorAddress, &validator); !ok {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		hmCommon.StakingLogger.Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace).Result()
	}

	// Add deactivation time for validator
	if err := k.AddDeactivationEpoch(ctx, validator); err != nil {
		hmCommon.StakingLogger.Error("Error while setting deactivation epoch to validator", "error", err, "validatorAddress", validator.Address.String())
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace).Result()
	}

	return sdk.Result{}
}
