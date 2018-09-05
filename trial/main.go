package main

import (
	"fmt"
	"log"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/xsleonard/go-merkle"
	"encoding/hex"
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
	if start>end{
		return ""
	}
	current:=start
	var result [][32]byte
	for current <= end {
		//TODO run this in different goroutines and use channels to fetch results(how to maintian order)
		blockheader,err:=client.HeaderByNumber(context.Background(),big.NewInt(int64(current)))
		if err!=nil {
			fmt.Printf("not found")
			log.Fatal(err)
		}
		headerBytes := appendBytes32(	blockheader.Number.Bytes(),
										blockheader.Time.Bytes(),
										blockheader.TxHash.Bytes(),
										blockheader.ReceiptHash.Bytes() )

		fmt.Printf("attention !!! %v \n",getsha3frombyte("abcd"))
		header:= getsha3frombyte(hex.EncodeToString(headerBytes))
		var arr [32]byte
		copy(arr[:], header)
		result = append(result, arr)
		current++
	}
	fmt.Println("------")
	fmt.Printf(" the headers are ",result)
	merkelData:=convert(result)
	fmt.Println("------")
	fmt.Printf("merkel data is %v",merkelData)
	fmt.Println("------")
	tree := merkle.NewTree()
	err := tree.Generate(merkelData,sha3.New256())
	if err != nil {
		fmt.Println("*********ERROR***********")
		log.Fatal(err)
	}
	fmt.Printf("Root: %v\n", tree.Root())// return the hash of root
	fmt.Println(hex.EncodeToString(tree.Root().Hash))
	return hex.EncodeToString(tree.Root().Hash)
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

func appendBytes32(data... []byte) []byte {
	var result []byte
	for _, v := range data {
		paddedV, err := convertTo32(v)
		if err == nil {
			result = append(result, paddedV[:]...)
		}
	}
	return result
}
func getsha3frombyte(input string) string{
	hash := sha3.NewKeccak256()
	fmt.Println(input)
	v,err:=hex.DecodeString(input)
	if err!=nil{
		fmt.Println("dsd")
		log.Fatal(err)
	}
	fmt.Println("v is ")
	fmt.Println(v)
	var buf []byte
	hash.Write(v)
	buf = hash.Sum(buf)
	fmt.Printf("the input was %v and the output is %v",input,hex.EncodeToString(buf))
	return hex.EncodeToString(buf)
}

