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
	MaticTokenAddress     = "matic_token_address"
	StakingManagerAddress = "staking_manager_address"
	SlashManagerAddress   = "slash_manager_address"
	RootChainAddress      = "root_chain_address"
	StakingInfoAddress    = "staking_info_address"
	StateSenderAddress    = "state_sender_address"

	// Bor Chain Contracts
	StateReceiverAddress = "state_receiver_address"
	ValidatorSetAddress  = "validator_set_address"
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

	var maticchainTxConfirmations uint64
	simState.AppParams.GetOrGenerate(simState.Cdc, MaticchainTxConfirmations, &maticchainTxConfirmations, simState.Rand,
		func(r *rand.Rand) { maticchainTxConfirmations = GenMaticchainTxConfirmations(r) },
	)

	var borChainID string
	simState.AppParams.GetOrGenerate(simState.Cdc, BorChainID, &borChainID, simState.Rand,
		func(r *rand.Rand) { borChainID = GenBorChainId(r) },
	)

	var maticTokenAddress = GenHeimdallAddress()
	var stakingManagerAddress = GenHeimdallAddress()
	var slashManagerAddress = GenHeimdallAddress()
	var rootChainAddress = GenHeimdallAddress()
	var stakingInfoAddress = GenHeimdallAddress()
	var stateSenderAddress = GenHeimdallAddress()
	var stateReceiverAddress = GenHeimdallAddress()
	var validatorSetAddress = GenHeimdallAddress()
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
