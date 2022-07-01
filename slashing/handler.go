package slashing

import (
	"bytes"
	"encoding/hex"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
)

// NewHandler creates an sdk.Handler for all the slashing type messages
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgTick:
			return handlerMsgTick(ctx, msg, k, contractCaller)
		case types.MsgTickAck:
			return handleMsgTickAck(ctx, msg, k, contractCaller)
		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in slashing module").Result()
		}
	}
}

// handlerMsgTick  - handles slashing of validators
// 0. check if slashLimit is exceeded or not.
// 1. Validate input slashing info hash data
// 2. If hash matches, copy slashBuffer into latestTickData
// 3. flushes slashBuffer, totalSlashedAmount
// 4. iterate and reduce the power of slashed validators.
// 5. Also update the jailStatus of Validator
// 6. emit event TickConfirmation
func handlerMsgTick(ctx sdk.Context, msg types.MsgTick, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating tick msg",
		"msgID", msg.ID,
		"SlashInfoBytes", msg.SlashingInfoBytes.String(),
	)

	// check if slash limit is exceeded or not
	totalSlashedAmount := k.GetTotalSlashedAmount(ctx)
	if totalSlashedAmount == 0 {
		k.Logger(ctx).Error("Slashed amount is zero")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Slashed amount is zero").Result()
	}

	// check if tick msgs are in continuity
	tickCount := k.GetTickCount(ctx)
	if msg.ID != tickCount+1 {
		k.Logger(ctx).Error("Tick not in countinuity", "msgID", msg.ID, "expectedMsgID", tickCount+1)
		return hmCommon.ErrTickNotInContinuity(k.Codespace()).Result()
	}

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

	// compare slashingInfoHash with msg hash
	k.Logger(ctx).Info("SlashInfo bytes generated", "SlashInfoBytes", hex.EncodeToString(slashingInfoBytes))

	if !bytes.Equal(slashingInfoBytes, msg.SlashingInfoBytes) {
		k.Logger(ctx).Error("slashingInfoBytes of current buffer state", "bufferSlashingInfoBytes", hex.EncodeToString(slashingInfoBytes),
			"doesn't match with slashingInfoBytes of msg", "msgSlashInfoBytes", msg.SlashingInfoBytes.String())
		return hmCommon.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("SlashInfoHash matches")

	// ensure latestTickData is empty
	tickSlashingInfos, err := k.GetTickValSlashingInfos(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching slash Info list from tick", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	if len(tickSlashingInfos) > 0 {
		k.Logger(ctx).Error("Waiting for tick data to be pushed to contract", "tickSlashingInfo", tickSlashingInfos)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg types.MsgUnjail, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating unjail msg",
		"validatorId", msg.ID,
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// pull validator from store
	validator, ok := k.sk.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	if !validator.Jailed {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

/*
	handleMsgTickAck - handle msg tick ack event
	1. validate the tx hash in the event
	2. flush the last tick slashing info
*/
func handleMsgTickAck(ctx sdk.Context, msg types.MsgTickAck, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating TickAck msg",
		"ID", msg.ID,
		"SlashedAmount", msg.SlashedAmount,
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// check if tick ack msgs are in continuity
	tickCount := k.GetTickCount(ctx)
	if msg.ID != tickCount {
		k.Logger(ctx).Error("Tick-ack not in countinuity", "msgID", msg.ID, "expectedMsgID", tickCount)
		return hmCommon.ErrTickAckNotInContinuity(k.Codespace()).Result()
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
