package slashing

import (
	"bytes"
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
		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k, contractCaller)
		case types.MsgTick:
			return handlerMsgTick(ctx, msg, k, contractCaller)
		case types.MsgTickAck:
			return handleMsgTickAck(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in slashing module").Result()
		}
	}
}

// Validators must submit a transaction to unjail itself after
// having been jailed (and thus unbonded) for downtime
func handleMsgUnjail(ctx sdk.Context, msg types.MsgUnjail, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	// err := k.Unjail(ctx, msg.ValidatorAddr)
	// if err != nil {
	// 	return nil, err
	// }

	k.Logger(ctx).Info("Handling unjail event", "msg", msg)
	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash(), params.TxConfirmationTime)
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace(), params.TxConfirmationTime).Result()
	}

	// decode unjail event
	eventLog, err := contractCaller.DecodeUnJailedEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	// sequence id
	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUnjail,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
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
	// err := k.Unjail(ctx, msg.ValidatorAddr)
	// if err != nil {
	// 	return nil, err
	// }

	// check if slash limit is exceeded or not
	if !k.IsSlashedLimitExceeded(ctx) {
		k.Logger(ctx).Error("TotalSlashedAmount is less than SlashLimit")
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	valSlashingInfos, err := k.GetBufferValSlashingInfos(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching slash Info list from buffer", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	slashingInfoHash, err := types.GetSlashingInfoHash(valSlashingInfos)
	if err != nil {
		k.Logger(ctx).Info("Error generating slashing info hash", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// compare slashingInfoHash with msg hash
	k.Logger(ctx).Info("SlashInfo hash generated", "SlashInfoHash", hmTypes.BytesToHeimdallHash(slashingInfoHash).String())

	if !bytes.Equal(slashingInfoHash, msg.SlashingInfoHash.Bytes()) {
		k.Logger(ctx).Error("SlashInfoHash of current buffer state", "bufferSlashInfoHash", hmTypes.BytesToHeimdallHash(slashingInfoHash).String(),
			"doesn't match with SlashInfoHash of msg", "msgSlashInfoHash", msg.SlashingInfoHash)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("SlashInfoHash matches")

	// ensure latestTickData is empty
	tickSlashingInfos, err := k.GetTickValSlashingInfos(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching slash Info list from tick", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	if tickSlashingInfos != nil && len(tickSlashingInfos) > 0 {
		k.Logger(ctx).Error("Waiting for tick data to be pushed to contract", "tickSlashingInfo", tickSlashingInfos)
	}

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

	// slash validator - Iterate lastTickData and reduce power of each validator along with jailing if needed
	if err := k.SlashAndJailTickValSlashingInfos(ctx); err != nil {
		k.Logger(ctx).Error("Error slashing and jailing validator in tick handler", "error", err)
		return common.ErrSlashInfoDetails(k.Codespace()).Result()
	}

	// -slashing.  Emit TickConfirmation Event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTickConfirm,
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeySlashInfoHash, msg.SlashingInfoHash.String()),
		),
	)

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
	k.Logger(ctx).Info("Handling TickAck", "msg", msg)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash(), params.TxConfirmationTime)
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace(), params.TxConfirmationTime).Result()
	}

	// get event log for slashed event
	eventLog, err := contractCaller.DecodeSlashedEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasSlashingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// remove validator slashing infos from tick data
	k.FlushBufferValSlashingInfos(ctx)

	// save staking sequence
	k.SetSlashingSequence(ctx, sequence.String())

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTickAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
