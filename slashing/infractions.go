package slashing

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// HandleValidatorSignature handles a validator signature, must be called once per validator per block.
func (k *Keeper) HandleValidatorSignature(ctx sdk.Context, addr []byte, power int64, signed bool) error {
	height := ctx.BlockHeight()
	signerAddress := hmTypes.BytesToHeimdallAddress(addr)
	k.Logger(ctx).Debug("Processing downtime request for validator", "address", signerAddress, "signed", signed, "power", power)

	// fetch validator Info
	validator, err := k.sk.GetValidatorInfo(ctx, signerAddress.Bytes())
	if err != nil {
		k.Logger(ctx).Error("validator info not found", "address", signerAddress)
		return err
	}

	// fetch signing info
	signInfo, found := k.GetValidatorSigningInfo(ctx, validator.ID)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", addr))
	}

	k.Logger(ctx).Debug("validator signing info", "valID", validator.ID, "address", signerAddress, "signingInfo", signInfo)

	params := k.GetParams(ctx)
	// this is a relative index, so it counts blocks the validator *should* have signed
	// will use the 0-value default signing info if not present, except for start height
	index := signInfo.IndexOffset % params.SignedBlocksWindow
	signInfo.IndexOffset++

	// Update signed block bit array & counter
	// This counter just tracks the sum of the bit array
	// That way we avoid needing to read/write the whole array each time
	previous := k.GetValidatorMissedBlockBitArray(ctx, validator.ID, index)
	k.Logger(ctx).Debug("validator signing status", "valID", validator.ID, "address", signerAddress, "previous", previous, "current", signed)
	missed := !signed
	switch {
	case !previous && missed:
		// Array value has changed from not missed to missed, increment counter
		k.SetValidatorMissedBlockBitArray(ctx, validator.ID, index, true)
		signInfo.MissedBlocksCounter++
		k.Logger(ctx).Debug("Array value has changed from not missed to missed, increment counter", "missedBlocksCounter", signInfo.MissedBlocksCounter)
	case previous && !missed:
		// Array value has changed from missed to not missed, decrement counter
		k.SetValidatorMissedBlockBitArray(ctx, validator.ID, index, false)
		signInfo.MissedBlocksCounter--
		k.Logger(ctx).Debug("Array value has changed from missed to not missed, decrement counter", "missedBlocksCounter", signInfo.MissedBlocksCounter)
	default:
		// Array value at this index has not changed, no need to update counter
		k.Logger(ctx).Debug("Array value has not changed. missedBlocksCounter remains same", "signingInfo", signInfo)
	}

	if missed {
		k.Logger(ctx).Info(
			fmt.Sprintf("Absent validator %s at height %d, %d missed, threshold %d", validator.ID, height, signInfo.MissedBlocksCounter, k.MinSignedPerWindow(ctx)))
	}

	minHeight := signInfo.StartHeight + params.SignedBlocksWindow
	maxMissed := params.SignedBlocksWindow - k.MinSignedPerWindow(ctx)

	// SLASH - if we are past the minimum height and the validator has missed too many blocks, punish them
	if height > minHeight && signInfo.MissedBlocksCounter > maxMissed {

		valSlashInfo, found := k.GetBufferValSlashingInfo(ctx, validator.ID)
		// if val is already in jailed state(in buffer or fixed), don't slash him anymore.
		if validator.Jailed || (found && valSlashInfo.IsJailed) {
			// Validator was (a) not found or (b) already jailed, don't slash
			k.Logger(ctx).Info(fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed", validator.ID))
		} else {
			// Downtime confirmed: slash and jail the validator
			k.Logger(ctx).Info(fmt.Sprintf("Validator %s past min height of %d and below signed blocks threshold of %d",
				validator.ID, minHeight, k.MinSignedPerWindow(ctx)))

			slashedAmount := k.SlashInterim(ctx, validator.ID, params.SlashFractionDowntime)
			k.Logger(ctx).Debug("Interim uptime slashing successful", "valID", validator.ID, "slashedAmount", slashedAmount)

			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
			signInfo.MissedBlocksCounter = 0
			signInfo.IndexOffset = 0
			k.clearValidatorMissedBlockBitArray(ctx, validator.ID)

		}
	}

	// Set the updated signing info
	k.SetValidatorSigningInfo(ctx, validator.ID, signInfo)
	return nil
}

// HandleDoubleSign implements an equivocation evidence handler. Assuming the
// evidence is valid, the validator committing the misbehavior will be slashed,
// jailed
//
// The evidence is considered invalid if:
// - the evidence is too old
// - the validator does not exist
// - the signing info does not exist (will panic)
// - is already jailed
func (k *Keeper) HandleDoubleSign(ctx sdk.Context, evidence types.Equivocation) error {
	consAddr := evidence.GetConsensusAddress()
	signerAddress := hmTypes.BytesToHeimdallAddress(consAddr)

	validator, err := k.sk.GetValidatorInfo(ctx, signerAddress.Bytes())
	if err != nil {
		k.Logger(ctx).Error("Error fetching validator", "signerAddress", signerAddress)
		return err
	}

	infractionHeight := evidence.GetHeight()
	k.Logger(ctx).Debug("Processing doubleSign request for validator", "address", signerAddress, "height", infractionHeight)

	// calculate the age of the evidence
	blockTime := ctx.BlockHeader().Time
	age := blockTime.Sub(evidence.GetTime())
	params := k.GetParams(ctx)

	// reject evidence if the double-sign is too old
	if age > params.MaxEvidenceAge {
		k.Logger(ctx).Error("Ignored double sign from %s at height %d, age of %d past max age of %d",
			signerAddress, infractionHeight, age, params.MaxEvidenceAge)
		return errors.New("double sign too old")
	}

	if ok := k.HasValidatorSigningInfo(ctx, validator.ID); !ok {
		panic(fmt.Sprintf("expected signing info for validator %s but not found", validator.ID))
	}

	k.Logger(ctx).Info(fmt.Sprintf("confirmed double sign from %s at height %d, age of %d", validator.ID, infractionHeight, age))

	// Slash validator. The `power` is the int64 power of the validator as provided
	// to/by Tendermint.
	valSlashInfo, found := k.GetBufferValSlashingInfo(ctx, validator.ID)
	// if val is already in jailed state(in buffer or fixed), don't slash him anymore.
	if validator.Jailed || (found && valSlashInfo.IsJailed) {
		// Validator was (a) not found or (b) already jailed, don't slash
		k.Logger(ctx).Info(fmt.Sprintf("Validator %s would have been slashed for double time, but was either not found in store or already jailed", validator.ID))
	} else {
		slashedAmount := k.SlashInterim(ctx, validator.ID, params.SlashFractionDoubleSign)
		k.Logger(ctx).Debug("Interim uptime slashing successful", "valID", validator.ID, "slashedAmount", slashedAmount)
	}

	return nil
}
