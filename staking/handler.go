package staking

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var ZeroBytes []byte = []byte("0x0000000000000000000000000000000000000000")

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
	if err != nil || bytes.Equal(validator.Address.Bytes(), ZeroBytes) {
		hmCommon.StakingLogger.Error("Unable to fetch validator from rootchain", "Error", err, "MsgValidator", msg.ValidatorAddress.String(), "MainchainValidator", validator.Address.String())
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}
	hmCommon.StakingLogger.Debug("Fetched validator from rootchain successfully", "Validator", validator.String())

	// Generate PubKey from Pubkey in message
	pubkey := helper.BytesToPubkey(msg.SignerPubKey)

	// check validator address in message corresponds
	if !bytes.Equal(msg.ValidatorAddress.Bytes(), validator.Address.Bytes()) {
		hmCommon.StakingLogger.Error("Validator address doesn't match", "MsgValidator", msg.ValidatorAddress.String(), "MainchainValidator", validator.Address.String())
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// Check if validator has been validator before
	var savedValidator hmTypes.Validator
	err = k.GetValidatorInfo(ctx, pubkey.Address().Bytes(), &savedValidator)
	if err == nil {
		hmCommon.StakingLogger.Error("Validator has been validator before, cannot join with same address", "PresentValidator", savedValidator.String())
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace).Result()
	}

	// add pubkey and override signer in validator
	validator.PubKey = pubkey
	validator.Signer = common.HexToAddress(pubkey.Address().String())

	// add validator to store
	hmCommon.StakingLogger.Debug("Adding new validator to state", "Validator", validator.String())
	err = k.AddValidator(ctx, validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to add validator to state", "Error", err, "Validator", validator.String())
		return hmCommon.ErrValidatorSave(k.Codespace).Result()
	}

	// add validator to validatorAddress => SignerAddress map
	k.SetValidatorAddrToSignerAddr(ctx, validator.Address, validator.Signer)

	return sdk.Result{}
}

func handleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k hmCommon.Keeper) sdk.Result {
	var validator hmTypes.Validator

	// pull val from store
	err := k.GetValidatorFromValAddr(ctx, msg.ValidatorAddress, &validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "error", err, "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	pubKey := helper.BytesToPubkey(msg.NewSignerPubKey)

	// check for already updated
	if bytes.Equal(pubKey.Address().Bytes(), validator.Signer.Bytes()) {
		hmCommon.StakingLogger.Error("No signer update on stakemanager found or signer already updated", "error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		return hmCommon.ErrValidatorAlreadySynced(k.Codespace).Result()
	}

	// fetch new signer
	newSigner := common.BytesToAddress(pubKey.Address().Bytes())

	// update
	err = k.UpdateSigner(ctx, newSigner, pubKey, validator.Signer)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to update signer", "Error", err, "currentSigner", validator.Signer.String(), "signerFromMsg", pubKey.Address().String())
		panic(err)
	}

	// update ValidatorAddress to SignerAddress Map
	k.SetValidatorAddrToSignerAddr(ctx, validator.Address, newSigner)

	return sdk.Result{}
}

func handleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k hmCommon.Keeper) sdk.Result {
	var validator hmTypes.Validator

	// fetch validator from store
	err := k.GetValidatorFromValAddr(ctx, msg.ValidatorAddress, &validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "error", err, "validatorAddress", msg.ValidatorAddress)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		hmCommon.StakingLogger.Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace).Result()
	}

	// means exit has been processed but validator in unbonding period
	// if validator.Power != 0 {
	// 	hmCommon.StakingLogger.Error("Validator already unbonded")
	// 	return hmCommon.ErrValUnbonded(k.Codespace).Result()
	// }

	// Add deactivation time for validator
	err = k.AddDeactivationEpoch(ctx, validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Error while setting deactivation epoch to validator", "error", err, "validatorAddress", validator.Address.String())
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace).Result()
	}

	return sdk.Result{}
}
