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
)

var cdc = amino.NewCodec()

func GenerateAuthObj(client *ethclient.Client, callMsg ethereum.CallMsg) (auth *bind.TransactOpts, err error) {
	// get priv key
	pkObject := GetPrivKey()

	// create ecdsa private key
	ecdsaPrivateKey, err := crypto.ToECDSA(pkObject[:])
	if err != nil {
		return
	}

	// from address
	fromAddress := common.BytesToAddress(pkObject.PubKey().Address().Bytes())
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
	// TODO check proposer address from voteSignBytes (have to unmarshall RLP to get address bytes)

	validatorSetInstance, err := GetValidatorSetInstance()
	if err != nil {
		return
	}

	// get ValidatorSet Instance
	validatorSetABI, err := GetValidatorSetABI()
	if err != nil {
		return
	}
	data, err := validatorSetABI.Pack("validate", voteSignBytes, sigs, txData)
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
		Logger.Error("Error while submitting checkpoint", "Error", err)
	} else {
		Logger.Info("Submitted new header successfully ", "txHash", tx.Hash().String())
	}
}
