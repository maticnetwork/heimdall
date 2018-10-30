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
	"github.com/tendermint/tendermint/types"
	"github.com/ethereum/go-ethereum/rlp"
	"bytes"
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
	var vote types.CanonicalRLPVote
	err := rlp.DecodeBytes(voteSignBytes, &vote)
	if err != nil {
		Logger.Error("Unable to decode vote while sending checkpoint","vote",string(voteSignBytes))
	}

	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		return
	}

	// get stakeManager Instance
	stakeManagerABI, err := GetStakeManagerABI()
	if err != nil {
		return
	}

	data, err := stakeManagerABI.Pack("validate", voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Unable to pack tx for validate", "error", err)
		return
	}

	stakeManagerAddress := GetStakeManagerAddress()
	auth, err := GenerateAuthObj(GetMainClient(), ethereum.CallMsg{
		To:   &stakeManagerAddress,
		Data: data,
	})

	if !bytes.Equal(GetPrivKey().PubKey().Address().Bytes(),vote.Proposer){
		Logger.Info("You are not proposer","Proposer",vote.Proposer,"Validator",pkObject.PubKey().Address().Bytes())

	}else{

		Logger.Info("We are proposer , sending checkpoint","Vote",string(voteSignBytes),"Sigs",string(sigs),"TxData",string(txData))

		tx, err := stakeManagerInstance.Validate(auth, voteSignBytes, sigs, txData)
		if err != nil {
			Logger.Error("Error while submitting checkpoint", "Error", err)

		} else {
			Logger.Info("Submitted new header successfully ", "txHash", tx.Hash().String())

		}

	}


}
