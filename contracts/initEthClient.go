package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/basecoin/contracts/StakeManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/xsleonard/go-merkle"
	"log"
	"math/big"
)

var (
	stakeManagerAddress = "0x8b28d78eb59c323867c43b4ab8d06e0f1efa1573"
)

func getValidatorByIndex(_index int64) abci.Validator {
	client := initKovan()
	stakeManagerInstance, err := StakeManager.NewContracts(common.HexToAddress(stakeManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	validator, _ := stakeManagerInstance.Validators(nil, big.NewInt(_index))
	var _pubkey secp256k1.PubKeySecp256k1
	_pub, _ := hex.DecodeString(validator.Pubkey)
	copy(_pubkey[:], _pub[:])
	_address, _ := hex.DecodeString(_pubkey.Address().String())

	abciValidator := abci.Validator{
		Address: _address,
		Power:   validator.Power.Int64(),
		PubKey:  tmtypes.TM2PB.PubKey(_pubkey),
	}
	return abciValidator

}

func getLastValidator() int64 {
	client := initKovan()
	stakeManagerInstance, err := StakeManager.NewContracts(common.HexToAddress(stakeManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	last, _ := stakeManagerInstance.LastValidatorIndex(nil)
	return last.Int64()
}

func sendCheckpoint() {
	//clientKovan := initKovan()
	clientMatic := initMatic()
	rootHash := getHeaders(4733028, 4733031, clientMatic)
	fmt.Printf("Root hash obtained for blocks from %v to %v is %v", 4733028, 4733031, rootHash)

}

func initKovan() *ethclient.Client {
	client, err := ethclient.Dial("https://kovan.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	return client
}
func initMatic() *ethclient.Client {
	client, err := ethclient.Dial("https://testnet.matic.network")
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getHeaders(start int, end int, client *ethclient.Client) string {
	if (start-end+1)%2 != 0 {
		return "Not Defined , make sure difference is even "
	}
	if start > end {
		return ""
	}
	current := start
	var result [][32]byte
	for current <= end {
		//TODO run this in different goroutines and use channels to fetch results(how to maintian order)
		blockheader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(current)))
		if err != nil {
			fmt.Printf("not found")
			log.Fatal(err)
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
		fmt.Println("*********ERROR***********")
		log.Fatal(err)
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
