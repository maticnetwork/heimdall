package helper

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"
	"math/big"

	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var cdc = amino.NewCodec()

func GenerateAuthObj(client *ethclient.Client) (auth *bind.TransactOpts) {
	config := GetConfig()
	privVal := privval.LoadFilePV(config.validatorFilePVPath)
	var pkObject secp256k1.PrivKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privVal.PrivKey.Bytes(), &pkObject)

	// create ecdsa private key
	ecdsaPrivateKey, err := crypto.ToECDSA(pkObject[:])
	if err != nil {
		panic(err)
	}

	// from address
	fromAddress := common.BytesToAddress(privVal.Address)
	// fetch gas price
	gasprice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	// fetch nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	// create auth
	auth = bind.NewKeyedTransactor(ecdsaPrivateKey)
	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(3000000)
	return auth
}

func SelectProposer() {
	validatorSetInstance := getValidatorSetInstance(kovanClient)
	auth := GenerateAuthObj(kovanClient)
	tx, err := validatorSetInstance.SelectProposer(auth)
	if err != nil {
		fmt.Printf("Unable to send transaction for proposer selection ")
	}
	fmt.Printf("New Proposer Selected ! %v", tx)
}
