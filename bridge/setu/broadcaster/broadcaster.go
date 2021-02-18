package broadcaster

import (
	"context"
	"fmt"
	"sync"

	bor "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"

	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/tendermint/tendermint/libs/log"
)

// TxBroadcaster uses to broadcast transaction to each chain
type TxBroadcaster struct {
	logger log.Logger

	cliCtx client.Context

	heimdallMutex sync.Mutex
	maticMutex    sync.Mutex

	lastSeqNo uint64
	accNum    uint64
}

// NewTxBroadcaster creates new broadcaster
func NewTxBroadcaster(cdc codec.Marshaler) *TxBroadcaster {
	cliCtx := client.Context{}.WithJSONMarshaler(cdc)
	cliCtx.BroadcastMode = flags.BroadcastAsync

	// current address
	address := hmCommonTypes.BytesToHeimdallAddress(helper.GetAddress())
	account, err := util.GetAccount(cliCtx, address)
	if err != nil {
		panic("Error connecting to rest-server, please start server before bridge.")

	}

	txBroadcaster := TxBroadcaster{
		logger:    util.Logger().With("module", "txBroadcaster"),
		cliCtx:    cliCtx,
		lastSeqNo: account.GetSequence(),
		accNum:    account.GetAccountNumber(),
	}

	return &txBroadcaster
}

//
// BroadcastToHeimdall broadcast to heimdall
func (tb *TxBroadcaster) BroadcastToHeimdall(msg sdk.Msg) error {
	tb.heimdallMutex.Lock()
	defer tb.heimdallMutex.Unlock()

	// tx encoder
	//txEncoder := helper.GetTxEncoder(tb.cliCtx.JSONMarshaler)
	// chain id
	chainID := helper.GetGenesisDoc().ChainID

	// get account number and sequence
	txBldr := tx.NewFactoryCLI(tb.cliCtx, nil).
		WithAccountNumber(tb.accNum).
		WithSequence(tb.lastSeqNo).
		WithChainID(chainID)

	txResponse, err := helper.BuildAndBroadcastMsgs(tb.cliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		tb.logger.Error("Error while broadcasting the heimdall transaction", "error", err)
		// current address
		address := hmCommonTypes.BytesToHeimdallAddress(helper.GetAddress())

		// fetch from APIs
		account, errAcc := util.GetAccount(tb.cliCtx, address)
		if errAcc != nil {
			tb.logger.Error("Error fetching account from rest-api", "url", helper.GetHeimdallServerEndpoint(fmt.Sprintf(util.AccountDetailsURL, helper.GetAddress())))
			return errAcc
		}

		// update seqNo for safety
		tb.lastSeqNo = account.GetSequence()

		return err
	}

	tb.logger.Info("Tx sent on heimdall", "txHash", txResponse.TxHash, "accSeq", tb.lastSeqNo, "accNum", tb.accNum)
	tb.logger.Debug("Tx successful on heimdall", "txResponse", txResponse)
	// increment account sequence
	tb.lastSeqNo += 1
	return nil
}

//
// BroadcastToMatic broadcast to matic
func (tb *TxBroadcaster) BroadcastToMatic(msg bor.CallMsg) error {
	tb.maticMutex.Lock()
	defer tb.maticMutex.Unlock()

	// get matic client
	maticClient := helper.GetMaticClient()

	// get auth
	auth, err := helper.GenerateAuthObj(maticClient, *msg.To, msg.Data)

	if err != nil {
		tb.logger.Error("Error generating auth object", "error", err)
		return err
	}

	// Create the transaction, sign it and schedule it for execution
	rawTx := types.NewTransaction(auth.Nonce.Uint64(), *msg.To, msg.Value, auth.GasLimit, auth.GasPrice, msg.Data)

	// signer
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, rawTx)
	if err != nil {
		tb.logger.Error("Error signing the transaction", "error", err)
		return err
	}

	tb.logger.Info("Sending transaction to bor", "txHash", signedTx.Hash())

	// broadcast transaction
	if err := maticClient.SendTransaction(context.Background(), signedTx); err != nil {
		tb.logger.Error("Error while broadcasting the transaction to maticchain", "error", err)
		return err
	}

	return nil
}

// BroadcastToRootchain broadcast to rootchain
func (tb *TxBroadcaster) BroadcastToRootchain() {}
