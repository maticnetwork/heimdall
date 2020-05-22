package types

import (
	"bytes"
	"errors"

	"github.com/cbergoon/merkletree"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/rpc"
	"github.com/tendermint/crypto/sha3"
	"golang.org/x/sync/errgroup"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateCheckpoint - Validates if checkpoint rootHash matches or not
func ValidateCheckpoint(start uint64, end uint64, rootHash hmTypes.HeimdallHash, checkpointLength uint64, contractCaller helper.IContractCaller) (bool, error) {
	// Check if blocks exist locally
	if !contractCaller.CheckIfBlocksExist(end) {
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
	var list []merkletree.Content

	for i := 0; i < len(dividendAccounts); i++ {
		list = append(list, dividendAccounts[i])
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
	var list []merkletree.Content
	var account hmTypes.DividendAccount
	index := uint64(0)
	for i := 0; i < len(dividendAccounts); i++ {
		list = append(list, dividendAccounts[i])
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
		return false, nil
	}

	// check proof bytes
	if bytes.Equal(common.FromHex(proofToVerify), proof) {
		return true, nil
	}

	return false, nil
}

func convert(input []([32]byte)) [][]byte {
	var output [][]byte
	for _, in := range input {
		newInput := make([]byte, len(in[:]))
		copy(newInput, in[:])
		output = append(output, newInput)

	}
	return output
}

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

func nextPowerOfTwo(n uint64) uint64 {
	if n == 0 {
		return 1
	}
	// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n |= n >> 32
	n++
	return n
}

// spins go-routines to fetch batch elements to allow creation of large merkle trees
func fetchBatchElements(rpcClient *rpc.Client, elements []rpc.BatchElem, checkpointLength uint64) (err error) {
	var batchLength = int(checkpointLength)
	// group
	var g errgroup.Group

	for i := 0; i < len(elements); i += batchLength {
		var newBatch []rpc.BatchElem
		if len(elements) < i+batchLength {
			newBatch = elements[i:]
		} else {
			newBatch = elements[i : i+batchLength]
		}

		// common.CheckpointLogger.Info("Batching requests", "index", i, "length", len(newBatch))

		// spawn go-routine
		g.Go(func() error {
			// Batch call
			err := rpcClient.BatchCall(newBatch)
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	// common.CheckpointLogger.Info("Fetched all headers", "len", len(elements))
	return nil
}
