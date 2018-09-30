package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"

	"github.com/xsleonard/go-merkle"
)

func main() {
	client, err := ethclient.Dial("https://testnet.matic.network")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections
	//TODO reject transaction if the difference in numbers is not even
	root := getHeaders(4733028, 4733032, client)
	fmt.Printf("the root hash is %v", root)

}

func getHeaders(start int, end int, client *ethclient.Client) string {
	if start > end {
		return ""
	}
	fmt.Printf("start from %v /n end to %v", start, end)
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
	fmt.Println("------")
	merkelData := convert(result)
	//fmt.Println("------")
	fmt.Printf("merkel data is \n %v", len(merkelData))
	fmt.Println("------")
	//tree := merkle.NewTree()
	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})

	err := tree.Generate(merkelData, sha3.NewKeccak256())
	if err != nil {
		fmt.Println("*********ERROR***********")
		log.Fatal(err)
	}
	fmt.Println("------")
	//fmt.Printf("tree: %v\n", tree.Leaves())// return the hash of root
	fmt.Println("------")
	fmt.Println(hex.EncodeToString(tree.Root().Hash))
	return hex.EncodeToString(tree.Root().Hash)
}
func convert(input []([32]byte)) [][]byte {
	var output [][]byte
	for _, in := range input {
		newInput := make([]byte, len(in[:]))
		copy(newInput, in[:])
		output = append(output, newInput)
		//fmt.Printf("------- \n input is %v \n output is %v",newInput,output )

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
	//v,err:=hex.DecodeString(input)
	//if err!=nil{
	//	fmt.Println("Error occured in getsha3")
	//	log.Fatal(err)
	//}
	var buf []byte
	hash.Write(input)
	buf = hash.Sum(buf)
	return buf
}
