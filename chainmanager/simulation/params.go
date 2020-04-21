package simulation

import (
	"fmt"
	"math/rand"

	"github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/simulation"
	simtypes "github.com/maticnetwork/heimdall/types/simulation"
)

const (
	KeyTxConfirmationTime = "TxConfirmationTime"
	KeyChainParams        = "ChainParams"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, KeyTxConfirmationTime,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenTxConfirmationTime())
			},
		),
		simulation.NewSimParamChange(types.ModuleName, KeyChainParams,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBorChainId(r))
			},
		),
	}
}
