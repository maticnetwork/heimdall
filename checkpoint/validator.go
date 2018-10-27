package checkpoint

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/xsleonard/go-merkle"
)

func ValidateCheckpoint(start uint64, end uint64, rootHash string) bool {
	client := helper.GetMaticClient()

	if (start-end+1)%2 != 0 {
		return false
	}

	root := "0x" + GetHeaders(start, end, client)
	if strings.Compare(root, rootHash) == 0 {
		CheckpointLogger.Info("RootHash matched!")
		return true
	} else {
		CheckpointLogger.Error("RootHash does not match", "rootHashTx", rootHash, "rootHash", root)
		return false
	}
}

func GetHeaders(start uint64, end uint64, client *ethclient.Client) string {
	if start > end {
		return ""
	}

	// TODO add check for even difference
	current := start
	var result [][32]byte
	for current <= end {
		blockheader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(current)))
		if err != nil {
			CheckpointLogger.Error("Error getting block from Matic", "error", err, "start", start, "end", end, "current", current)
			return ""
		}

		headerBytes := appendBytes32(
			blockheader.Number.Bytes(),
			blockheader.Time.Bytes(),
			blockheader.TxHash.Bytes(),
			blockheader.ReceiptHash.Bytes(),
		)

		header := getsha3frombyte(headerBytes)
		var arr [32]byte
		copy(arr[:], header)
		result = append(result, arr)
		current++
	}
	merkelData := convert(result)
	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})

	err := tree.Generate(merkelData, sha3.NewKeccak256())
	if err != nil {
		CheckpointLogger.Error("Error generating tree", "error", err)
		return ""
	}

	return hex.EncodeToString(tree.Root().Hash)
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
		err = fmt.Errorf("Input length is greater than 32")
		CheckpointLogger.Error("Input length is greater than 32 while converting", "error", err)
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

func getsha3frombyte(input []byte) []byte {
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(input)
	buf = hash.Sum(buf)
	return buf
}
