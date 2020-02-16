package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/cbergoon/merkletree"
	"github.com/maticnetwork/bor/common/hexutil"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/rpc"
	"github.com/tendermint/crypto/sha3"
	"github.com/xsleonard/go-merkle"
	"golang.org/x/sync/errgroup"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateCheckpoint - Validates if checkpoint rootHash matches or not
func ValidateCheckpoint(start uint64, end uint64, rootHash hmTypes.HeimdallHash) (bool, error) {
	// Check if blocks exist locally
	if !CheckIfBlocksExist(end) {
		return false, errors.New("blocks not found locally")
	}

	// Compare RootHash
	root, err := GetHeaders(start, end)
	if err != nil {
		return false, err
	}

	if bytes.Equal(root, rootHash.Bytes()) {
		return true, nil
	}

	return false, nil
}

// CheckIfBlocksExist - check if latest block number is greater than end block
func CheckIfBlocksExist(end uint64) bool {
	// Get Latest block number.
	rpcClient := helper.GetMaticRPCClient()
	var latestBlock *types.Header

	err := rpcClient.Call(&latestBlock, "eth_getBlockByNumber", "latest", false)
	if err != nil {
		return false
	}

	if end > latestBlock.Number.Uint64() {
		return false
	}

	return true
}

func GetHeaders(start uint64, end uint64) ([]byte, error) {
	rpcClient := helper.GetMaticRPCClient()

	if start > end {
		return nil, errors.New("start is greater than end")
	}

	batchElements := make([]rpc.BatchElem, end-start+1)
	for i := range batchElements {
		param := new(big.Int)
		param.SetUint64(uint64(i) + start)

		batchElements[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{hexutil.EncodeBig(param), true},
			Result: &types.Header{},
		}
	}

	// Batch call
	err := fetchBatchElements(rpcClient, batchElements)
	if err != nil {
		return nil, err
	}

	// Fetch result and draft header and add into tree
	expectedLength := nextPowerOfTwo(end - start + 1)
	headers := make([][32]byte, expectedLength)
	for i, batchElement := range batchElements {
		if batchElement.Error != nil {
			return nil, batchElement.Error
		}

		blockHeader := batchElement.Result.(*types.Header)
		header := crypto.Keccak256(appendBytes32(
			blockHeader.Number.Bytes(),
			new(big.Int).SetUint64(blockHeader.Time).Bytes(),
			blockHeader.TxHash.Bytes(),
			blockHeader.ReceiptHash.Bytes(),
		))

		var arr [32]byte
		copy(arr[:], header)

		// set header
		headers[i] = arr
	}

	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})
	if err := tree.Generate(convert(headers), sha3.NewLegacyKeccak256()); err != nil {
		return nil, err
	}

	return tree.Root().Hash, nil
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
	dividendAccounts = hmTypes.SortDividendAccountByID(dividendAccounts)
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
func GetAccountProof(dividendAccounts []hmTypes.DividendAccount, dividendAccountID hmTypes.DividendAccountID) ([]byte, error) {
	// Sort the dividendAccounts by ID
	dividendAccounts = hmTypes.SortDividendAccountByID(dividendAccounts)
	var list []merkletree.Content
	var account hmTypes.DividendAccount

	for i := 0; i < len(dividendAccounts); i++ {
		list = append(list, dividendAccounts[i])
		if dividendAccounts[i].ID == dividendAccountID {
			account = dividendAccounts[i]
		}
	}

	tree, err := merkletree.NewTreeWithHashStrategy(list, sha3.NewLegacyKeccak256)
	if err != nil {
		return nil, err
	}

	branchArray, _, err := tree.GetMerklePath(account)

	// concatenate branch array
	proof := appendBytes32(branchArray...)
	return proof, err
}

// VerifyAccountProof returns proof of dividend Account
func VerifyAccountProof(dividendAccounts []hmTypes.DividendAccount, dividendAccountID hmTypes.DividendAccountID, proofToVerify string) (bool, error) {

	proof, err := GetAccountProof(dividendAccounts, dividendAccountID)
	if err != nil {
		return false, nil
	}

	if proofToVerify == hex.EncodeToString(proof) {
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
func fetchBatchElements(rpcClient *rpc.Client, elements []rpc.BatchElem) (err error) {
	var batchLength = int(helper.GetConfig().AvgCheckpointLength)
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
