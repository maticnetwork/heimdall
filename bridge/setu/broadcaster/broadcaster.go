package broadcaster

import (
	"context"
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bor "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/tendermint/tendermint/libs/log"
)

// TxBroadcaster uses to broadcast transaction to each chain
type TxBroadcaster struct {
	logger log.Logger

	cliCtx cliContext.CLIContext

	heimdallMutex sync.Mutex
	maticMutex    sync.Mutex
}

// NewTxBroadcaster creates new broadcaster
func NewTxBroadcaster(cdc *codec.Codec) *TxBroadcaster {
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	cliCtx.TrustNode = true

	txBroadcaster := TxBroadcaster{
		logger: util.Logger().With("module", "txBroadcaster"),
		cliCtx: cliCtx,
	}

	return &txBroadcaster
}

// BroadcastToHeimdall broadcast to heimdall
func (tb *TxBroadcaster) BroadcastToHeimdall(msg sdk.Msg) error {
	tb.heimdallMutex.Lock()
	defer tb.heimdallMutex.Unlock()

	// tx encoder
	txEncoder := helper.GetTxEncoder(tb.cliCtx.Codec)
	// chain id
	chainID := helper.GetGenesisDoc().ChainID
	// current address
	address := hmTypes.BytesToHeimdallAddress(helper.GetAddress())

	// fetch from APIs
	var account authTypes.Account
	response, err := util.FetchFromAPI(tb.cliCtx, util.GetHeimdallServerEndpoint(fmt.Sprintf(util.AccountDetailsURL, address)))
	if err != nil {
		tb.logger.Error("Error fetching account from rest-api", "url", util.GetHeimdallServerEndpoint(fmt.Sprintf(util.AccountDetailsURL, address)))
		panic("Error connecting to rest-server, please start server before bridge.")
	}

	// get proposer from response
	if err := tb.cliCtx.Codec.UnmarshalJSON(response.Result, &account); err != nil && len(response.Result) != 0 {
		tb.logger.Error("Error unmarshalling account details", "url", util.GetHeimdallServerEndpoint(fmt.Sprintf(util.AccountDetailsURL, address)))
		return err
	}

	// get account number and sequence
	accNum := account.GetAccountNumber()
	accSeq := account.GetSequence()

	txBldr := authTypes.NewTxBuilderFromCLI().
		WithTxEncoder(txEncoder).
		WithAccountNumber(accNum).
		WithSequence(accSeq).
		WithChainID(chainID)

	txResponse, err := helper.BuildAndBroadcastMsgs(tb.cliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		tb.logger.Error("Error while broadcasting the heimdall transaction", "error", err)
		return err
	}

	// increment account sequence
	// accSeq = accSeq + 1
	tb.logger.Info("Tx sent on heimdall", "txHash", txResponse.TxHash, "accSeq", accSeq, "accNum", accNum)
	tb.logger.Debug("Tx successful on heimdall", "txResponse", txResponse)
	return nil
}

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
		tb.logger.Error("Error while broadcasting the transaction to maticchain", "txHash", signedTx.Hash(), "error", err)
		return err
	}

	return nil
}

// BroadcastToRootchain broadcast to rootchain
func (tb *TxBroadcaster) BroadcastToRootchain() {}
