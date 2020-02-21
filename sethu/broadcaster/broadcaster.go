package broadcaster

import (
	"context"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bor "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/util"
	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/tendermint/tendermint/libs/log"
)

type TxBroadcaster struct {

	// logger
	logger log.Logger

	// tx encoder
	cliCtx cliContext.CLIContext
}

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func NewTxBroadcaster(cdc *codec.Codec) *TxBroadcaster {

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	txBroadcaster := TxBroadcaster{
		logger: Logger.With("module", util.Broadcaster),
		cliCtx: cliCtx,
	}

	return &txBroadcaster
}

func (tb *TxBroadcaster) BroadcastToHeimdall(msg sdk.Msg) bool {
	// tx encoder
	txEncoder := helper.GetTxEncoder()
	// chain id
	chainID := helper.GetGenesisDoc().ChainID
	// current address
	address := hmTypes.BytesToHeimdallAddress(helper.GetAddress())
	// fetch from APIs
	var account authTypes.Account
	response, err := pier.FetchFromAPI(tb.cliCtx, pier.GetHeimdallServerEndpoint(fmt.Sprintf(pier.AccountDetailsURL, address)))
	if err != nil {
		tb.logger.Error("Error fetching account from rest-api", "url", pier.GetHeimdallServerEndpoint(fmt.Sprintf(pier.AccountDetailsURL, address)))
		panic("Error connecting to rest-server, please start server before bridge")
	}

	// get proposer from response
	if err := tb.cliCtx.Codec.UnmarshalJSON(response.Result, &account); err != nil && len(response.Result) != 0 {
		panic(err)
	}

	// get account number and sequence
	accNum := account.GetAccountNumber()
	accSeq := account.GetSequence()

	txBldr := authTypes.NewTxBuilderFromCLI().
		WithTxEncoder(txEncoder).
		WithAccountNumber(accNum).
		WithSequence(accSeq).
		WithChainID(chainID)
	if _, err := helper.BuildAndBroadcastMsgs(tb.cliCtx, txBldr, []sdk.Msg{msg}); err != nil {
		tb.logger.Error("Error while broadcasting the heimdall transaction", "error", err)
		return false
	}

	// increment account sequence
	accSeq = accSeq + 1

	return true
}

func (tb *TxBroadcaster) BroadcastToMatic(msg bor.CallMsg) bool {
	maticClient := helper.GetMaticClient()

	// get auth
	auth, err := helper.GenerateAuthObj(maticClient, *msg.To, msg.Data)

	if err != nil {
		tb.logger.Error("Error while fetching the transaction param details", "error", err)
		return false
	}

	// Create the transaction, sign it and schedule it for execution
	rawTx := types.NewTransaction(auth.Nonce.Uint64(), *msg.To, msg.Value, auth.GasLimit, auth.GasPrice, msg.Data)
	// signer
	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, rawTx)
	if err != nil {
		tb.logger.Error("Error while signing the transaction", "error", err)
		return false
	}

	tb.logger.Debug("Sending transaction to bor", "TxHash", signedTx.Hash())

	// broadcast transaction
	if err := maticClient.SendTransaction(context.Background(), signedTx); err != nil {
		tb.logger.Error("Error while broadcasting the transaction", "error", err)
		return false
	}

	// amqp msg
	return true
}

func (tb *TxBroadcaster) BroadcastToRootchain() {

}
