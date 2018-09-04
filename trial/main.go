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

	getHeaders(4733028,4733029,client)
	//blockheader,err:=client.HeaderByNumber(context.Background(),big.NewInt(4345947))
	//if err!=nil {
	//	fmt.Printf("not found")
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("the header for the  blocks is %v", blockheader.Hash().Hex())




	//block,err := client.BlockByHash(context.Background(),common.HexToHash("0x330c180b100187e9e61d20a8a4af351103eb6c295366e09d379acea90523e4b8"))
	//if err!=nil {
	//	fmt.Printf("not found")
	//	log.Fatal(err)
	//}
	//fmt.Println(block.ReceiptHash().Hex())
	//fmt.Println(block.Hash().Hex())
	//fmt.Println(block.TxHash().Hex())
}
//func getHeaderBlockSha3(blockNumber big.Int) bytes.Buffer {
//
//}

func getHeaders(start int,end int,client *ethclient.Client)  {
	//if start<end{
	//	return ""
	//}
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
	//err := tree.Generate(result,sha3.New256())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	fmt.Printf("Root: %v\n", tree.Root())
	//return "lol"

}
func convert(input [][32]byte) [][]byte {
	var output [][]byte
	for _,in := range input{
		newInput:=in[:]
		output:= append(output, newInput)
		fmt.Printf("for loop output is %v",output)
	}
	fmt.Printf("the output is %v",output)
	return output
}
