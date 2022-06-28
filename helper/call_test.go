package helper

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/maticnetwork/bor/common"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
)

const (
	testTendermintNode = "tcp://localhost:26657"
)

//  Test - to decode signers from checkpoint sigs data
func TestCheckpointsigs(t *testing.T) {
	t.Parallel()

	viper.Set(TendermintNodeFlag, testTendermintNode)
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
	const sigLength = 65

	signersList := make([]string, len(sigInput))

	// Calculate total stake Power of all Signers.
	for i := 0; i < len(sigInput); i += sigLength {
		signature := sigInput[i : i+sigLength]

		pKey, err := authTypes.RecoverPubkey(voteBytes, signature)
		if err != nil {
			fmt.Println("Error Recovering PubKey", "Error", err)
			return nil, err
		}

		signersList[i] = types.NewPubKey(pKey).Address().String()
	}

	return signersList, nil
}
