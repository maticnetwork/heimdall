package simulation

import (
	"math/rand"
	"testing"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// GenRandomVal generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []hmTypes.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := hmCommonTypes.NewPubKey(privKey1.PubKey().Bytes())
		if randomise {
			startBlock := uint64(rand.Intn(10))
			// todo find a way to genrate non zero random number
			if startBlock == 0 {
				startBlock = 1
			}
			power := uint64(rand.Intn(100))
			if power == 0 {
				power = 1
			}
		}

		newVal := hmTypes.Validator{
			ID:               hmTypes.NewValidatorID(startID + uint64(i)),
			StartEpoch:       startBlock,
			EndEpoch:         startBlock + timeAlive,
			VotingPower:      power,
			Signer:           pubkey.Address().String(),
			PubKey:           pubkey.String(),
			ProposerPriority: 0,
		}
		validators = append(validators, newVal)
	}
	return
}

// LoadValidatorSet loads validator set
func LoadValidatorSet(count int, t *testing.T, keeper keeper.Keeper, ctx sdk.Context, randomise bool, timeAlive int) hmTypes.ValidatorSet {
	validators := GenRandomVal(count, 0, 10, uint64(timeAlive), randomise, 1)
	var valSet hmTypes.ValidatorSet

	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.NoError(t, err, "Unable to set validator, Error: %v", err)
		valSet.UpdateWithChangeSet([]*hmTypes.Validator{&validator})
	}

	err := keeper.UpdateValidatorSetInStore(ctx, &valSet)
	require.NoError(t, err, "Unable to update validator set")
	vals := keeper.GetAllValidators(ctx)
	require.NotNil(t, vals)
	return valSet
}
