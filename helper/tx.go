package helper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"
)

var cdc = amino.NewCodec()

func GenerateAuthObj(client *ethclient.Client, callMsg ethereum.CallMsg) (auth *bind.TransactOpts, err error) {
	// get config
	config := GetConfig()
	// load file and unmarshall
	privVal := privval.LoadFilePV(config.ValidatorFilePVPath)
	var pkObject secp256k1.PrivKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privVal.PrivKey.Bytes(), &pkObject)

	// create ecdsa private key
	ecdsaPrivateKey, err := crypto.ToECDSA(pkObject[:])
	if err != nil {
		return
	}

	// from address
	fromAddress := common.BytesToAddress(privVal.Address)
	// fetch gas price
	gasprice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}
	// fetch nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return
	}

	// fetch gas limit
	callMsg.From = fromAddress
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)

	// create auth
	auth = bind.NewKeyedTransactor(ecdsaPrivateKey)
	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(gasLimit) // uint64(gasLimit)

	return
}

func SelectProposer() {
	// get ValidatorSet Instance
	validatorSetInstance, err := GetValidatorSetInstance()
	if err != nil {
		return
	}

	// get ValidatorSet Instance
	validatorSetABI, err := GetValidatorSetABI()
	if err != nil {
		return
	}

	data, err := validatorSetABI.Pack("selectProposer")
	if err != nil {
		Logger.Error("Unable to pack tx for SelectProposer", "error", err)
		return
	}

	validatorAddress := GetValidatorSetAddress()

	// get auth Obj
	auth, err := GenerateAuthObj(GetMainClient(), ethereum.CallMsg{
		To:   &validatorAddress,
		Data: data,
	})

	if err != nil {
		Logger.Error("Unable to draft auth for proposer selection", "error", err)
		return
	}

	// send tx
	tx, err := validatorSetInstance.SelectProposer(auth)
	if err != nil {
		Logger.Error("Unable to send transaction for proposer selection", "error", err)
	} else {
		Logger.Info("Transaction hash for proposing transaction", "txHash", tx.Hash().String())
	}
}
