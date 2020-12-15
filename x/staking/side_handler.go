package staking

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	"github.com/maticnetwork/heimdall/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "staking" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgValidatorJoin:
			return SideHandleMsgValidatorJoin(ctx, *msg, k, contractCaller)
		case *types.MsgValidatorExit:
			return SideHandleMsgValidatorExit(ctx, *msg, k, contractCaller)
		case *types.MsgSignerUpdate:
			return SideHandleMsgSignerUpdate(ctx, *msg, k, contractCaller)
		case *types.MsgStakeUpdate:
			return SideHandleMsgStakeUpdate(ctx, *msg, k, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(6), // TODO should be changed like `sdk.CodeUnknownRequest`
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgValidatorJoin:
			return PostHandleMsgValidatorJoin(ctx, k, *msg, sideTxResult)
		case *types.MsgValidatorExit:
			return PostHandleMsgValidatorExit(ctx, k, *msg, sideTxResult)
		case *types.MsgSignerUpdate:
			return PostHandleMsgSignerUpdate(ctx, k, *msg, sideTxResult)
		case *types.MsgStakeUpdate:
			return PostHandleMsgStakeUpdate(ctx, k, *msg, sideTxResult)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// SideHandleMsgValidatorJoin side msg validator join
func SideHandleMsgValidatorJoin(ctx sdk.Context, msg types.MsgValidatorJoin, k keeper.Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ Validating External call for validator join msg",
		"txHash", hmCommonTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// TODO uncomment and fix the issue
	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// decode validator join event
	eventLog, err := contractCaller.DecodeValidatorJoinEvent(chainParams.StakingInfoAddress, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// check signer pubkey in message corresponds
	if !bytes.Equal(pubkey.Bytes()[1:], eventLog.SignerPubkey) {
		k.Logger(ctx).Error(
			"Signer Pubkey does not match",
			"msgValidator", pubkey.String(),
			"mainchainValidator", hmTypes.BytesToHexBytes(eventLog.SignerPubkey),
		)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(signer.Bytes(), eventLog.Signer.Bytes()) {
		k.Logger(ctx).Error(
			"Signer Address from Pubkey does not match",
			"Validator", signer.String(),
			"mainchainValidator", eventLog.Signer.Hex(),
		)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check msg id
	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check ActivationEpoch
	if eventLog.ActivationEpoch.Uint64() != msg.ActivationEpoch {
		k.Logger(ctx).Error("ActivationEpoch in message doesn't match with ActivationEpoch in log", "msgActivationEpoch", msg.ActivationEpoch, "activationEpochFromTx", eventLog.ActivationEpoch.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check Amount
	if eventLog.Amount.Cmp(msg.Amount.BigInt()) != 0 {
		k.Logger(ctx).Error("Amount in message doesn't match Amount in event logs", "MsgAmount", msg.Amount, "AmountFromEvent", eventLog.Amount)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check Blocknumber
	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for validator join msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

// SideHandleMsgStakeUpdate handles stake update message
func SideHandleMsgStakeUpdate(ctx sdk.Context, msg types.MsgStakeUpdate, k keeper.Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for stake update msg",
		"txHash", hmCommonTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	eventLog, err := contractCaller.DecodeValidatorStakeUpdateEvent(chainParams.StakingInfoAddress, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check Amount
	if eventLog.NewAmount.Cmp(msg.NewAmount.BigInt()) != 0 {
		k.Logger(ctx).Error("NewAmount in message doesn't match NewAmount in event logs", "MsgNewAmount", msg.NewAmount, "NewAmountFromEvent", eventLog.NewAmount)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for stake update msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

// SideHandleMsgSignerUpdate handles signer update message
func SideHandleMsgSignerUpdate(ctx sdk.Context, msg types.MsgSignerUpdate, k keeper.Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for signer update msg",
		"txHash", hmCommonTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	eventLog, err := contractCaller.DecodeSignerUpdateEvent(chainParams.StakingInfoAddress, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if bytes.Compare(eventLog.SignerPubkey, newPubKey.Bytes()[1:]) != 0 {
		k.Logger(ctx).Error("Newsigner pubkey in txhash and msg dont match", "msgPubKey", newPubKey.String(), "pubkeyTx", hmCommonTypes.NewPubKey(eventLog.SignerPubkey[:]).String())
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(newSigner.Bytes(), eventLog.NewSigner.Bytes()) {
		k.Logger(ctx).Error("Signer Address from Pubkey does not match", "Validator", newSigner.String(), "mainchainValidator", eventLog.NewSigner.Hex())
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for signer update msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

// SideHandleMsgValidatorExit  handle  side msg validator exit
func SideHandleMsgValidatorExit(ctx sdk.Context, msg types.MsgValidatorExit, k keeper.Keeper, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for validator exit msg",
		"txHash", hmCommonTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// decode validator exit
	eventLog, err := contractCaller.DecodeValidatorExitEvent(chainParams.StakingInfoAddress, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	if eventLog.DeactivationEpoch.Uint64() != msg.DeactivationEpoch {
		k.Logger(ctx).Error("DeactivationEpoch in message doesn't match with deactivationEpoch in log", "msgDeactivationEpoch", msg.DeactivationEpoch, "deactivationEpochFromTx", eventLog.DeactivationEpoch.Uint64)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	// check nonce
	if eventLog.Nonce.Uint64() != msg.Nonce {
		k.Logger(ctx).Error("Nonce in message doesn't match with nonce in log", "msgNonce", msg.Nonce, "nonceFromTx", eventLog.Nonce)
		return hmCommon.ErrorSideTx(hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for validator exit msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

/*
	Post Handlers - update the state of the tx
**/

// PostHandleMsgValidatorJoin msg validator join
func PostHandleMsgValidatorJoin(ctx sdk.Context, k keeper.Keeper, msg types.MsgValidatorJoin, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {

	// Skip handler if validator join is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping new validator-join since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Adding validator to state", "sideTxResult", sideTxResult)

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address().String()

	// get voting power from amount
	votingPower, err := helper.GetPowerFromAmount(msg.Amount.BigInt())
	if err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          msg.ID,
		StartEpoch:  msg.ActivationEpoch,
		EndEpoch:    0,
		Nonce:       msg.Nonce,
		VotingPower: votingPower.Int64(),
		PubKey:      pubkey.String(),
		Signer:      signer,
		LastUpdated: "",
	}

	// update last updated
	newValidator.LastUpdated = sequence.String()

	// add validator to store
	k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		k.Logger(ctx).Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
		return nil, hmCommon.ErrValidatorSave
	}

	// Add Validator signing info. It is required for slashing module
	k.Logger(ctx).Debug("Adding signing info for new validator")
	valSigningInfo := hmTypes.NewValidatorSigningInfo(newValidator.ID, ctx.BlockHeight(), int64(0), int64(0))
	err = k.AddValidatorSigningInfo(ctx, newValidator.ID, valSigningInfo)
	if err != nil {
		k.Logger(ctx).Error("Unable to add validator signing info to state", "error", err, "valSigningInfo", valSigningInfo.String())
		return nil, hmCommon.ErrValidatorSigningInfoSave
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
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeyTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()), // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(newValidator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeySigner, newValidator.Signer),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &sdk.Result{}, nil
}

// PostHandleMsgStakeUpdate handles stake update message
func PostHandleMsgStakeUpdate(ctx sdk.Context, k keeper.Keeper, msg types.MsgStakeUpdate, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	// Skip handler if stakeUpdate is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping stake update since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Updating validator stake", "sideTxResult", sideTxResult)

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return nil, hmCommon.ErrNoValidator
	}

	// update last updated
	validator.LastUpdated = sequence.String()

	// update nonce
	validator.Nonce = msg.Nonce

	// set validator amount
	p, err := helper.GetPowerFromAmount(msg.NewAmount.BigInt())
	if err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}
	validator.VotingPower = p.Int64()

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return nil, hmCommon.ErrSignerUpdateError
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
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &sdk.Result{}, nil
}

// PostHandleMsgSignerUpdate handles signer update message
func PostHandleMsgSignerUpdate(ctx sdk.Context, k keeper.Keeper, msg types.MsgSignerUpdate, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	// Skip handler if signer update is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping signer update since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))
	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Persisting signer update", "sideTxResult", sideTxResult)

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return nil, hmCommon.ErrNoValidator
	}
	oldValidator := validator.Copy()

	// update last udpated
	validator.LastUpdated = sequence.String()

	// update nonce
	validator.Nonce = msg.Nonce

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), []byte(validator.Signer)) {
		// Update signer in prev Validator
		validator.Signer = newSigner.String()
		validator.PubKey = newPubKey.String()
		k.Logger(ctx).Debug("Updating new signer", "newSigner", newSigner, "oldSigner", oldValidator.Signer, "validatorID", msg.ID)
	} else {
		k.Logger(ctx).Error("No signer change", "newSigner", newSigner, "oldSigner", oldValidator.Signer, "validatorID", msg.ID)
		return nil, hmCommon.ErrSignerUpdateError
	}

	k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

	// remove old validator from HM
	oldValidator.EndEpoch = k.ModuleCommunicator.GetACKCount(ctx)

	// remove old validator from TM
	oldValidator.VotingPower = 0
	// updated last
	oldValidator.LastUpdated = sequence.String()

	// updated nonce
	oldValidator.Nonce = msg.Nonce

	// save old validator
	if err := k.AddValidator(ctx, *oldValidator); err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "validatorId", validator.ID)
		return nil, hmCommon.ErrSignerUpdateError
	}

	// adding new validator
	k.Logger(ctx).Debug("Adding new validator", "validator", validator.String())

	// save validator
	err := k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return nil, hmCommon.ErrSignerUpdateError
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
	coins := k.ModuleCommunicator.GetCoins(ctx, oldValidator.Signer)
	maticBalance := coins.AmountOf(types.FeeToken)
	if !maticBalance.IsZero() {
		k.Logger(ctx).Info("Transferring fee", "from", oldValidator.Signer, "to", validator.Signer, "balance", maticBalance.String())
		maticCoins := sdk.Coins{sdk.Coin{Denom: types.FeeToken, Amount: maticBalance}}
		if err := k.ModuleCommunicator.SendCoins(ctx, oldValidator.Signer, validator.Signer, maticCoins); err != nil {
			k.Logger(ctx).Info("Error while transferring fee", "from", oldValidator.Signer, "to", validator.Signer, "balance", maticBalance.String())
			return nil, err
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &sdk.Result{}, nil
}

// PostHandleMsgValidatorExit handle msg validator exit
func PostHandleMsgValidatorExit(ctx sdk.Context, k keeper.Keeper, msg types.MsgValidatorExit, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	// Skip handler if validator exit is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping validator exit since side-tx didn't get yes votes")
		return nil, common.ErrSideTxValidation
	}

	// Check for replay attack
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Persisting validator exit", "sideTxResult", sideTxResult)

	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorID", msg.ID)
		return nil, hmCommon.ErrNoValidator
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
		return nil, hmCommon.ErrValidatorNotDeactivated
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &sdk.Result{}, nil
}
