package simulation

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

func RandomizedGenState(simState *module.SimulationState) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5
	accounts := simulation.RandomAccounts(r1, n)
	stakingSequence := make([]string, n)

	validators := make([]*hmTypes.Validator, n)

	for i := range stakingSequence {
		stakingSequence[i] = strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	}

	for i := 0; i < len(validators); i++ {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
	}

	// validator set
	validatorSet := hmTypes.NewValidatorSet(validators)

	genesisState := types.NewGenesisState(validators, *validatorSet, stakingSequence)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
