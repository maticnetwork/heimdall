package staking

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "staking" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgValidatorJoin:
			return SideHandleMsgValidatorJoin(ctx, msg, k, contractCaller)
		case types.MsgValidatorExit:
			return SideHandleMsgValidatorExit(ctx, msg, k, contractCaller)
		case types.MsgSignerUpdate:
			return SideHandleMsgSignerUpdate(ctx, msg, k, contractCaller)
		case types.MsgStakeUpdate:
			return SideHandleMsgStakeUpdate(ctx, msg, k, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgValidatorJoin:
			return PostHandleMsgValidatorJoin(ctx, k, msg, sideTxResult)
		case types.MsgValidatorExit:
			return PostHandleMsgValidatorExit(ctx, k, msg, sideTxResult)
		case types.MsgSignerUpdate:
			return PostHandleMsgSignerUpdate(ctx, k, msg, sideTxResult)
		case types.MsgStakeUpdate:
			return PostHandleMsgStakeUpdate(ctx, k, msg, sideTxResult)
		default:
			return sdk.ErrUnknownRequest("Unrecognized Staking Msg type").Result()
		}
	}
}

// SideHandleMsgValidatorJoin side msg validator join
func SideHandleMsgValidatorJoin(ctx sdk.Context, msg types.MsgValidatorJoin, k Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for validator join msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	// decode validator join event
	eventLog, err := contractCaller.DecodeValidatorJoinEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	if !util.IsPubKeyFirstByteValid(pubkey[0:1]) {
		k.Logger(ctx).Error("public key first byte mismatch", "expected", "0x04", "received", pubkey[0:1])
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check signer pubkey in message corresponds
	if !bytes.Equal(pubkey.Bytes()[1:], eventLog.SignerPubkey) {
		k.Logger(ctx).Error(
			"Signer Pubkey does not match",
			"msgValidator", pubkey.String(),
			"mainchainValidator", hmTypes.BytesToHexBytes(eventLog.SignerPubkey),
		)

		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(signer.Bytes(), eventLog.Signer.Bytes()) {
		k.Logger(ctx).Error(
			"Signer Address from Pubkey does not match",
			"Validator", signer.String(),
			"mainchainValidator", eventLog.Signer.Hex(),
		)

		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check msg id
	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check ActivationEpoch
	if eventLog.ActivationEpoch.Uint64() != msg.ActivationEpoch {
		k.Logger(ctx).Error("ActivationEpoch in message doesn't match with ActivationEpoch in log", "msgActivationEpoch", msg.ActivationEpoch, "activationEpochFromTx", eventLog.ActivationEpoch.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check Amount
	if eventLog.Amount.Cmp(msg.Amount.BigInt()) != 0 {
		k.Logger(ctx).Error("Amount in message doesn't match Amount in event logs", "MsgAmount", msg.Amount, "AmountFromEvent", eventLog.Amount)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check Blocknumber
	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for validator join msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

// SideHandleMsgStakeUpdate handles stake update message
func SideHandleMsgStakeUpdate(ctx sdk.Context, msg types.MsgStakeUpdate, k Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for stake update msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	eventLog, err := contractCaller.DecodeValidatorStakeUpdateEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check Amount
	if eventLog.NewAmount.Cmp(msg.NewAmount.BigInt()) != 0 {
		k.Logger(ctx).Error("NewAmount in message doesn't match NewAmount in event logs", "MsgNewAmount", msg.NewAmount, "NewAmountFromEvent", eventLog.NewAmount)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for stake update msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

// SideHandleMsgSignerUpdate handles signer update message
func SideHandleMsgSignerUpdate(ctx sdk.Context, msg types.MsgSignerUpdate, k Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for signer update msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	eventLog, err := contractCaller.DecodeSignerUpdateEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if !util.IsPubKeyFirstByteValid(newPubKey.Bytes()[0:1]) {
		k.Logger(ctx).Error("public key first byte mismatch", "expected", "0x04", "received", newPubKey.Bytes()[0:1])
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if !bytes.Equal(eventLog.SignerPubkey, newPubKey.Bytes()[1:]) {
		k.Logger(ctx).Error("Newsigner pubkey in txhash and msg dont match", "msgPubKey", newPubKey.String(), "pubkeyTx", hmTypes.NewPubKey(eventLog.SignerPubkey[:]).String())
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(newSigner.Bytes(), eventLog.NewSigner.Bytes()) {
		k.Logger(ctx).Error("Signer Address from Pubkey does not match", "Validator", newSigner.String(), "mainchainValidator", eventLog.NewSigner.Hex())
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for signer update msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

// SideHandleMsgValidatorExit  handle  side msg validator exit
func SideHandleMsgValidatorExit(ctx sdk.Context, msg types.MsgValidatorExit, k Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for validator exit msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	// decode validator exit
	eventLog, err := contractCaller.DecodeValidatorExitEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.DeactivationEpoch.Uint64() != msg.DeactivationEpoch {
		k.Logger(ctx).Error("DeactivationEpoch in message doesn't match with deactivationEpoch in log", "msgDeactivationEpoch", msg.DeactivationEpoch, "deactivationEpochFromTx", eventLog.DeactivationEpoch.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for validator exit msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

/*
	Post Handlers - update the state of the tx
**/

// PostHandleMsgValidatorJoin msg validator join
func PostHandleMsgValidatorJoin(ctx sdk.Context, k Keeper, msg types.MsgValidatorJoin, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if validator join is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new validator-join since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Adding validator to state", "sideTxResult", sideTxResult)

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// get voting power from amount
	votingPower, err := helper.GetPowerFromAmount(msg.Amount.BigInt())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", msg.Amount, msg.ID)).Result()
	}

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          msg.ID,
		StartEpoch:  msg.ActivationEpoch,
		EndEpoch:    0,
		Nonce:       msg.Nonce,
		VotingPower: votingPower.Int64(),
		PubKey:      pubkey,
		Signer:      hmTypes.BytesToHeimdallAddress(signer.Bytes()),
		LastUpdated: "",
	}

	// update last updated
	newValidator.LastUpdated = sequence.String()

	// add validator to store
	k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())

	if err = k.AddValidator(ctx, newValidator); err != nil {
		k.Logger(ctx).Error("Unable to add validator to state", "validator", newValidator.String(), "error", err)
		return hmCommon.ErrValidatorSave(k.Codespace()).Result()
	}

	// Add Validator signing info. It is required for slashing module
	k.Logger(ctx).Debug("Adding signing info for new validator")

	valSigningInfo := hmTypes.NewValidatorSigningInfo(newValidator.ID, ctx.BlockHeight(), int64(0), int64(0))
	if err = k.AddValidatorSigningInfo(ctx, newValidator.ID, valSigningInfo); err != nil {
		k.Logger(ctx).Error("Unable to add validator signing info to state", "valSigningInfo", valSigningInfo.String(), "error", err)
		return hmCommon.ErrValidatorSigningInfoSave(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())
	k.Logger(ctx).Debug("✅ New validator successfully joined", "validator", strconv.FormatUint(newValidator.ID.Uint64(), 10))

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorJoin,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeyTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()), // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(newValidator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeySigner, newValidator.Signer.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgStakeUpdate handles stake update message
func PostHandleMsgStakeUpdate(ctx sdk.Context, k Keeper, msg types.MsgStakeUpdate, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if stakeUpdate is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping stake update since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Updating validator stake", "sideTxResult", sideTxResult)

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// update last updated
	validator.LastUpdated = sequence.String()

	// update nonce
	validator.Nonce = msg.Nonce

	// set validator amount
	p, err := helper.GetPowerFromAmount(msg.NewAmount.BigInt())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", msg.NewAmount, msg.ID)).Result()
	}

	validator.VotingPower = p.Int64()

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "ValidatorID", validator.ID, "error", err)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStakeUpdate,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgSignerUpdate handles signer update message
func PostHandleMsgSignerUpdate(ctx sdk.Context, k Keeper, msg types.MsgSignerUpdate, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if signer update is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping signer update since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))
	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting signer update", "sideTxResult", sideTxResult)

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	oldValidator := validator.Copy()

	// update last updated
	validator.LastUpdated = sequence.String()

	// update nonce
	validator.Nonce = msg.Nonce

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// Update signer in prev Validator
		validator.Signer = hmTypes.HeimdallAddress(newSigner)
		validator.PubKey = newPubKey

		k.Logger(ctx).Debug("Updating new signer", "newSigner", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	} else {
		k.Logger(ctx).Error("No signer change", "newSigner", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

	// remove old validator from HM
	oldValidator.EndEpoch = k.moduleCommunicator.GetACKCount(ctx)

	// remove old validator from TM
	oldValidator.VotingPower = 0
	// updated last
	oldValidator.LastUpdated = sequence.String()

	// updated nonce
	oldValidator.Nonce = msg.Nonce

	// save old validator
	if err := k.AddValidator(ctx, *oldValidator); err != nil {
		k.Logger(ctx).Error("Unable to update signer", "validatorId", validator.ID, "error", err)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// adding new validator
	k.Logger(ctx).Debug("Adding new validator", "validator", validator.String())

	// save validator
	err := k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "ValidatorID", validator.ID, "error", err)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	//
	// Move heimdall fee to new signer
	//

	// check if fee is already withdrawn
	coins := k.moduleCommunicator.GetCoins(ctx, oldValidator.Signer)

	maticBalance := coins.AmountOf(authTypes.FeeToken)
	if !maticBalance.IsZero() {
		k.Logger(ctx).Info("Transferring fee", "from", oldValidator.Signer.String(), "to", validator.Signer.String(), "balance", maticBalance.String())

		maticCoins := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: maticBalance}}
		if err := k.moduleCommunicator.SendCoins(ctx, oldValidator.Signer, validator.Signer, maticCoins); err != nil {
			k.Logger(ctx).Info("Error while transferring fee", "from", oldValidator.Signer.String(), "to", validator.Signer.String(), "balance", maticBalance.String())
			return err.Result()
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgValidatorExit handle msg validator exit
func PostHandleMsgValidatorExit(ctx sdk.Context, k Keeper, msg types.MsgValidatorExit, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if validator exit is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping validator exit since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting validator exit", "sideTxResult", sideTxResult)

	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorID", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// set end epoch
	validator.EndEpoch = msg.DeactivationEpoch

	// update last updated
	validator.LastUpdated = sequence.String()

	// update nonce
	validator.Nonce = msg.Nonce

	// Add deactivation time for validator
	if err := k.AddValidator(ctx, validator); err != nil {
		k.Logger(ctx).Error("Error while setting deactivation epoch to validator", "error", err, "validatorID", validator.ID.String())
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
