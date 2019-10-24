package checkpoint

import (
	"bytes"
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/xsleonard/go-merkle"
	"golang.org/x/crypto/sha3"
	"golang.org/x/sync/errgroup"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func ValidateCheckpoint(start uint64, end uint64, rootHash hmTypes.HeimdallHash) bool {
	root, err := GetHeaders(start, end)
	if err != nil {
		return false
	}

	if bytes.Equal(root, rootHash.Bytes()) {
		return true
	}

	return false
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
func GetAccountRootHash(validatorAccounts []hmTypes.ValidatorAccount) ([]byte, error) {
	// Sort the validatorAccounts by valID
	sort.Slice(validatorAccounts, func(i, j int) bool { return validatorAccounts[i].ID < validatorAccounts[j].ID })

	expectedLength := len(validatorAccounts)
	valAccountHashes := make([][32]byte, expectedLength)
	i := 0
	for _, va := range validatorAccounts {
		valAccountHash := crypto.Keccak256(appendBytes32(
			new(big.Int).SetUint64(uint64(va.ID)).Bytes(),
			va.RewardAmount,
			va.SlashedAmount,
		))
		var arr [32]byte
		copy(arr[:], valAccountHash)

		valAccountHashes[i] = arr
		i++
	}

	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})
	if err := tree.Generate(convert(valAccountHashes), sha3.NewLegacyKeccak256()); err != nil {
		return nil, err
	}
	return tree.Root().Hash, nil
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
