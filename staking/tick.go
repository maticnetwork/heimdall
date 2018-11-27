package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/common"
)

// EndBlocker refreshes validator set after block commit
func EndBlocker(ctx sdk.Context, k common.Keeper) (validators []abci.ValidatorUpdate) {
	// todo revive this when we get ACK in endBlock
	//StakingLogger.Info("Current validators fetched", "validators", helper.ValidatorsToString(k.GetAllValidators(ctx)))
	//
	//// flush exiting validator set
	//k.FlushValidatorSet(ctx)
	//// fetch current validator set
	//validatorSet := helper.GetValidators()
	//// update
	//k.SetValidatorSet(ctx, validatorSet)
	//
	//StakingLogger.Info("New validators set", "validators", helper.ValidatorsToString(k.GetAllValidators(ctx)))
	//return validatorSet
	return
}
