package types

import (
	"bytes"
	"errors"

	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/common"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/crypto/sha3"
)

// ValidateCheckpoint - Validates if checkpoint rootHash matches or not
func ValidateCheckpoint(start uint64, end uint64, rootHash hmTypes.HeimdallHash, checkpointLength uint64, contractCaller helper.IContractCaller, confirmations uint64) (bool, error) {
	// Check if blocks exist locally
	if !contractCaller.CheckIfBlocksExist(end + confirmations) {
		return false, errors.New("blocks not found locally")
	}

	// Compare RootHash
	root, err := contractCaller.GetRootHash(start, end, checkpointLength)
	if err != nil {
		return false, err
	}

	if bytes.Equal(root, rootHash.Bytes()) {
		return true, nil
	}

	return false, nil
}

// GetAccountRootHash returns roothash of Validator Account State Tree
func GetAccountRootHash(dividendAccounts []hmTypes.DividendAccount) ([]byte, error) {
	tree, err := GetAccountTree(dividendAccounts)
	if err != nil {
		return nil, err
	}

	return tree.Root.Hash, nil
}

// GetAccountTree returns roothash of Validator Account State Tree
func GetAccountTree(dividendAccounts []hmTypes.DividendAccount) (*merkletree.MerkleTree, error) {
	// Sort the dividendAccounts by ID
	dividendAccounts = hmTypes.SortDividendAccountByAddress(dividendAccounts)
	list := make([]merkletree.Content, len(dividendAccounts))

	for i := 0; i < len(dividendAccounts); i++ {
		list[i] = dividendAccounts[i]
	}

	tree, err := merkletree.NewTreeWithHashStrategy(list, sha3.NewLegacyKeccak256)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// GetAccountProof returns proof of dividend Account
func GetAccountProof(dividendAccounts []hmTypes.DividendAccount, userAddr hmTypes.HeimdallAddress) ([]byte, uint64, error) {
	// Sort the dividendAccounts by user address
	dividendAccounts = hmTypes.SortDividendAccountByAddress(dividendAccounts)

	var (
		list    = make([]merkletree.Content, len(dividendAccounts))
		account hmTypes.DividendAccount
	)

	index := uint64(0)

	for i := 0; i < len(dividendAccounts); i++ {
		list[i] = dividendAccounts[i]

		if dividendAccounts[i].User.Equals(userAddr) {
			account = dividendAccounts[i]
			index = uint64(i)
		}
	}

	tree, err := merkletree.NewTreeWithHashStrategy(list, sha3.NewLegacyKeccak256)
	if err != nil {
		return nil, 0, err
	}

	branchArray, _, err := tree.GetMerklePath(account)

	// concatenate branch array
	proof := appendBytes32(branchArray...)

	return proof, index, err
}

// VerifyAccountProof returns proof of dividend Account
func VerifyAccountProof(dividendAccounts []hmTypes.DividendAccount, userAddr hmTypes.HeimdallAddress, proofToVerify string) (bool, error) {
	proof, _, err := GetAccountProof(dividendAccounts, userAddr)
	if err != nil {
		// nolint: nilerr
		return false, nil
	}

	// check proof bytes
	if bytes.Equal(common.FromHex(proofToVerify), proof) {
		return true, nil
	}

	return false, nil
}

//nolint:unparam
func convertTo32(input []byte) (output [32]byte, err error) {
	l := len(input)
	if l > 32 || l == 0 {
		return
	}

	copy(output[32-l:], input[:])

	return
}

func appendBytes32(data ...[]byte) []byte {
	var result []byte

	for _, v := range data {
		paddedV, err := convertTo32(v)
		if err == nil {
			result = append(result, paddedV[:]...)
		}
	}

	return result
}
