package slashing

import (
	"bytes"
	"encoding/hex"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "topup" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgTick:
			return SideHandleMsgTick(ctx, k, msg, contractCaller)
		case types.MsgTickAck:
			return SideHandleMsgTickAck(ctx, k, msg, contractCaller)
		case types.MsgUnjail:
			return SideHandleMsgUnjail(ctx, k, msg, contractCaller)
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
		case types.MsgTick:
			return PostHandleMsgTick(ctx, k, msg, sideTxResult)
		case types.MsgTickAck:
			return PostHandleMsgTickAck(ctx, k, msg, sideTxResult)
		case types.MsgUnjail:
			return PostHandleMsgUnjail(ctx, k, msg, sideTxResult)
		default:
			errMsg := "Unrecognized slash Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// SideHandleMsgTick handles MsgTick message for external call
func SideHandleMsgTick(ctx sdk.Context, k Keeper, msg types.MsgTick, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for tick msg")
	k.Logger(ctx).Debug("✅ Succesfully validated External call for tick msg")
	result.Result = abci.SideTxResultType_Yes
	return
}

// SideHandleMsgTick handles MsgTick message for external call
func SideHandleMsgTickAck(ctx sdk.Context, k Keeper, msg types.MsgTickAck, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for tick-ack msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
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

	// get event log for slashed event
	eventLog, err := contractCaller.DecodeSlashedEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.Amount.Uint64() != msg.SlashedAmount {
		k.Logger(ctx).Error("SlashedAmount in message doesn't match SlashedAmount in event logs", "MsgSlashedAmount", msg.SlashedAmount, "SlashedAmountFromEvent", eventLog.Amount)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for tick-ack msg")
	result.Result = abci.SideTxResultType_Yes
	return
}

// SideHandleMsgUnjail handles MsgUnjail message for external call
func SideHandleMsgUnjail(ctx sdk.Context, k Keeper, msg types.MsgUnjail, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for unjail msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
		"validatorID", msg.ID,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	// get unjail event
	eventLog, err := contractCaller.DecodeUnJailedEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for tick msg")
	result.Result = abci.SideTxResultType_Yes
	return
}

// PostHandleMsgTick  - handles slashing of validators
// 1. copy slashBuffer into latestTickData
// 2. flush slashBuffer, totalSlashedAmount
// 3. iterate and reduce the power of slashed validators.
// 4. Also update the jailStatus of Validator
// 5. emit event TickConfirmation
func PostHandleMsgTick(ctx sdk.Context, k Keeper, msg types.MsgTick, sideTxResult abci.SideTxResultType) sdk.Result {

	// Skip handler if tick is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new tick since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay - tick should be in conitunity
	tickCount := k.GetTickCount(ctx)
	if msg.ID != tickCount+1 {
		k.Logger(ctx).Error("Tick not in countinuity. may be due to replay", "msgID", msg.ID, "expectedMsgID", tickCount+1)
		return hmCommon.ErrTickNotInContinuity(k.Codespace()).Result()
	}

	// check if state has changed between handler and side-handler blocks
	valSlashingInfos, err := k.GetBufferValSlashingInfos(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching slash Info list from buffer", "error", err)
		return hmCommon.ErrSlashInfoDetails(k.Codespace()).Result()
	}
	slashingInfoBytes, err := types.SortAndRLPEncodeSlashInfos(valSlashingInfos)
	if err != nil {
		k.Logger(ctx).Info("Error generating slashing info bytes", "error", err)
		return hmCommon.ErrSlashInfoDetails(k.Codespace()).Result()
	}
	k.Logger(ctx).Info("SlashInfo bytes generated", "SlashInfoBytes", hex.EncodeToString(slashingInfoBytes))
	// compare slashingInfoHash with msg hash
	if !bytes.Equal(slashingInfoBytes, msg.SlashingInfoBytes) {
		k.Logger(ctx).Error("slashingInfoBytes of current buffer state", "bufferSlashingInfoBytes", hex.EncodeToString(slashingInfoBytes),
			"doesn't match with slashingInfoBytes of msg", "msgSlashInfoBytes", msg.SlashingInfoBytes.String())
		return hmCommon.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting tick state", "sideTxResult", sideTxResult)

	// copy slashBuffer into latestTickData
	if err := k.CopyBufferValSlashingInfosToTickData(ctx); err != nil {
		k.Logger(ctx).Error("Error copying bufferSlashInfo to tickSlashInfo", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// Flush slashBuffer
	if err := k.FlushBufferValSlashingInfos(ctx); err != nil {
		k.Logger(ctx).Error("Error flushing buffer slash info in tick handler", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// Flush TotalSlashedAmount
	k.FlushTotalSlashedAmount(ctx)

	// Update Tick count
	k.IncrementTickCount(ctx)

	k.Logger(ctx).Debug("Successfully slashed and jailed")
	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTickConfirm,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeySlashInfoBytes, msg.SlashingInfoBytes.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgTickAck  - handles slashing of validators
/*
	handleMsgTickAck - handle msg tick ack event
	1. validate the tx hash in the event
	2. flush the last tick slashing info
*/
func PostHandleMsgTickAck(ctx sdk.Context, k Keeper, msg types.MsgTickAck, sideTxResult abci.SideTxResultType) sdk.Result {

	// Skip handler if topup is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new topup since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check if tick ack msgs are in continuity.
	tickCount := k.GetTickCount(ctx)
	if msg.ID != tickCount {
		k.Logger(ctx).Error("Tick-ack not in countinuity.", "msgID", msg.ID, "expectedMsgID", tickCount)
		return hmCommon.ErrTickAckNotInContinuity(k.Codespace()).Result()
	}

	// check for replay -  check if incoming tx is older
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	tickSlashInfos, err := k.GetTickValSlashingInfos(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching tick slash infos", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	if tickSlashInfos == nil || len(tickSlashInfos) == 0 {
		k.Logger(ctx).Error("tick ack already processed", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// slash validator - Iterate lastTickData and reduce power of each validator along with jailing if needed
	if err := k.SlashAndJailTickValSlashingInfos(ctx); err != nil {
		k.Logger(ctx).Error("Error slashing and jailing validator in tick handler", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// remove validator slashing infos from tick data
	if err := k.FlushTickValSlashingInfos(ctx); err != nil {
		k.Logger(ctx).Error("Error flushing tick slash info in tick-ack handler", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Successfully flushed tick slash info in tick-ack handler")

	// save staking sequence
	k.SetSlashingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTickAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func PostHandleMsgUnjail(ctx sdk.Context, k Keeper, msg types.MsgUnjail, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if topup is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new topup since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check if incoming tx is older
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting jail status", "sideTxResult", sideTxResult)

	// unjail validator
	k.sk.Unjail(ctx, msg.ID)

	// check if unjail is successful or not
	val, _ := k.sk.GetValidatorFromValID(ctx, msg.ID)
	if val.Jailed {
		k.Logger(ctx).Error("Error unjailing validator", "validatorId", msg.ID, "jailStatus", val.Jailed)
		return hmCommon.ErrUnjailValidator(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetSlashingSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnjail,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
