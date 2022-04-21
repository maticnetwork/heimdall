package simulation

// DONTCOVER

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// Simulation parameter constants
const (
	MaxMemoChars           = "max_memo_characters"
	TxSigLimit             = "tx_sig_limit"
	TxSizeCostPerByte      = "tx_size_cost_per_byte"
	SigVerifyCostED25519   = "sig_verify_cost_ed25519"
	SigVerifyCostSECP256K1 = "sig_verify_cost_secp256k1"

	MaxGas = "max_tx_gas"
	TxFees = "tx_fees"
)

// GenMaxMemoChars randomized MaxMemoChars
func GenMaxMemoChars(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 100, 200))
}

// GenTxSigLimit randomized TxSigLimit
// make sure that sigLimit is always high
// so that arbitrarily simulated messages from other
// modules can still create valid transactions
func GenTxSigLimit(r *rand.Rand) uint64 {
	return uint64(r.Intn(7) + 5)
}

// GenTxSizeCostPerByte randomized TxSizeCostPerByte
func GenTxSizeCostPerByte(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 5, 15))
}

// GenSigVerifyCostED25519 randomized SigVerifyCostED25519
func GenSigVerifyCostED25519(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 500, 1000))
}

// GenSigVerifyCostSECP256K1 randomized SigVerifyCostSECP256K1
func GenSigVerifyCostSECP256K1(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 500, 1000))
}

// GenMaxGas randomized MaxGas
func GenMaxGas(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 300000, 9000000))
}

// GenTxFees randomized GenTxFees
func GenTxFees(r *rand.Rand) string {
	base, _ := big.NewInt(0).SetString("1000000000000000", 10)
	t := simulation.RandIntBetween(r, 500, 1000)
	return big.NewInt(0).Mul(big.NewInt(0).SetInt64(int64(t)), base).String()
}

// RandomizedGenState generates a random GenesisState for auth
func RandomizedGenState(simState *module.SimulationState) {
	var maxMemoChars uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxMemoChars, &maxMemoChars, simState.Rand,
		func(r *rand.Rand) { maxMemoChars = GenMaxMemoChars(r) },
	)

	var txSigLimit uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TxSigLimit, &txSigLimit, simState.Rand,
		func(r *rand.Rand) { txSigLimit = GenTxSigLimit(r) },
	)

	var txSizeCostPerByte uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TxSizeCostPerByte, &txSizeCostPerByte, simState.Rand,
		func(r *rand.Rand) { txSizeCostPerByte = GenTxSizeCostPerByte(r) },
	)

	var sigVerifyCostED25519 uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SigVerifyCostED25519, &sigVerifyCostED25519, simState.Rand,
		func(r *rand.Rand) { sigVerifyCostED25519 = GenSigVerifyCostED25519(r) },
	)

	var sigVerifyCostSECP256K1 uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SigVerifyCostSECP256K1, &sigVerifyCostSECP256K1, simState.Rand,
		func(r *rand.Rand) { sigVerifyCostSECP256K1 = GenSigVerifyCostSECP256K1(r) },
	)

	var maxTxGas uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxGas, &maxTxGas, simState.Rand,
		func(r *rand.Rand) { maxTxGas = GenMaxGas(r) },
	)

	var txFees string
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TxFees, &txFees, simState.Rand,
		func(r *rand.Rand) { txFees = GenTxFees(r) },
	)

	params := types.NewParams(
		maxMemoChars,
		txSigLimit,
		txSizeCostPerByte,
		sigVerifyCostED25519,
		sigVerifyCostSECP256K1,

		maxTxGas,
		txFees,
	)
	genesisAccs := RandomGenesisAccounts(simState)

	authGenesis := types.NewGenesisState(params, genesisAccs)

	fmt.Printf("Selected randomly generated auth parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, authGenesis.Params))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(authGenesis)
}

// RandomGenesisAccounts returns randomly generated genesis accounts
func RandomGenesisAccounts(simState *module.SimulationState) (genesisAccs types.GenesisAccounts) {
	for _, acc := range simState.Accounts {
		bacc := types.NewBaseAccountWithAddress(acc.Address)
		bacc.Coins = simulation.RandomFeeCoins()
		gacc, _ := types.NewGenesisAccountI(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	return genesisAccs
}
