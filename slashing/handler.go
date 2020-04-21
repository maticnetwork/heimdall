package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
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

	// remove validator slashing infos from tick data
	k.FlushBufferValSlashingInfos(ctx)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handlerMsgTick  - handles slashing of validators
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
	valSlashingInfos := k.GetAllBufferValSlashingInfos(ctx)
	slashingInfoHash, err := GetSlashingInfoHash(valSlashingInfos)

	if err != nil {
		// compare slashingInfoHash with msg hash
		fmt.Println(slashingInfoHash)
	}

	// Validate nonce index of slashing message

	// ensure latestTickData is empty

	// copy slashBuffer into latestTickData
	k.CopyBufferValSlashingInfosToTickData(ctx)

	// Flush slashBuffer
	k.FlushBufferValSlashingInfos(ctx)

	// Flush TotalSlashedAmount

	// slash validator - Iterate lastTickData and reduce power of each validator along with jailing if needed
	k.SlashAndJailTickValSlashingInfos(ctx)

	// -slashing.  Emit TickConfirmation Event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTickConfirm,
			sdk.NewAttribute(types.AttributeKeySlashedAmount, fmt.Sprintf("%d", 10)),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// GetSlashingInfoHash returns hash of Validator slashing Info
func GetSlashingInfoHash(validatorSlashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {
	// TODO - Adding hashing logic
	return nil, nil
}
