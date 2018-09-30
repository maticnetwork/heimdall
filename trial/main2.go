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
	"crypto/ecdsa"
	"math/big"
)

func main() {

	// Connecting to contract and creating an instance

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

	//TODO create raw transaction and sign using validator private key
	// -----------------
	// https://ethereum.stackexchange.com/questions/16472/signing-a-raw-transaction-in-go
	//
	//privValObj := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")
	//ecdsaPk, _ := crypto.ToECDSA(privValObj.PrivKey.Bytes()[:])
	//auth := bind.NewKeyedTransactor(ecdsaPk)
	//fmt.Printf("prival is %v\n,%v\n,%v",privValObj,ecdsaPk,auth)
	//


	//------------------


	/*
	Connecting to priv_validator.json and signing using that private key
	 */
	//privVal := privval.LoadFilePV(config.DefaultBaseConfig().PrivValidatorFile())


	privValObj := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")
	pkBytes := privValObj.PrivKey.Bytes()

	fmt.Printf("Public key is :%v\n",hex.EncodeToString(privValObj.PubKey.Address()))
	fmt.Printf("Private key is : %v\n",hex.EncodeToString(pkBytes[len(pkBytes)-32:]))
	privateKey,err:= crypto.HexToECDSA(string(hex.EncodeToString(pkBytes[len(pkBytes)-32:])))
	if err != nil {
		panic(err)
	}
	publicKey:= privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Errorf(" Unable to cast ")
	}

	fromaddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	auth := bind.NewKeyedTransactor(privateKey)
	gasprice ,err:= client.SuggestGasPrice(context.Background())
	if err!=nil {
		fmt.Errorf("Unable to estimate gas")
	}
	nonce,err:= client.PendingNonceAt(context.Background(),fromaddress)
	if err != nil {
		panic(err)
	}
	auth.Nonce=big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasprice
	tx,err:=instance.SetHello(auth,"vaibhav")
	if err!=nil {
		fmt.Errorf(" Unable to send transaction ERROR %v",err)
	}
	fmt.Printf("Transaction send successfully ! %v/n",tx)


	
	}
