package simulation

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingKeeper "github.com/maticnetwork/heimdall/x/staking/keeper"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"

	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/common"
	"github.com/stretchr/testify/require"
)

// GenRandCheckpoint return headers
func GenRandCheckpoint(start uint64, headerSize uint64, maxCheckpointLength uint64) (headerBlock types.Checkpoint, err error) {
	end := start + headerSize
	borChainID := "1234"
	rootHash := common.HexToHeimdallHash("123")
	proposer := common.HeimdallAddress{}

	headerBlock = types.CreateBlock(
		start,
		end,
		rootHash,
		proposer,
		borChainID,
		uint64(time.Now().UTC().Unix()))

	return headerBlock, nil
}

// LoadValidatorSet loads validator set
func LoadValidatorSet(count int, t *testing.T, keeper stakingKeeper.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	validators := stakingSim.GenRandomVal(count, 0, 10, uint64(timeAlive), randomise, 1)
	var valSet types.ValidatorSet

	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.NoError(t, err, "Unable to set validator, Error: %v", err)
		valSet.UpdateWithChangeSet([]*types.Validator{&validator})
	}

	err := keeper.UpdateValidatorSetInStore(ctx, &valSet)
	require.NoError(t, err, "Unable to update validator set")
	vals := keeper.GetAllValidators(ctx)
	require.NotNil(t, vals)
	return valSet
}
