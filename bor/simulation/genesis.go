package simulation

import (
	"math/rand"
	"time"

	"github.com/maticnetwork/heimdall/bor/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// RandomizedGenState return dummy genesis
func RandomizedGenState(simState *module.SimulationState) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5
	params := types.DefaultParams()
	spans := []*hmTypes.Span{}
	chainID := "15001"
	start := uint64(0)
	end := uint64(0)
	accounts := simulation.RandomAccounts(r1, n)
	validators := make([]*hmTypes.Validator, n)
	producers := make([]hmTypes.Validator, n)

	for i := 0; i < len(validators); i++ {
		// validator
		val := hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
		validators[i] = val
		producers[0] = *val
	}
	validatorSet := hmTypes.NewValidatorSet(validators)

	for i := 0; i < n; i++ {
		start = end + 1
		end = end + 10
		span := hmTypes.NewSpan(
			uint64(i+1),
			start,
			end,
			*validatorSet,
			producers,
			chainID,
		)
		spans = append(spans, &span)
	}
	genesisState := types.NewGenesisState(
		params,
		spans,
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)

}
