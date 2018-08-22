package sideBlock
import (
	"fmt"
	"log"

	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

func getBlockDetails(hash string,txroot string,rRoot string) bool{
	client, err := ethclient.Dial("https://testnet.matic.network")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections

	block,err := client.BlockByHash(context.Background(),common.HexToHash(hash))
	if err!=nil {
		fmt.Printf("not found")
		log.Fatal(err)
	}
	if block.Hash().Hex() == hash && block.ReceiptHash().Hex() == rRoot && block.TxHash().Hex() == txroot{
		return true
	}else{
		return false
	}


}