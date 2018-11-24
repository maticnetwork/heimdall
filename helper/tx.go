package helper

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"

	"bytes"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
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

// send checkpoint to rootchain contract
// todo return err
func SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte) {
	var vote types.CanonicalRLPVote
	err := rlp.DecodeBytes(voteSignBytes, &vote)
	if err != nil {
		Logger.Error("Unable to decode vote while sending checkpoint", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))
	}

	rootchainInstance, err := GetRootChainInstance()
	if err != nil {
		return
	}

	// get stakeManager Instance
	rootchainABI, err := GetRootChainABI()
	if err != nil {
		return
	}

	data, err := rootchainABI.Pack("submitHeaderBlock", voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Unable to pack tx for submitHeaderBlock", "error", err)
		return
	}

	rootChainAddress := GetRootChainAddress()
	auth, err := GenerateAuthObj(GetMainClient(), ethereum.CallMsg{
		To:   &rootChainAddress,
		Data: data,
	})

	// get validator address
	validatorAddress := GetPubKey().Address().Bytes()

	// todo replace proposer check here with proposer check from store
	if !bytes.Equal(validatorAddress, vote.Proposer) {
		Logger.Info("You are not proposer", "proposer", hex.EncodeToString(vote.Proposer), "validator", hex.EncodeToString(validatorAddress))
	} else {
		Logger.Info("We are proposer. Sending new checkpoint", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))

		tx, err := rootchainInstance.SubmitHeaderBlock(auth, voteSignBytes, sigs, txData)
		if err != nil {
			Logger.Error("Error while submitting checkpoint", "error", err)
		} else {
			Logger.Info("Submitted new header successfully", "txHash", tx.Hash().String())
		}
	}
}
