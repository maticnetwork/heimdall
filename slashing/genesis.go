package slashing

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// stakingKeeper.IterateValidators(ctx,
	// 	func(index int64, validator exported.ValidatorI) bool {
	// 		keeper.AddPubkey(ctx, validator.GetConsPubKey())
	// 		return false
	// 	},
	// )

	for _, info := range data.SigningInfos {
		keeper.SetValidatorSigningInfo(ctx, info.ValID, info)
	}

	for valIDStr, array := range data.MissedBlocks {
		for _, missed := range array {
			valID, _ := strconv.ParseUint(valIDStr, 10, 64)
			keeper.SetValidatorMissedBlockBitArray(ctx, hmTypes.ValidatorID(valID), missed.Index, missed.Missed)
		}
	}

	for _, valSlashInfo := range data.BufferValSlashingInfo {
		keeper.SetBufferValSlashingInfo(ctx, valSlashInfo.ID, *valSlashInfo)
	}

	for _, tickValSlashInfo := range data.TickValSlashingInfo {
		keeper.SetTickValSlashingInfo(ctx, tickValSlashInfo.ID, *tickValSlashInfo)
	}

	keeper.SetParams(ctx, data.Params)

	// Set initial tick count
	keeper.UpdateTickCountWithValue(ctx, data.TickCount)

}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)
	signingInfos := make(map[string]hmTypes.ValidatorSigningInfo)
	missedBlocks := make(map[string][]types.MissedBlock)
	keeper.IterateValidatorSigningInfos(ctx, func(valID hmTypes.ValidatorID, info hmTypes.ValidatorSigningInfo) (stop bool) {
		signingInfos[valID.String()] = info
		localMissedBlocks := []types.MissedBlock{}

		keeper.IterateValidatorMissedBlockBitArray(ctx, valID, func(index int64, missed bool) (stop bool) {
			localMissedBlocks = append(localMissedBlocks, types.NewMissedBlock(index, missed))
			return false
		})
		missedBlocks[valID.String()] = localMissedBlocks
		return false
	})

	bufSlashInfos, _ := keeper.GetBufferValSlashingInfos(ctx)
	tickSlashInfos, _ := keeper.GetTickValSlashingInfos(ctx)
	return types.NewGenesisState(
		params,
		signingInfos,
		missedBlocks,
		bufSlashInfos,
		tickSlashInfos,
		keeper.GetTickCount(ctx))
}
