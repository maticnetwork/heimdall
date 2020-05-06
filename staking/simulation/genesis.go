package simulation

import (
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5
	accounts := simulation.RandomAccounts(r1, n)
	stakingSequence := make([]string, n)

	validators := make([]*hmTypes.Validator, n)
	dividendAccounts := make([]hmTypes.DividendAccount, n)

	for i := range stakingSequence {
		stakingSequence[i] = strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	}

	for i := 0; i < len(validators); i++ {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)

		// create dividend account for validator
		dividendAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.NewDividendAccountID(uint64(validators[i].ID)),
			big.NewInt(0).String(),
			big.NewInt(0).String(),
		)
	}

	// validator set
	validatorSet := hmTypes.NewValidatorSet(validators)

	genesisState := types.NewGenesisState(validators, *validatorSet, dividendAccounts, stakingSequence)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
