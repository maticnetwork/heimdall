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
	"github.com/stretchr/testify/assert"
)

const (
	testTendermintNode = "tcp://localhost:26657"
)

//  TestCheckpointSigs decodes signers from checkpoint sigs data
func TestCheckpointSigs(t *testing.T) {
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

// FetchSigners fetches the signers' list
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

//  TestPopulateABIs tests that package level ABIs cache works as expected
//  by not invoking json methods after contracts ABIs' init
func TestPopulateABIs(t *testing.T) {
	t.Parallel()

	viper.Set(TendermintNodeFlag, testTendermintNode)
	viper.Set("log_level", "info")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	assert.True(t, len(ContractsABIsMap) == 0)

	fmt.Println("Should create a new contract caller and populate its ABIs by decoding json")

	contractCallerObjFirst, err := NewContractCaller()
	if err != nil {
		fmt.Println("Error creating contract caller")
	}

	assert.Equalf(t, ContractsABIsMap[RootChainABI], &contractCallerObjFirst.RootChainABI,
		"values for %s not equals", RootChainABI)
	assert.Equalf(t, ContractsABIsMap[StakingInfoABI], &contractCallerObjFirst.StakingInfoABI,
		"values for %s not equals", StakingInfoABI)
	assert.Equalf(t, ContractsABIsMap[StateReceiverABI], &contractCallerObjFirst.StateReceiverABI,
		"values for %s not equals", StateReceiverABI)
	assert.Equalf(t, ContractsABIsMap[StateSenderABI], &contractCallerObjFirst.StateSenderABI,
		"values for %s not equals", StateSenderABI)
	assert.Equalf(t, ContractsABIsMap[StakeManagerABI], &contractCallerObjFirst.StakeManagerABI,
		"values for %s not equals", StakeManagerABI)
	assert.Equalf(t, ContractsABIsMap[SlashManagerABI], &contractCallerObjFirst.SlashManagerABI,
		"values for %s not equals", SlashManagerABI)
	assert.Equalf(t, ContractsABIsMap[MaticTokenABI], &contractCallerObjFirst.MaticTokenABI,
		"values for %s not equals", MaticTokenABI)

	fmt.Println("Should create a new contract caller and populate its ABIs by using cached map")

	contractCallerObjSecond, err := NewContractCaller()
	if err != nil {
		fmt.Println("Error creating contract caller")
	}

	assert.Emptyf(t, &contractCallerObjSecond.RootChainABI, "contract caller %s not empty", RootChainABI)
	assert.Emptyf(t, &contractCallerObjSecond.StakingInfoABI, "contract caller %s not empty", StakingInfoABI)
	assert.Emptyf(t, &contractCallerObjSecond.StateReceiverABI, "contract caller %s not empty", StateReceiverABI)
	assert.Emptyf(t, &contractCallerObjSecond.StateSenderABI, "contract caller %s not empty", StateSenderABI)
	assert.Emptyf(t, &contractCallerObjSecond.StakeManagerABI, "contract caller %s not empty", StakeManagerABI)
	assert.Emptyf(t, &contractCallerObjSecond.SlashManagerABI, "contract caller %s not empty", SlashManagerABI)
	assert.Emptyf(t, &contractCallerObjSecond.MaticTokenABI, "contract caller %s not empty", MaticTokenABI)
}
