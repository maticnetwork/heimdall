package helper

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/types"
)

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

// SendCheckpoint sends checkpoint to rootchain contract
// todo return err
func (c *ContractCaller) SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte) {
	var vote types.CanonicalRLPVote
	err := rlp.DecodeBytes(voteSignBytes, &vote)
	if err != nil {
		Logger.Error("Unable to decode vote while sending checkpoint", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))
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
	GetPubKey().VerifyBytes(voteSignBytes, sigs)

	Logger.Info("Sending new checkpoint", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))

	tx, err := c.RootChainInstance.SubmitHeaderBlock(auth, voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Error while submitting checkpoint", "error", err)
	} else {
		Logger.Info("Submitted new header successfully", "txHash", tx.Hash().String())
	}
}

// CommitSpan sends commit span transaction to validator set contract on bor to change committee for next span
func (c *ContractCaller) CommitSpan(voteSignBytes []byte, sigs []byte, txData []byte, proof []byte) {
	// validator set ABI
	validatorSetABI, err := GetValidatorSetABI()
	if err != nil {
		return
	}

	// commit span
	data, err := validatorSetABI.Pack("commitSpan", voteSignBytes, sigs, txData, proof)
	if err != nil {
		Logger.Error("Unable to pack tx for submitHeaderBlock", "error", err)
		return
	}

	validatorSetAddress := GetValidatorSetAddress()
	auth, err := GenerateAuthObj(GetMaticClient(), ethereum.CallMsg{
		To:   &validatorSetAddress,
		Data: data,
	})
	GetPubKey().VerifyBytes(voteSignBytes, sigs)
	Logger.Info("Submitting new commitee", "Vote", hex.EncodeToString(voteSignBytes), "Sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData), "proof", hex.EncodeToString(proof))

	// commit span
	tx, err := c.ValidatorSetInstance.CommitSpan(auth, voteSignBytes, sigs, txData, proof)
	if err != nil {
		Logger.Error("Error while submitting commit commitee for next span", "error", err)
	} else {
		Logger.Info("Submitted new committee", "txHash", tx.Hash().String())
	}
}
