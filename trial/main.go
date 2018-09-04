package main

import (
	"fmt"
	"log"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/xsleonard/go-merkle"
)

func main()  {
	client, err := ethclient.Dial("https://testnet.matic.network")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections

	root:=getHeaders(4733028,4733033,client)
	fmt.Printf("the root hash is %v",root)

}

func getHeaders(start int,end int,client *ethclient.Client) string {
	fmt.Printf(" the start is %v and the end is %v and value is %v",start,end,start<end)
	if start>end{
		return ""
	}
	//TODO fetch block header by making a goroutine , when we get result take sha3 of information and put in array
	current:=start
	var result [][32]byte
	for current <= end {
		//TODO run this in different goroutines and use channels to fetch results(how to maintian order)
		blockheader,err:=client.HeaderByNumber(context.Background(),big.NewInt(int64(current)))
		if err!=nil {
			fmt.Printf("not found")
			log.Fatal(err)
		}
		//fmt.Printf("the block number is %v",blockheader.Number)
		fmt.Println(blockheader.Number)
		fmt.Println(blockheader.Hash().Hex())
		headerBytes:= blockheader.Number.Bytes()
		input,err:= convertTo32(blockheader.Number.Bytes())
		fmt.Printf("blocknumber bytes %v,%v",blockheader.Number,input)
		fmt.Println()
		headerBytes = append(headerBytes,blockheader.Time.Bytes()...)
		headerBytes = append(headerBytes,blockheader.TxHash.Bytes()...)
		headerBytes = append(headerBytes,blockheader.ReceiptHash.Bytes()...)



		header:= sha3.Sum256(headerBytes)
		fmt.Printf("the header is %v",header)


		result = append(result, header)
		current++
	}
	fmt.Printf("loop ended ")
	for _,number := range result{
		// we get 32 bytes headers in a list
		fmt.Println(len(number))
	}
	merkelData:=convert(result)
	fmt.Printf("merkel data is %v",merkelData)
	tree := merkle.NewTree()
	err := tree.Generate(merkelData,sha3.New256())
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Printf("Root: %v\n", tree.Root())// return the hash of root
	fmt.Println(tree.Root().Hash)
	return string(tree.Root().Hash)
}
func convert(input [][32]byte) [][]byte {
	var output [][]byte
	for _,in := range input{
		newInput:=in[:]
		output = append(output, newInput)
		//fmt.Printf("for loop output is %v",output)

	}
	fmt.Printf("the output is %v",output)
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
