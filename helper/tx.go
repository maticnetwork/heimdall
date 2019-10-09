package helper

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

	Logger.Debug("Sending new checkpoint",
		"vote", hex.EncodeToString(voteSignBytes),
		"sigs", hex.EncodeToString(sigs),
		"txData", hex.EncodeToString(txData))

	tx, err := c.RootChainInstance.SubmitHeaderBlock(auth, voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Error while submitting checkpoint", "error", err)
	} else {
		Logger.Info("Submitted new header successfully", "txHash", tx.Hash().String())
	}
}

// GetCheckpointSign returns sigs input of committed checkpoint tranasction
func (c *ContractCaller) GetCheckpointSign(ctx sdk.Context, txHash common.Hash) ([]byte, []byte, []byte, error) {
	mainChainClient := GetMainClient()
	transaction, isPending, err := mainChainClient.TransactionByHash(ctx, txHash)
	if err != nil {
		Logger.Error("Error while Fetching Transaction By hash from MainChain", "error", err)
		return []byte{}, []byte{}, []byte{}, err
	} else if isPending {
		return []byte{}, []byte{}, []byte{}, errors.New("Transaction is still pending")
	}

	payload := transaction.Data()
	abi := c.RootChainABI
	return UnpackSigAndVotes(payload, abi)
}

// UnpackSigAndVotes Unpacks Sig and Votes from Tx Payload
func UnpackSigAndVotes(payload []byte, abi abi.ABI) ([]byte, []byte, []byte, error) {
	// recover Method from signature and ABI
	method := abi.Methods["submitHeaderBlock"]
	decodedPayload := payload[4:]
	inputDataMap := make(map[string]interface{})
	// unpack method inputs
	err := method.Inputs.UnpackIntoMap(inputDataMap, decodedPayload)
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}
	inputSigs := inputDataMap["sigs"].([]byte)
	txData := inputDataMap["extradata"].([]byte)
	voteSignBytes := inputDataMap["vote"].([]byte)
	Logger.Debug("Sigs of committed checkpoint transaction - ", hex.EncodeToString(inputSigs))
	return voteSignBytes, inputSigs, txData, nil
}
