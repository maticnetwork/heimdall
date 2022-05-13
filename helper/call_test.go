package helper

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/maticnetwork/bor/common"
	"github.com/spf13/viper"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
)

//  Test - to decode signers from checkpoint sigs data
func TestCheckpointsigs(t *testing.T) {
	tendermintNode := "tcp://localhost:26657"
	viper.Set(NodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	contractCallerObj, err := NewContractCaller()
	if err != nil {
		fmt.Println("Error creating contract caller")
	}

	txHashStr := "0x9c2a9e20e1fecdae538f72b01dd0fd5008cc90176fd603b92b59274d754cbbd8"
	txHash := common.HexToHash(txHashStr)
	voteSignBytes, sigs, txData, err := contractCallerObj.GetCheckpointSign(txHash)
	if err != nil {
		fmt.Println("Error fetching checkpoint tx input args")
	}
	fmt.Println("checkpoint args", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))

	signerList, err := FetchSigners(voteSignBytes, sigs)
	if err != nil {
		fmt.Println("Error fetching signer list from tx input args")
	}
	fmt.Println("signers list", signerList)
}

// CalculateSignerRewards calculates new rewards for signers
func FetchSigners(voteBytes []byte, sigInput []byte) ([]string, error) {
	signersList := []string{}
	const sigLength = 65

	// Calculate total stake Power of all Signers.
	for i := 0; i < len(sigInput); i += sigLength {
		signature := sigInput[i : i+sigLength]
		pKey, err := authTypes.RecoverPubkey(voteBytes, []byte(signature))
		if err != nil {
			fmt.Println("Error Recovering PubKey", "Error", err)
			return nil, err
		}

		pubKey := types.NewPubKey(pKey)
		signerAddress := pubKey.Address().String()
		signersList = append(signersList, signerAddress)
	}
	return signersList, nil
}
