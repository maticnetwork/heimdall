package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"log"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://ropsten.infura.io")
	if err != nil {
		fmt.Printf("Error while creating main chain client", "error", err)
	}
	receipt,err:=client.TransactionReceipt(context.Background(),common.HexToHash("0xcb2712e485f4680decd0f9a00e7d26c35286970f275f2f30f50c9bb5e672140b"))
	if err!=nil{
		fmt.Printf("Unable to get transaction receiipt by hash")

	}
	contractAbi, err := abi.JSON(strings.NewReader(string(rootchain.RootchainABI)))
	if err != nil {
		log.Fatal(err)
	}
	for _, vLog := range receipt.Logs {
		var event rootchain.RootchainNewHeaderBlock
		err := contractAbi.Unpack(&event, "NewHeaderBlock", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(event.Start.String()))   // foo
		fmt.Println(string(event.End.String())) // bar
		fmt.Println(common.BytesToHash(event.Root[:]).String())
		fmt.Println(event.Proposer.String())
	}
	type txExtraInfo struct {
		BlockNumber *string         `json:"blockNumber,omitempty"`
		BlockHash   *common.Hash    `json:"blockHash,omitempty"`
		From        *common.Address `json:"from,omitempty"`
	}
	type rpcTransaction struct {
		tx *types.Transaction
		txExtraInfo
	}
	var json *rpcTransaction
	maticRPCClient, err := rpc.Dial("https://ropsten.infura.io")
	err=maticRPCClient.CallContext(context.Background(),  &json, "eth_getTransactionByHash",common.HexToHash("0xcb2712e485f4680decd0f9a00e7d26c35286970f275f2f30f50c9bb5e672140b"))
	if err!=nil{
		fmt.Println("error %v",err)
	}
	var blocknum big.Int
	blocknumber,ok:=blocknum.SetString(*json.BlockNumber,0)
	if !ok{
		fmt.Println("Not found: %v",ok)
	}
	fmt.Println("receipt %v",blocknumber)
	fmt.Printf("Transaction fetched %v ",receipt.Logs)
}
