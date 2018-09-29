package main

import (
	"github.com/tendermint/tendermint/privval"
	"fmt"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"github.com/ethereum/go-ethereum/common"
	contracts "../contracts"
	)

func main() {
	//privVal := privval.LoadFilePV(config.DefaultBaseConfig().PrivValidatorFile())
	privValObj := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")
	sig,err:=privValObj.PrivKey.Sign(crypto.Keccak256([]byte("hello")))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Signature is %v",hex.EncodeToString(sig))
	client, err := ethclient.Dial("https://kovan.infura.io/TJSJL5u9maRXnaZrSvnv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("*** WE are connected to kovan ******")
	contractAddress  := common.HexToAddress("0xd8854cc230a393b31d7f0ef99b4dea14b91a5b41")
	instance,err := contracts.NewContract(contractAddress,client)
	if err!=nil {
		fmt.Println("error error")
	}
	fmt.Println(instance)
	
	}
