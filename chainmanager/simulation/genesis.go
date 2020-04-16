package simulation

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/chainmanager/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// Parameter keys
const (
	TxConfirmationTime    = "tx_confirmation_time"
	BorChainID            = "bor_chain_id"
	MaticTokenAddress     = "matic_token_address"
	StakingManagerAddress = "staking_manager_address"
	RootChainAddress      = "root_chain_address"
	StakingInfoAddress    = "staking_info_address"
	StateSenderAddress    = "state_sender_address"

	// Bor Chain Contracts
	StateReceiverAddress = "state_receiver_address"
	ValidatorSetAddress  = "validator_set_address"
)

func GenTxConfirmationTime() time.Duration {
	// create seed
	seed := rand.New(rand.NewSource(int64(rand.Int())))
	return time.Duration(simulation.RandTimestamp(seed).Unix())
}

func GenHeimdallAddress() hmTypes.HeimdallAddress {
	return hmTypes.BytesToHeimdallAddress(simulation.RandHex(20))
}

// GenBorChainId returns randomc chainID
func GenBorChainId(r *rand.Rand) string {
	return strconv.Itoa(simulation.RandIntBetween(r, 0, 4294967295))
}

func RandomizedGenState(simState *module.SimulationState) {
	var txConfirmationTime time.Duration
	simState.AppParams.GetOrGenerate(simState.Cdc, TxConfirmationTime, &txConfirmationTime, simState.Rand,
		func(r *rand.Rand) { txConfirmationTime = GenTxConfirmationTime() },
	)

	var borChainID string
	simState.AppParams.GetOrGenerate(simState.Cdc, BorChainID, &borChainID, simState.Rand,
		func(r *rand.Rand) { borChainID = GenBorChainId(r) },
	)

	var maticTokenAddress = GenHeimdallAddress()
	var stakingManagerAddress = GenHeimdallAddress()
	var rootChainAddress = GenHeimdallAddress()
	var stakingInfoAddress = GenHeimdallAddress()
	var stateSenderAddress = GenHeimdallAddress()
	var stateReceiverAddress = GenHeimdallAddress()
	var validatorSetAddress = GenHeimdallAddress()
	chainParams := types.ChainParams{
		BorChainID:            borChainID,
		MaticTokenAddress:     maticTokenAddress,
		StakingManagerAddress: stakingManagerAddress,
		RootChainAddress:      rootChainAddress,
		StakingInfoAddress:    stakingInfoAddress,
		StateSenderAddress:    stateSenderAddress,
		StateReceiverAddress:  stateReceiverAddress,
		ValidatorSetAddress:   validatorSetAddress,
	}
	params := types.NewParams(txConfirmationTime, chainParams)
	chainManagerGenesis := types.NewGenesisState(params)
	fmt.Printf("Selected randomly generated chainmanager parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, chainManagerGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(chainManagerGenesis)
}
