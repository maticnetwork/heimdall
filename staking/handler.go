package staking

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func NewHandler(k hmCommon.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return HandleMsgValidatorJoin(ctx, msg, k, contractCaller)
		case MsgValidatorExit:
			return HandleMsgValidatorExit(ctx, msg, k, contractCaller)
		case MsgSignerUpdate:
			return HandleMsgSignerUpdate(ctx, msg, k, contractCaller)
		case MsgPowerUpdate:
			return HandleMsgPowerUpdate(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func HandleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k hmCommon.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	hmCommon.StakingLogger.Debug("Handing new validator join", "Msg", msg)

	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash); !confirmed {
		return hmCommon.ErrWaitFrConfirmation(k.Codespace).Result()
	}

	//fetch validator from mainchain
	validator, err := contractCaller.GetValidatorInfo(msg.ID)
	if err != nil || bytes.Equal(validator.Signer.Bytes(), helper.ZeroAddress.Bytes()) {
		hmCommon.StakingLogger.Error(
			"Unable to fetch validator from rootchain",
			"error", err,
			"msgValidator", msg.ID,
			"mainChainSigner", validator.Signer.String(),
		)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}
	hmCommon.StakingLogger.Debug("Fetched validator from rootchain successfully", "validator", validator.String())

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// check signer in message corresponds
	if !bytes.Equal(signer.Bytes(), validator.Signer.Bytes()) {
		hmCommon.StakingLogger.Error(
			"Signer Address does not match",
			"msgValidator", msg.SignerPubKey.Address().String(),
			"mainchainValidator", validator.Signer.String(),
		)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// Check if validator has been validator before
	if _, ok := k.GetSignerFromValidatorID(ctx, msg.ID); ok {
		hmCommon.StakingLogger.Error("Validator has been validator before, cannot join with same ID", "valID", msg.ID)
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace).Result()
	}

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          validator.ID,
		StartEpoch:  validator.StartEpoch,
		EndEpoch:    validator.EndEpoch,
		Power:       validator.Power,
		PubKey:      pubkey,
		Signer:      validator.Signer,
		LastUpdated: big.NewInt(0),
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

// Handle signer update message
func HandleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k hmCommon.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	hmCommon.StakingLogger.Debug("Handling signer update", "Validator", msg.ID, "Signer", msg.NewSignerPubKey.Address())

	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash); !confirmed {
		return hmCommon.ErrWaitFrConfirmation(k.Codespace).Result()
	}

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	id, newSignerTx, _, err := contractCaller.SigUpdateEvent(msg.TxHash)
	if err != nil {
		hmCommon.StakingLogger.Error("Error fetching log from txhash", "Error", err)
		return hmCommon.ErrInvalidMsg(k.Codespace, "Unable to fetch logs for txHash. Error: %v", err).Result()
	}

	if int(id) != msg.ID.Int() {
		hmCommon.StakingLogger.Error("ID in message doesnt match id in logs", "MsgID", msg.ID, "IdFromTx", id)
		return hmCommon.ErrInvalidMsg(k.Codespace, "Invalid txhash, id's dont match. Id from tx hash is %v", id).Result()
	}

	if bytes.Compare(newSignerTx.Bytes(), newSigner.Bytes()) != 0 {
		hmCommon.StakingLogger.Error("Signer in txhash and msg dont match", "MsgSigner", newSigner.String(), "SignerTx", newSignerTx.String())
		return hmCommon.ErrInvalidMsg(k.Codespace, "Signer in txhash and msg dont match").Result()
	}

	// pull val from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "validatorAddress", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}
	oldValidator := validator.Copy()

	// check if txhash has been used before
	blockNum, _ := contractCaller.GetBlockNoFromTxHash(msg.TxHash)
	if err := k.SetLastUpdated(ctx, msg.ID, &blockNum); err != nil {
		hmCommon.StakingLogger.Error("Error occured while updating last updated", "Error", err)
		return err.Result()
	}

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// Update signer in prev Validator
		validator.Signer = newSigner
		validator.PubKey = newPubKey
		hmCommon.StakingLogger.Debug("Updating new signer", "signer", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	}

	hmCommon.StakingLogger.Error("Removing old validator", "Validator", oldValidator.String())
	// remove old validator from HM
	oldValidator.EndEpoch = k.GetACKCount(ctx)
	// remove old validator from TM
	oldValidator.Power = 0
	err = k.AddValidator(ctx, *oldValidator)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace).Result()
	}

	hmCommon.StakingLogger.Error("Adding new validator", "Validator", validator.String())
	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace).Result()
	}

	return sdk.Result{}
}

// handle validator exit transactions
func HandleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k hmCommon.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	hmCommon.StakingLogger.Info("Handling validator exit", "ValidatorID", msg.ID)

	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash); !confirmed {
		return hmCommon.ErrWaitFrConfirmation(k.Codespace).Result()
	}
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "validatorID", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	hmCommon.StakingLogger.Debug("validator in store", "validator", validator)
	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		hmCommon.StakingLogger.Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace).Result()
	}

	// get validator from mainchain
	updatedVal, err := contractCaller.GetValidatorInfo(validator.ID)
	if err != nil {
		hmCommon.StakingLogger.Error("Cannot fetch validator info while unstaking", "Error", err, "validatorID", validator.ID)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}

	// Add deactivation time for validator
	if err := k.AddDeactivationEpoch(ctx, validator, updatedVal); err != nil {
		hmCommon.StakingLogger.Error("Error while setting deactivation epoch to validator", "error", err, "validatorID", validator.ID)
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace).Result()
	}

	return sdk.Result{}
}

// handle power update transactions
func HandleMsgPowerUpdate(ctx sdk.Context, msg MsgPowerUpdate, k hmCommon.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	hmCommon.StakingLogger.Info("Handling power update", "ValidatorID", msg.ID)

	// make sure tx is has received specified confirmation
	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash); !confirmed {
		hmCommon.StakingLogger.Error("Tx has not received enough confirmations", "ReqConf", helper.GetConfig().ConfirmationBlocks, "Tx", msg.TxHash)
		return hmCommon.ErrWaitFrConfirmation(k.Codespace).Result()
	}

	// get event from tx
	valID, power, err := contractCaller.PowerUpdateEvent(msg.TxHash)
	if err != nil {
		hmCommon.StakingLogger.Error("Error fetching log from txhash", "Error", err)
		return hmCommon.ErrInvalidMsg(k.Codespace, "Unable to fetch logs for txHash. Error: %v", err).Result()
	}

	// check if validator id in event is correct
	if valID != msg.ID {
		// return id mismatch error
	}

	// get validator information from validator ID
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		hmCommon.StakingLogger.Error("Fetching of validator from store failed", "validatorID", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace).Result()
	}
	if validator.Power != power {
		hmCommon.StakingLogger.Info("Updating validator power", "OldPower", validator.Power, "NewPower", power)
	}
	validator.Power = power

	err = k.AddValidator(ctx, validator)
	if err != nil {
		hmCommon.StakingLogger.Error("Unable to update power", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace).Result()
	}

	// update last updated field as well

	return sdk.Result{}

}
