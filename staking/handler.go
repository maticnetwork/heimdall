package staking

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

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
	if err != nil {
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// validate if start epoch is after current tip
	// ackCount := k.GetACKCount(ctx)
	// if int(validator.StartEpoch) < ackCount {
	// 	// TODO add log
	// 	return hmCommon.ErrOldValidator(k.Codespace).Result()
	// }

	pubKey := helper.BytesToPubkey(msg.ValidatorPubKey)

	// check if the address of signer matches address from pubkey
	// if !bytes.Equal(pubKey.Address().Bytes(), validator.Signer.Bytes()) {
	// 	// TODO add log
	// 	return hmCommon.ErrValSignerMismatch(k.Codespace).Result()
	// }

	// add pubkey generated to validator
	// validator.PubKey = pubKey

	if false {
		if !bytes.Equal(validator.Address.Bytes(), msg.ValidatorAddress.Bytes()) ||
			validator.Power != msg.Amount ||
			validator.StartEpoch != msg.StartEpoch {
			// TODO revert if mainchain doesn't match with incoming data
		}
	}

	var savedValidator hmTypes.Validator
	err := k.GetValidator(ctx, msg.ValidatorAddress.Bytes(), &savedValidator)
	if err != nil {
		//
		// Ignore if not present
	} else if savedValidator.Power != 0 {
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace).Result()
	}

	newValidator := hmTypes.Validator{
		Address:    msg.ValidatorAddress,
		Power:      msg.Amount,
		StartEpoch: msg.StartEpoch,
		EndEpoch:   0,
		PubKey:     pubKey,
		Signer:     common.BytesToAddress(pubKey.Address().Bytes()),
	}

	// add validator to store
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		return hmCommon.ErrValidatorSave(k.Codespace).Result()
	}

	// validator set changed
	k.SetValidatorSetChangedFlag(ctx, true)

	return sdk.Result{}
}

func handleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k hmCommon.Keeper) sdk.Result {
	// pull val from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddress)
	if err != nil {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "error", err, "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// pull val from mainchain
	newValidator, err := helper.GetValidatorInfo(msg.ValidatorAddress)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to fetch validator from stakemanager", "error", err, "currentValidatorAddress", msg.ValidatorAddress)
	}

	pubKey := helper.BytesToPubkey(msg.NewValidatorPubKey)
	// check if the address of signer matches address from pubkey
	if !bytes.Equal(pubKey.Address().Bytes(), newValidator.Signer.Bytes()) {
		// TODO add log
		return hmCommon.ErrValSignerMismatch(k.Codespace).Result()
	}

	// check for already updated
	if !bytes.Equal(newValidator.Signer.Bytes(), validator.Signer.Bytes()) {
		hmCommon.StakingLogger.Error("No signer update on stakemanager found or signer already updated", "error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		return hmCommon.ErrValidatorAlreadySynced(k.Codespace).Result()
	}

	// update
	err = k.UpdateSigner(ctx, newValidator.Signer, pubKey, msg.ValidatorAddress)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to update signer", "Error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		panic(err)
	}

	// TODO: Not sure how to communicate signer changes to TM

	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k hmCommon.Keeper) sdk.Result {
	// fetch validator from store
	validator, err := k.GetValidatorInfo(ctx, msg.ValidatorAddress)
	if err != nil {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "error", err, "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		hmCommon.StakingLogger.Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace).Result()
	}

	// allow only validators to exit from validator set
	if !validator.IsCurrentValidator(k.GetACKCount(ctx)) {
		hmCommon.StakingLogger.Error("Validator is not in validator set, exit not possible")
		return hmCommon.ErrValIsNotCurrentVal(k.Codespace).Result()
	}

	// means exit has been processed but validator in unbonding period
	if validator.Power != 0 {
		hmCommon.StakingLogger.Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace).Result()
	}

	// Add deactivation time for validator
	k.AddDeactivationEpoch(ctx, msg.ValidatorAddress, validator)

	return sdk.Result{}
}
