package helper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"
)

var cdc = amino.NewCodec()

func GenerateAuthObj(client *ethclient.Client) (auth *bind.TransactOpts) {
	// get config
	config := GetConfig()
	// load file and unmarshall
	privVal := privval.LoadFilePV(config.ValidatorFilePVPath)
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
	// get ValidatorSet Instance
	validatorSetInstance := GetValidatorSetInstance(MainChainClient)
	// get auth Obj
	auth := GenerateAuthObj(MainChainClient)
	// send tx
	tx, err := validatorSetInstance.SelectProposer(auth)
	if err != nil {
		Logger.Error("Unable to send transaction for proposer selection ")
	} else {
		Logger.Info("New Proposer Selected ! TxHash :  %v", tx.Hash().String())
	}

}
