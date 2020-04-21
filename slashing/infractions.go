package slashing

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr []byte, power int64, signed bool) {
	logger := k.Logger(ctx)
	height := ctx.BlockHeight()
	address := hmTypes.BytesToHeimdallAddress(addr)

	/* 	// fetch the validator public key
	   	consAddr := sdk.ConsAddress(addr)
	   	if _, err := k.GetPubkey(ctx, addr); err != nil {
	   		panic(fmt.Sprintf("Validator consensus-address %s not found", consAddr))
	   	} */

	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, addr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", addr))
	}

	params := k.GetParams(ctx)
	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := signInfo.IndexOffset % params.SignedBlocksWindow
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedBlockBitArray(ctx, addr, index)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedBlockBitArray(ctx, addr, index, true)
		signInfo.MissedBlocksCounter++
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedBlockBitArray(ctx, addr, index, false)
		signInfo.MissedBlocksCounter--
	default:
		// Array value at this index has not changed, no need to update counter
	}

	if missed {
		fmt.Println("Entered23")
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
				sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", height)),
			),
		)

		logger.Info(
			fmt.Sprintf("Absent validator %s at height %d, %d missed, threshold %d", address, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}

	minHeight := signInfo.StartHeight + params.SignedBlocksWindow
	maxMissed := params.SignedBlocksWindow - k.MinSignedPerWindow(ctx)

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {
		validator, err := k.sk.GetValidatorInfo(ctx, addr)
		if err != nil {
			logger.Error("Error fetching validator")
		}
		logger.Error("Error fetching validator", validator)
		if err == nil /* && !validator.IsJailed()  */ {

			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				address, minHeight, k.MinSignedPerWindow(ctx)))

			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the LastCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			// That's fine since this is just used to filter unbonding delegations & redelegations.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeSlash,
					sdk.NewAttribute(types.AttributeKeyAddress, address.String()),
					sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
					sdk.NewAttribute(types.AttributeKeyJailed, address.String()),
				),
			)

			// update slash buffer present in slash keeper. Also add slashedAmount totalSlashedAmount.
			val, _ := k.sk.GetValidatorInfo(ctx, signInfo.Signer.Bytes())
			amount := ""
			k.SlashInterim(ctx, val.ID, amount)
			
			k.sk.Slash(ctx, addr, distributionHeight, power, params.SlashFractionDowntime)
			k.sk.Jail(ctx, addr)

			signInfo.JailedUntil = ctx.BlockHeader().Time.Add(params.DowntimeJailDuration)

			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearValidatorMissedBlockBitArray(ctx, addr)
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed", address),
			)
		}
	}

	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, addr, signInfo)
}

// HandleDoubleSign implements an equivocation evidence handler. Assuming the
// evidence is valid, the validator committing the misbehavior will be slashed,
// jailed and tombstoned. Once tombstoned, the validator will not be able to
// recover. Note, the evidence contains the block time and height at the time of
// the equivocation.
//
// The evidence is considered invalid if:
// - the evidence is too old
// - the validator is unbonded or does not exist
// - the signing info does not exist (will panic)
// - is already tombstoned
//
// TODO: Some of the invalid constraints listed above may need to be reconsidered
// in the case of a lunatic attack.
func (k Keeper) HandleDoubleSign(ctx sdk.Context, evidence types.Equivocation) {
	logger := k.Logger(ctx)
	consAddr := evidence.GetConsensusAddress()
	infractionHeight := evidence.GetHeight()

	// calculate the age of the evidence
	blockTime := ctx.BlockHeader().Time
	age := blockTime.Sub(evidence.GetTime())
	params := k.GetParams(ctx)

	// if _, err := k.slashingKeeper.GetPubkey(ctx, consAddr.Bytes()); err != nil {
	// 	// Ignore evidence that cannot be handled.
	// 	//
	// 	// NOTE: We used to panic with:
	// 	// `panic(fmt.Sprintf("Validator consensus-address %v not found", consAddr))`,
	// 	// but this couples the expectations of the app to both Tendermint and
	// 	// the simulator.  Both are expected to provide the full range of
	// 	// allowable but none of the disallowed evidence types.  Instead of
	// 	// getting this coordination right, it is easier to relax the
	// 	// constraints and ignore evidence that cannot be handled.
	// 	return
	// }

	// reject evidence if the double-sign is too old
	if age > params.MaxEvidenceAge {
		logger.Info(
			fmt.Sprintf(
				"ignored double sign from %s at height %d, age of %d past max age of %d",
				consAddr, infractionHeight, age, params.MaxEvidenceAge,
			),
		)
		return
	}

	// validator, _ := k.sk.GetValidatorInfo(ctx, consAddr)
	// TODO - slashing
	// if validator == nil || validator.IsUnbonded() {
	// 	// Defensive: Simulation doesn't take unbonding periods into account, and
	// 	// Tendermint might break this assumption at some point.
	// 	return
	// }

	if ok := k.HasValidatorSigningInfo(ctx, consAddr); !ok {
		panic(fmt.Sprintf("expected signing info for validator %s but not found", consAddr))
	}

	// ignore if the validator is already tombstoned
	// TODO - slashing
	/* 	if k.IsTombstoned(ctx, consAddr) {
		logger.Info(
			fmt.Sprintf(
				"ignored double sign from %s at height %d, validator already tombstoned",
				consAddr, infractionHeight,
			),
		)
		return
	} */

	logger.Info(fmt.Sprintf("confirmed double sign from %s at height %d, age of %d", consAddr, infractionHeight, age))

	// We need to retrieve the stake distribution which signed the block, so we
	// subtract ValidatorUpdateDelay from the evidence height.
	// Note, that this *can* result in a negative "distributionHeight", up to
	// -ValidatorUpdateDelay, i.e. at the end of the
	// pre-genesis block (none) = at the beginning of the genesis block.
	// That's fine since this is just used to filter unbonding delegations & redelegations.
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay

	// Slash validator. The `power` is the int64 power of the validator as provided
	// to/by Tendermint. This value is validator.Tokens as sent to Tendermint via
	// ABCI, and now received as evidence. The fraction is passed in to separately
	// to slash unbonding and rebonding delegations.
	k.Slash(
		ctx,
		consAddr,
		params.SlashFractionDoubleSign,
		evidence.GetValidatorPower(), distributionHeight,
	)

	// Jail the validator if not already jailed. This will begin unbonding the
	// validator if not already unbonding (tombstoned).
	// TODO - slashing
	/* 	if !validator.IsJailed() {
	   		k.Jail(ctx, consAddr)
	   	}

	   	k.JailUntil(ctx, consAddr, params.DoubleSignJailEndTime)
	   	k.slashingKeeper.Tombstone(ctx, consAddr) */
}
