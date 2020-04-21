package slashing

import (
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

	for addr, info := range data.SigningInfos {
		address := hmTypes.HexToHeimdallAddress(addr)
		keeper.SetValidatorSigningInfo(ctx, address.Bytes(), info)
	}

	for addr, array := range data.MissedBlocks {
		address := hmTypes.HexToHeimdallAddress(addr)
		for _, missed := range array {
			keeper.SetValidatorMissedBlockBitArray(ctx, address.Bytes(), missed.Index, missed.Missed)
		}
	}

	keeper.SetParams(ctx, data.Params)

}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)
	signingInfos := make(map[string]hmTypes.ValidatorSigningInfo)
	missedBlocks := make(map[string][]types.MissedBlock)
	keeper.IterateValidatorSigningInfos(ctx, func(address []byte, info hmTypes.ValidatorSigningInfo) (stop bool) {
		bechAddr := hmTypes.BytesToHeimdallAddress(address).String()
		signingInfos[bechAddr] = info
		localMissedBlocks := []types.MissedBlock{}

		keeper.IterateValidatorMissedBlockBitArray(ctx, address, func(index int64, missed bool) (stop bool) {
			localMissedBlocks = append(localMissedBlocks, types.NewMissedBlock(index, missed))
			return false
		})
		missedBlocks[bechAddr] = localMissedBlocks

		return false
	})

	return types.NewGenesisState(params, signingInfos, missedBlocks)
}
