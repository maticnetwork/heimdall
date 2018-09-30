package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"github.com/ethereum/go-ethereum/common"
	contracts "../contracts"
	"github.com/tendermint/tendermint/privval"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"context"
	"math/big"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {


	client, err := ethclient.Dial("https://kovan.infura.io/TJSJL5u9maRXnaZrSvnv")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n *** WE are connected to kovan ****** ")
	contractAddress  := common.HexToAddress("0xd8854cc230a393b31d7f0ef99b4dea14b91a5b41")
	instance,err := contracts.NewContract(contractAddress,client)
	if err!=nil {
		fmt.Println("error error")
	}
	fmt.Println(instance)

	privValObj := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")
	pkBytes := privValObj.PrivKey.Bytes()

	fmt.Printf("Public key is :%v\n",hex.EncodeToString(privValObj.PubKey.Address()))
	fmt.Printf("Private key is : %v\n",hex.EncodeToString(pkBytes[len(pkBytes)-32:]))

	privateKey, err := crypto.ToECDSA(pkBytes[len(pkBytes)-32:])
	if err != nil {
		panic(err)
	}


	fromAddress := common.BytesToAddress(privValObj.Address)

	auth := bind.NewKeyedTransactor(privateKey)
	gasprice ,err:= client.SuggestGasPrice(context.Background())
	if err!=nil {
		fmt.Errorf("Unable to estimate gas")
	}
	nonce,err:= client.PendingNonceAt(context.Background(),fromAddress)
	if err != nil {
		panic(err)
	}
	transferFnSignature := []byte("setHello(string)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

	auth.Nonce=big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasprice
	var data []byte
	hello:=[]byte("vaibhav")
	paddedData := common.LeftPadBytes(hello, 32)
	data = append(data, methodID...)
	data = append(data, paddedData...)
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), auth.GasLimit, auth.GasPrice, data)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Transaction sent with Tx Hash : %s", signedTx.Hash().Hex())

	
	}
