package chainmanager

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/x/chainmanager/keeper"
	"github.com/maticnetwork/heimdall/x/chainmanager/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if genState.Params.ChainParams.MaticTokenAddress != "" {
		genState.Params.ChainParams.MaticTokenAddress = strings.ToLower(genState.Params.ChainParams.MaticTokenAddress)
	}
	if genState.Params.ChainParams.StakingManagerAddress != "" {
		genState.Params.ChainParams.StakingManagerAddress = strings.ToLower(genState.Params.ChainParams.StakingManagerAddress)
	}
	if genState.Params.ChainParams.SlashManagerAddress != "" {
		genState.Params.ChainParams.SlashManagerAddress = strings.ToLower(genState.Params.ChainParams.SlashManagerAddress)
	}
	if genState.Params.ChainParams.RootChainAddress != "" {
		genState.Params.ChainParams.RootChainAddress = strings.ToLower(genState.Params.ChainParams.RootChainAddress)
	}
	if genState.Params.ChainParams.StakingInfoAddress != "" {
		genState.Params.ChainParams.StakingInfoAddress = strings.ToLower(genState.Params.ChainParams.StakingInfoAddress)
	}
	if genState.Params.ChainParams.StateSenderAddress != "" {
		genState.Params.ChainParams.StateSenderAddress = strings.ToLower(genState.Params.ChainParams.StateSenderAddress)
	}
	if genState.Params.ChainParams.StateReceiverAddress != "" {
		genState.Params.ChainParams.StateReceiverAddress = strings.ToLower(genState.Params.ChainParams.StateReceiverAddress)
	}
	if genState.Params.ChainParams.ValidatorSetAddress != "" {
		genState.Params.ChainParams.ValidatorSetAddress = strings.ToLower(genState.Params.ChainParams.ValidatorSetAddress)
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
