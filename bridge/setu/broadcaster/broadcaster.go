package broadcaster

import (
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	lastSeqNo uint64
	accNum    uint64
}

// NewTxBroadcaster creates new broadcaster
func NewTxBroadcaster(cdc *codec.Codec) *TxBroadcaster {
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	cliCtx.TrustNode = true

	// current address
	address := hmTypes.BytesToHeimdallAddress(helper.GetAddress())
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

// BroadcastToHeimdall broadcast to heimdall
func (tb *TxBroadcaster) BroadcastToHeimdall(msg sdk.Msg) error {
	tb.heimdallMutex.Lock()
	defer tb.heimdallMutex.Unlock()

	// tx encoder
	txEncoder := helper.GetTxEncoder(tb.cliCtx.Codec)
	// chain id
	chainID := helper.GetGenesisDoc().ChainID

	// get account number and sequence
	txBldr := authTypes.NewTxBuilderFromCLI().
		WithTxEncoder(txEncoder).
		WithAccountNumber(tb.accNum).
		WithSequence(tb.lastSeqNo).
		WithChainID(chainID)

	if tb.cliCtx.Height > helper.TxWithGasHeight {
		txBldr = txBldr.WithFees(helper.GetConfig().HeimdallTxFee).
			WithGasAdjustment(helper.GetConfig().GasAdjustment).
			WithSimulateAndExecute(true)
	}

	txResponse, err := helper.BuildAndBroadcastMsgs(tb.cliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		tb.logger.Error("Error while broadcasting the heimdall transaction", "error", err)

		// current address
		address := hmTypes.BytesToHeimdallAddress(helper.GetAddress())

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
