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

func validateCheckpoint(start int, end int, rootHash string) bool {
	var logger = helper.Logger.With("module", "checkpoint/validate")

	client, err := ethclient.Dial(helper.GetConfig().MaticRPCUrl)
	if err != nil {
		logger.Error("Error Dialing to matic via RPC", "Error", err)
	}

	if (start-end+1)%2 != 0 {
		return false
	}

	root := "0x" + getHeaders(start, end, client)
	if strings.Compare(root, rootHash) == 0 {
		logger.Info("root hash matched ! ")
		return true
	} else {
		logger.Info("root hash does not match ", "Root Hash From Message", rootHash, "Root Hash Generated", root)
		return false
	}
}

func getHeaders(start int, end int, client *ethclient.Client) string {
	var logger = helper.Logger.With("module", "checkpoint/validate")

	if start > end {
		return ""
	}
	//todo add check for even difference

	current := start
	var result [][32]byte
	for current <= end {
		blockheader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(current)))
		if err != nil {
			logger.Error(" Error Getting Block from Matic ", " Error ", err)
		}
		headerBytes := appendBytes32(blockheader.Number.Bytes(),
			blockheader.Time.Bytes(),
			blockheader.TxHash.Bytes(),
			blockheader.ReceiptHash.Bytes())

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
		logger.Error(" Error generating tree ", " Error ", err)
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
		err = fmt.Errorf("input length is greater than 32")
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
