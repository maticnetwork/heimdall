package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://kovan.infura.io")
	header, err := client.HeaderByNumber(context.Background(), nil)
	if header != nil && err == nil {
		fmt.Println("header fetched %v", header.Number)
	} else {
		fmt.Println("error here %v", err)
	}
}
