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

func SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte) {

	validatorSetInstance, err := GetValidatorSetInstance()
	if err != nil {
		return
	}

	// get ValidatorSet Instance
	validatorSetABI, err := GetValidatorSetABI()
	if err != nil {
		return
	}
	data, err := validatorSetABI.Pack("validate")
	if err != nil {
		Logger.Error("Unable to pack tx for validate", "error", err)
		return
	}

	validatorAddress := GetValidatorSetAddress()
	auth, err := GenerateAuthObj(GetMainClient(), ethereum.CallMsg{
		To:   &validatorAddress,
		Data: data,
	})

	tx, err := validatorSetInstance.Validate(auth, voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Checkpoint Submission Errored", "Error", err)
	} else {
		Logger.Info("Submitted Proof Successfully ", "txHash", tx.Hash().String())
	}

}
