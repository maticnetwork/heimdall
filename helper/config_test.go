package helper

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

//  Test - to check heimdall config
func TestHeimdallConfig(t *testing.T) {
	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set(NodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	// cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastSync
	// cliCtx.TrustNode = true

	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	fmt.Println("Address", GetAddress())
	pubKey := GetPubKey()
	fmt.Println("PublicKey", pubKey.String())
	// fmt.Println("CryptoPublicKey", pubKey.CryptoPubKey().String())
}

func TestHeimdallConfigNewSelectionAlgoHeight(t *testing.T) {
	var data map[string]bool = map[string]bool{"mumbai": false, "mainnet": false, "local": true}
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
