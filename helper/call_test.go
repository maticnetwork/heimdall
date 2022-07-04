package helper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/maticnetwork/bor/common"
	ethTypes "github.com/maticnetwork/bor/core/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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

// TestDecodeValidatorStakeUpdateEvent
func TestDecodeValidatorStakeUpdateEvent(t *testing.T) {
	contractCallerObj, err := NewContractCaller()
	if err != nil {
		fmt.Println("Error creating contract caller")
	}
	testContractAddress := common.HexToAddress("0x29c40836c17f22d16a7fe953fb25da670c96d69e")
	testTxReceipt := &ethTypes.Receipt{
		BlockNumber: big.NewInt(1),
		Logs: []*ethTypes.Log{
			{
				Address: common.HexToAddress("0x29c40836c17f22d16a7fe953fb25da670c96d69e"),
				Topics: []common.Hash{
					common.HexToHash("0x31d1715032654fde9867c0f095aecce1113049e30b9f4ecbaa6954ed6c63b8df"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
					common.HexToHash("0x000000000000000000000000e29d3d4d72997b31ccdf8188113c189f1106f6b8"),
					common.HexToHash("0x000000000000000000000000000000000000000000000001c1e7de9a29a2ae50"),
				},
				Index: 10,
			},
		},
	}
	event, err := contractCallerObj.DecodeValidatorStakeUpdateEvent(testContractAddress, testTxReceipt, 10)
	assert.Nil(t, event)
	assert.Error(t, err)

	testTxReceipt = &ethTypes.Receipt{
		BlockNumber: big.NewInt(1),
		Logs: []*ethTypes.Log{
			{
				Address: common.HexToAddress("0x29c40836c17f22d16a7fe953fb25da670c96d69e"),
				Topics: []common.Hash{
					common.HexToHash("0x35af9eea1f0e7b300b0a14fae90139a072470e44daa3f14b5069bebbc1265bda"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001"),
					common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000026b09"),
					common.HexToHash("0x0000000000000000000000000000000000000000002ea05349e479938e4ed6e6"),
				},
				Index: 20,
			},
		},
	}
	event, err = contractCallerObj.DecodeValidatorStakeUpdateEvent(testContractAddress, testTxReceipt, 20)
	assert.NotNil(t, event)
	assert.NoError(t, err)
}
