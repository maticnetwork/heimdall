package helper

import (
	"context"
	"encoding/hex"
	"math/big"

	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi/bind"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rlp"
	"github.com/tendermint/tendermint/types"
)

func GenerateAuthObj(client *ethclient.Client, address common.Address, data []byte) (auth *bind.TransactOpts, err error) {
	// generate call msg
	callMsg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

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
func (c *ContractCaller) SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte) error {
	var vote types.CanonicalRLPVote
	err := rlp.DecodeBytes(voteSignBytes, &vote)
	if err != nil {
		Logger.Error("Unable to decode vote while sending checkpoint", "vote", hex.EncodeToString(voteSignBytes), "sigs", hex.EncodeToString(sigs), "txData", hex.EncodeToString(txData))
		return err
	}

	data, err := c.RootChainABI.Pack("submitHeaderBlock", voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Unable to pack tx for submitHeaderBlock", "error", err)
		return err
	}

	rootChainAddress := GetRootChainAddress()
	auth, err := GenerateAuthObj(GetMainClient(), rootChainAddress, data)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return err
	}
	GetPubKey().VerifyBytes(voteSignBytes, sigs)

	Logger.Debug("Sending new checkpoint",
		"vote", hex.EncodeToString(voteSignBytes),
		"sigs", hex.EncodeToString(sigs),
		"txData", hex.EncodeToString(txData))

	tx, err := c.RootChainInstance.SubmitHeaderBlock(auth, voteSignBytes, sigs, txData)
	if err != nil {
		Logger.Error("Error while submitting checkpoint", "error", err)
		return err
	} else {
		Logger.Info("Submitted new header successfully", "txHash", tx.Hash().String())
	}
	return nil
}

// StakeFor stakes for a validator
func (c *ContractCaller) StakeFor(val common.Address, stakeAmount *big.Int, feeAmount *big.Int, acceptDelegation bool) error {
	// stake the amount
	stakeManagerAddress := GetStakeManagerAddress()

	signerPubkey := GetPubKey()
	signerAddress := common.HexToAddress(signerPubkey.Address().String())

	data, err := c.StakeManagerABI.Pack("stakeFor", val, stakeAmount, feeAmount, signerAddress, acceptDelegation)
	if err != nil {
		Logger.Error("Unable to pack tx for stakeFor", "error", err)
		return err
	}

	auth, err := GenerateAuthObj(GetMainClient(), stakeManagerAddress, data)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return err
	}

	// stake for stake manager
	tx, err := c.StakeManagerInstance.StakeFor(
		auth,
		val,
		stakeAmount,
		feeAmount,
		signerAddress,
		acceptDelegation,
	)

	if err != nil {
		Logger.Error("Error while submitting stake", "error", err)
		return err
	}

	Logger.Info("Submitted stake sucessfully", "txHash", tx.Hash().String())
	return nil
}

// ApproveTokens approves matic token for stake
func (c *ContractCaller) ApproveTokens(amount *big.Int) error {
	tokenAddress := GetMaticTokenAddress()
	stakeManager := GetStakeManagerAddress()

	data, err := c.MaticTokenABI.Pack("approve", stakeManager, amount)
	if err != nil {
		Logger.Error("Unable to pack tx for approve", "error", err)
		return err
	}

	auth, err := GenerateAuthObj(GetMainClient(), tokenAddress, data)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return err
	}

	tx, err := c.MaticTokenInstance.Approve(auth, stakeManager, amount)
	if err != nil {
		Logger.Error("Error while approving approve", "error", err)
		return err
	}

	Logger.Info("Sent approve tx sucessfully", "txHash", tx.Hash().String())
	return nil
}
