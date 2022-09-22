package helper

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/contracts/erc20"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/slashmanager"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statereceiver"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/types"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	testTendermintNode = "tcp://localhost:26657"
)

// TestCheckpointSigs decodes signers from checkpoint sigs data
func TestCheckpointSigs(t *testing.T) {
	t.Parallel()

	viper.Set(TendermintNodeFlag, testTendermintNode)
	viper.Set("log_level", "info")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	contractCallerObj, err := NewContractCaller()
	if err != nil {
		t.Error("Error creating contract caller")
	}

	txHashStr := "0x9c2a9e20e1fecdae538f72b01dd0fd5008cc90176fd603b92b59274d754cbbd8"
	txHash := common.HexToHash(txHashStr)

	voteSignBytes, sigs, txData, err := contractCallerObj.GetCheckpointSign(txHash)
	if err != nil {
		t.Error("Error fetching checkpoint tx input args")
	}

	t.Log("checkpoint args", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))

	signerList, err := FetchSigners(voteSignBytes, sigs)
	if err != nil {
		t.Error("Error fetching signer list from tx input args")
	}

	t.Log("signers list", signerList)
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
			return nil, err
		}

		signersList[i] = types.NewPubKey(pKey).Address().String()
	}

	return signersList, nil
}

// TestPopulateABIs tests that package level ABIs cache works as expected
// by not invoking json methods after contracts ABIs' init
func TestPopulateABIs(t *testing.T) {
	t.Parallel()

	viper.Set(TendermintNodeFlag, testTendermintNode)
	viper.Set("log_level", "info")
	InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	t.Log("ABIs map should be empty and all ABIs not found")
	assert.True(t, len(ContractsABIsMap) == 0)
	_, found := ContractsABIsMap[rootchain.RootchainABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[stakinginfo.StakinginfoABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[statereceiver.StatereceiverABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[statesender.StatesenderABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[stakemanager.StakemanagerABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[slashmanager.SlashmanagerABI]
	assert.False(t, found)
	_, found = ContractsABIsMap[erc20.Erc20ABI]
	assert.False(t, found)

	t.Log("Should create a new contract caller and populate its ABIs by decoding json")

	contractCallerObjFirst, err := NewContractCaller()
	if err != nil {
		t.Error("Error creating contract caller")
	}

	assert.Equalf(t, ContractsABIsMap[rootchain.RootchainABI], &contractCallerObjFirst.RootChainABI,
		"values for %s not equals", rootchain.RootchainABI)
	assert.Equalf(t, ContractsABIsMap[stakinginfo.StakinginfoABI], &contractCallerObjFirst.StakingInfoABI,
		"values for %s not equals", stakinginfo.StakinginfoABI)
	assert.Equalf(t, ContractsABIsMap[statereceiver.StatereceiverABI], &contractCallerObjFirst.StateReceiverABI,
		"values for %s not equals", statereceiver.StatereceiverABI)
	assert.Equalf(t, ContractsABIsMap[statesender.StatesenderABI], &contractCallerObjFirst.StateSenderABI,
		"values for %s not equals", statesender.StatesenderABI)
	assert.Equalf(t, ContractsABIsMap[stakemanager.StakemanagerABI], &contractCallerObjFirst.StakeManagerABI,
		"values for %s not equals", stakemanager.StakemanagerABI)
	assert.Equalf(t, ContractsABIsMap[slashmanager.SlashmanagerABI], &contractCallerObjFirst.SlashManagerABI,
		"values for %s not equals", slashmanager.SlashmanagerABI)
	assert.Equalf(t, ContractsABIsMap[erc20.Erc20ABI], &contractCallerObjFirst.MaticTokenABI,
		"values for %s not equals", erc20.Erc20ABI)

	t.Log("ABIs map should not be empty and all ABIs found")
	assert.True(t, len(ContractsABIsMap) == 8)
	_, found = ContractsABIsMap[rootchain.RootchainABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[stakinginfo.StakinginfoABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[statereceiver.StatereceiverABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[statesender.StatesenderABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[stakemanager.StakemanagerABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[slashmanager.SlashmanagerABI]
	assert.True(t, found)
	_, found = ContractsABIsMap[erc20.Erc20ABI]
	assert.True(t, found)

	t.Log("Should create a new contract caller and populate its ABIs by using cached map")

	contractCallerObjSecond, err := NewContractCaller()
	if err != nil {
		t.Log("Error creating contract caller")
	}

	assert.Equalf(t, ContractsABIsMap[rootchain.RootchainABI], &contractCallerObjSecond.RootChainABI,
		"values for %s not equals", rootchain.RootchainABI)
	assert.Equalf(t, ContractsABIsMap[stakinginfo.StakinginfoABI], &contractCallerObjSecond.StakingInfoABI,
		"values for %s not equals", stakinginfo.StakinginfoABI)
	assert.Equalf(t, ContractsABIsMap[statereceiver.StatereceiverABI], &contractCallerObjSecond.StateReceiverABI,
		"values for %s not equals", statereceiver.StatereceiverABI)
	assert.Equalf(t, ContractsABIsMap[statesender.StatesenderABI], &contractCallerObjSecond.StateSenderABI,
		"values for %s not equals", statesender.StatesenderABI)
	assert.Equalf(t, ContractsABIsMap[stakemanager.StakemanagerABI], &contractCallerObjSecond.StakeManagerABI,
		"values for %s not equals", stakemanager.StakemanagerABI)
	assert.Equalf(t, ContractsABIsMap[slashmanager.SlashmanagerABI], &contractCallerObjSecond.SlashManagerABI,
		"values for %s not equals", slashmanager.SlashmanagerABI)
	assert.Equalf(t, ContractsABIsMap[erc20.Erc20ABI], &contractCallerObjSecond.MaticTokenABI,
		"values for %s not equals", erc20.Erc20ABI)
}
