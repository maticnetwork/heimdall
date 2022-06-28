package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/chainmanager/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// Parameter keys
const (
	MainchainTxConfirmations  = "mainchain_tx_confirmations"
	MaticchainTxConfirmations = "maticchain_tx_confirmations"

	BorChainID            = "bor_chain_id"
	MaticTokenAddress     = "matic_token_address"     //nolint
	StakingManagerAddress = "staking_manager_address" //nolint
	SlashManagerAddress   = "slash_manager_address"   //nolint
	RootChainAddress      = "root_chain_address"      //nolint
	StakingInfoAddress    = "staking_info_address"    //nolint
	StateSenderAddress    = "state_sender_address"    //nolint

	// Bor Chain Contracts
	StateReceiverAddress = "state_receiver_address" //nolint
	ValidatorSetAddress  = "validator_set_address"  //nolint
)

func GenMainchainTxConfirmations(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1, 100))
}

func GenMaticchainTxConfirmations(r *rand.Rand) uint64 {
	return uint64(simulation.RandIntBetween(r, 1, 100))
}

func GenHeimdallAddress() hmTypes.HeimdallAddress {
	return hmTypes.BytesToHeimdallAddress(simulation.RandHex(20))
}

// GenBorChainId returns randomc chainID
func GenBorChainId(r *rand.Rand) string {
	return strconv.Itoa(simulation.RandIntBetween(r, 0, math.MaxInt32))
}

func RandomizedGenState(simState *module.SimulationState) {
	var mainchainTxConfirmations uint64

	simState.AppParams.GetOrGenerate(simState.Cdc, MainchainTxConfirmations, &mainchainTxConfirmations, simState.Rand,
		func(r *rand.Rand) { mainchainTxConfirmations = GenMainchainTxConfirmations(r) },
	)

	var (
		maticchainTxConfirmations uint64
		borChainID                string
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, MaticchainTxConfirmations, &maticchainTxConfirmations, simState.Rand,
		func(r *rand.Rand) { maticchainTxConfirmations = GenMaticchainTxConfirmations(r) },
	)

	simState.AppParams.GetOrGenerate(simState.Cdc, BorChainID, &borChainID, simState.Rand,
		func(r *rand.Rand) { borChainID = GenBorChainId(r) },
	)

	var (
		maticTokenAddress     = GenHeimdallAddress()
		stakingManagerAddress = GenHeimdallAddress()
		slashManagerAddress   = GenHeimdallAddress()
		rootChainAddress      = GenHeimdallAddress()
		stakingInfoAddress    = GenHeimdallAddress()
		stateSenderAddress    = GenHeimdallAddress()
		stateReceiverAddress  = GenHeimdallAddress()
		validatorSetAddress   = GenHeimdallAddress()
	)

	chainParams := types.ChainParams{
		BorChainID:            borChainID,
		MaticTokenAddress:     maticTokenAddress,
		StakingManagerAddress: stakingManagerAddress,
		SlashManagerAddress:   slashManagerAddress,
		RootChainAddress:      rootChainAddress,
		StakingInfoAddress:    stakingInfoAddress,
		StateSenderAddress:    stateSenderAddress,
		StateReceiverAddress:  stateReceiverAddress,
		ValidatorSetAddress:   validatorSetAddress,
	}
	params := types.NewParams(mainchainTxConfirmations, maticchainTxConfirmations, chainParams)
	chainManagerGenesis := types.NewGenesisState(params)
	fmt.Printf("Selected randomly generated chainmanager parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, chainManagerGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(chainManagerGenesis)
}
