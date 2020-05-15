package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/slashing/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// Simulation parameter constants
const (
	SignedBlocksWindow      = "signed_blocks_window"
	MinSignedPerWindow      = "min_signed_per_window"
	DowntimeJailDuration    = "downtime_jail_duration"
	SlashFractionDoubleSign = "slash_fraction_double_sign"
	SlashFractionDowntime   = "slash_fraction_downtime"
	SlashFractionLimit      = "slash_fraction_limit"
	JailFractionLimit       = "jail_fraction_limit"
	MaxEvidenceAge          = "max_evidence_age"
)

// GenSignedBlocksWindow randomized SignedBlocksWindow
func GenSignedBlocksWindow(r *rand.Rand) int64 {
	return int64(simulation.RandIntBetween(r, 10, 1000))
}

// GenMinSignedPerWindow randomized MinSignedPerWindow
func GenMinSignedPerWindow(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(10)), 1)
}

// GenDowntimeJailDuration randomized DowntimeJailDuration
func GenDowntimeJailDuration(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24)) * time.Second
}

// GenSlashFractionDoubleSign randomized SlashFractionDoubleSign
func GenSlashFractionDoubleSign(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(50) + 1)))
}

// GenSlashFractionDowntime randomized SlashFractionDowntime
func GenSlashFractionDowntime(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(200) + 1)))
}

// GenSlashFractionLimit randomized SlashFractionLimit
func GenSlashFractionLimit(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(200) + 1)))
}

// GenJailFractionLimit randomized JailFractionLimit
func GenJailFractionLimit(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(1).Quo(sdk.NewDec(int64(r.Intn(200) + 1)))
}

// GenMaxEvidenceAge randomized MaxEvidenceAge
func GenMaxEvidenceAge(r *rand.Rand) time.Duration {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24)) * time.Second
}

// GenMaxEvidenceAge randomized MaxEvidenceAge
func GenEnableslashing(r *rand.Rand) bool {
	return (r.Intn(200)%10 == 0)
}

// RandomizedGenState generates a random GenesisState for slashing
func RandomizedGenState(simState *module.SimulationState) {
	var signedBlocksWindow int64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SignedBlocksWindow, &signedBlocksWindow, simState.Rand,
		func(r *rand.Rand) { signedBlocksWindow = GenSignedBlocksWindow(r) },
	)

	var minSignedPerWindow sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinSignedPerWindow, &minSignedPerWindow, simState.Rand,
		func(r *rand.Rand) { minSignedPerWindow = GenMinSignedPerWindow(r) },
	)

	var downtimeJailDuration time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, DowntimeJailDuration, &downtimeJailDuration, simState.Rand,
		func(r *rand.Rand) { downtimeJailDuration = GenDowntimeJailDuration(r) },
	)

	var slashFractionDoubleSign sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFractionDoubleSign, &slashFractionDoubleSign, simState.Rand,
		func(r *rand.Rand) { slashFractionDoubleSign = GenSlashFractionDoubleSign(r) },
	)

	var slashFractionDowntime sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFractionDowntime, &slashFractionDowntime, simState.Rand,
		func(r *rand.Rand) { slashFractionDowntime = GenSlashFractionDowntime(r) },
	)

	var slashFractionLimit sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SlashFractionLimit, &slashFractionLimit, simState.Rand,
		func(r *rand.Rand) { slashFractionLimit = GenSlashFractionLimit(r) },
	)

	var jailFractionLimit sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, JailFractionLimit, &jailFractionLimit, simState.Rand,
		func(r *rand.Rand) { jailFractionLimit = GenJailFractionLimit(r) },
	)

	var maxEvidenceAge time.Duration
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxEvidenceAge, &maxEvidenceAge, simState.Rand,
		func(r *rand.Rand) { maxEvidenceAge = GenMaxEvidenceAge(r) },
	)

	var enableSlashing bool
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxEvidenceAge, &maxEvidenceAge, simState.Rand,
		func(r *rand.Rand) { enableSlashing = GenEnableslashing(r) },
	)

	params := types.NewParams(
		signedBlocksWindow, minSignedPerWindow, downtimeJailDuration,
		slashFractionDoubleSign, slashFractionDowntime, slashFractionLimit, jailFractionLimit, maxEvidenceAge, enableSlashing,
	)

	slashingGenesis := types.NewGenesisState(params, nil, nil, nil, nil, uint64(0))

	fmt.Printf("Selected randomly generated slashing parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, slashingGenesis.Params))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(slashingGenesis)
}
