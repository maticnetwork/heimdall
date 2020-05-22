package simulation

import (
	"github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/types/module"
)

// RandomizedGenState returns bank genesis
func RandomizedGenState(simState *module.SimulationState) {
	bankGenesis := types.NewGenesisState(true)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(bankGenesis)
}
