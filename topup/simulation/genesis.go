package simulation

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// SequenceNumber ...
var SequenceNumber = "sequence_number"

// GenSequenceNumber returns randomc chainID
func GenSequenceNumber(r *rand.Rand) string {
	return strconv.Itoa(simulation.RandIntBetween(r, 0, 1000000000000))
}

// RandomizeGenState returns topup genesis
func RandomizeGenState(simState *module.SimulationState) {
	var sequences []string
	dividendAccounts := make([]hmTypes.DividendAccount, 5)

	for i := 0; i < 5; i++ {
		var sequenceNumber string
		simState.AppParams.GetOrGenerate(simState.Cdc, SequenceNumber, &sequenceNumber, simState.Rand,
			func(r *rand.Rand) {
				sequenceNumber = GenSequenceNumber(r)
			},
		)
		sequences = append(sequences, sequenceNumber)

		// create dividend account for validator
		dividendAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.NewDividendAccountID(uint64(i)),
			big.NewInt(0).String(),
		)
	}

	topupGenesis := types.NewGenesisState(sequences, dividendAccounts)
	fmt.Printf("Selected randomly generated topup sequences:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, topupGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(topupGenesis)
}
