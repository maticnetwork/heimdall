package helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"

	hmTypes "github.com/maticnetwork/heimdall/types"

	cfg "github.com/tendermint/tendermint/config"
)

// Test - to check heimdall config
func TestHeimdallConfig(t *testing.T) {
	t.Parallel()

	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set(TendermintNodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	// cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastSync
	// cliCtx.TrustNode = true

	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	fmt.Println("Address", GetAddress())

	pubKey := GetPubKey()

	fmt.Println("PublicKey", pubKey.String())
}

func TestHeimdallConfigNewSelectionAlgoHeight(t *testing.T) {
	t.Parallel()

	data := map[string]bool{"mumbai": false, "mainnet": false, "local": true}
	for chain, shouldBeZero := range data {
		conf.BorRPCUrl = "" // allow config to be loaded again

		viper.Set("chain", chain)

		InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

		nsah := GetNewSelectionAlgoHeight()
		if nsah == 0 && !shouldBeZero || nsah != 0 && shouldBeZero {
			t.Errorf("Invalid GetNewSelectionAlgoHeight = %d for chain %s", nsah, chain)
		}
	}
}

func TestGetChainManagerAddressMigration(t *testing.T) {
	t.Parallel()

	newMaticContractAddress := "0x0000000000000000000000000000000000001234"

	chainManagerAddressMigrations["mumbai"] = map[int64]ChainManagerAddressMigration{
		350: {MaticTokenAddress: hmTypes.HexToHeimdallAddress(newMaticContractAddress)},
	}

	viper.Set("chain", "mumbai")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	migration, found := GetChainManagerAddressMigration(350)

	if !found {
		t.Errorf("Expected migration to be found")
	}

	if migration.MaticTokenAddress.String() != newMaticContractAddress {
		t.Errorf("Expected matic token address to be %s, got %s", newMaticContractAddress, migration.MaticTokenAddress.String())
	}

	// test for non existing migration
	_, found = GetChainManagerAddressMigration(351)
	if found {
		t.Errorf("Expected migration to not be found")
	}

	// test for non existing chain
	conf.BorRPCUrl = ""

	viper.Set("chain", "newChain")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	_, found = GetChainManagerAddressMigration(350)
	if found {
		t.Errorf("Expected migration to not be found")
	}
}

func TestHeimdallConfigUpdateTendermintConfig(t *testing.T) {
	t.Parallel()

	type teststruct struct {
		chain string
		viper string
		def   string
		value string
	}

	data := []teststruct{
		{chain: "mumbai", viper: "viper", def: "default", value: "viper"},
		{chain: "mumbai", viper: "viper", def: "", value: "viper"},
		{chain: "mumbai", viper: "", def: "default", value: "default"},
		{chain: "mumbai", viper: "", def: "", value: DefaultMumbaiTestnetSeeds},
		{chain: "amoy", viper: "", def: "", value: DefaultAmoyTestnetSeeds},
		{chain: "mainnet", viper: "viper", def: "default", value: "viper"},
		{chain: "mainnet", viper: "viper", def: "", value: "viper"},
		{chain: "mainnet", viper: "", def: "default", value: "default"},
		{chain: "mainnet", viper: "", def: "", value: DefaultMainnetSeeds},
		{chain: "local", viper: "viper", def: "default", value: "viper"},
		{chain: "local", viper: "viper", def: "", value: "viper"},
		{chain: "local", viper: "", def: "default", value: "default"},
		{chain: "local", viper: "", def: "", value: ""},
	}

	oldConf := conf.Chain
	viperObj := viper.New()
	tendermintConfig := cfg.DefaultConfig()

	for _, ts := range data {
		conf.Chain = ts.chain
		tendermintConfig.P2P.Seeds = ts.def
		viperObj.Set(SeedsFlag, ts.viper)
		UpdateTendermintConfig(tendermintConfig, viperObj)

		if tendermintConfig.P2P.Seeds != ts.value {
			t.Errorf("Invalid UpdateTendermintConfig, tendermintConfig.P2P.Seeds not set correctly")
		}
	}

	conf.Chain = oldConf
}
