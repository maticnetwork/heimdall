package cmd

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"os/signal"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTx "github.com/maticnetwork/heimdall/checkpoint/rest"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	tendermintProposerBridge = "TendermintProposerBridge"
	defaultPollInterval      = 5 * 1000                // in milliseconds
	defaultCheckpointLength  = 256                     // checkpoint number starts with 0, so length = defaultCheckpointLength -1
	maxCheckpointLength      = 4096                    // max blocks in one checkpoint
	defaultForcePushInterval = maxCheckpointLength * 2 // in seconds (4096 * 2 seconds)
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		// initialize tendermint viper config
		InitTendermintViperConfig(cmd)

		// start bridge
		startBridge()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Bridge to propose
type Bridge struct {
	// Base service
	common.BaseService

	// ETH client
	MaticClient *ethclient.Client
	// ETH RPC client
	MaticRPCClient *rpc.Client
	// Mainchain client
	MainClient *ethclient.Client
	// Rootchain instance
	RootChainInstance *rootchain.Rootchain
	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewBridge returns new service object
func NewBridge() *Bridge {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", tendermintProposerBridge)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// creating bridge object
	bridge := &Bridge{
		MaticClient:       helper.GetMaticClient(),
		MaticRPCClient:    helper.GetMaticRPCClient(),
		MainClient:        helper.GetMainClient(),
		RootChainInstance: rootchainInstance,
		HeaderChannel:     make(chan *types.Header),
	}

	bridge.BaseService = *common.NewBaseService(logger, tendermintProposerBridge, bridge)
	return bridge
}

// StartHeaderProcess starts header process when they get new header
func (bridge *Bridge) StartHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-bridge.HeaderChannel:
			bridge.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

func (bridge *Bridge) sendRequest(newHeader *types.Header) {
	bridge.Logger.Debug("New block detected", "blockNumber", newHeader.Number)
	lastCheckpointEnd, err := bridge.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		bridge.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}

	latest := newHeader.Number.Uint64()
	start := lastCheckpointEnd.Uint64()
	var end uint64

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := latest - start + 1

	// process if diff > 0 (positive)
	if diff > 0 {
		expectedDiff := diff - diff%defaultCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}

		// cap with max checkpoint length
		if expectedDiff > maxCheckpointLength-1 {
			expectedDiff = maxCheckpointLength - 1
		}

		// get end result
		end = expectedDiff + start

		bridge.Logger.Debug("Calculating checkpoint eligibility", "latest", latest, "start", start, "end", end)
	}

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < defaultCheckpointLength) {
		currentHeaderBlockNumber, err := bridge.RootChainInstance.CurrentHeaderBlock(nil)
		if err != nil {
			bridge.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
			return
		}

		// fetch current header block
		currentHeaderBlock, err := bridge.RootChainInstance.GetHeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(1)))
		if err != nil {
			bridge.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			return
		}

		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > defaultForcePushInterval {
			bridge.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval)
			end = latest
		}
	}

	if end == 0 || start >= end {
		return
	}

	// Get root hash
	root := checkpoint.GetHeaders(start, end)
	bridge.Logger.Info("New checkpoint header created", "latest", latest, "start", start, "end", end, "root", root)

	// TODO submit checkcoint
	txBytes, err := checkpointTx.CreateTxBytes(checkpointTx.EpochCheckpoint{
		RootHash:   root,
		StartBlock: start,
		EndBlock:   end,
	})

	if err != nil {
		bridge.Logger.Error("Error while creating tx bytes", "error", err)
		return
	}

	resp, err := checkpointTx.SendTendermintRequest(cliContext.NewCLIContext(), txBytes)
	if err != nil {
		bridge.Logger.Error("Error while sending request to Tendermint", "error", err)
		return
	}

	bridge.Logger.Error("Checkpoint sent successfully", "hash", hex.EncodeToString(resp.Hash), "start", start, "end", end, "root", root)
}

// OnStart starts new block subscription
func (bridge *Bridge) OnStart() error {
	bridge.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	bridge.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	bridge.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go bridge.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := bridge.MaticClient.SubscribeNewHead(ctx, bridge.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go bridge.StartPolling(ctx, defaultPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go bridge.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	bridge.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (bridge *Bridge) OnStop() {
	bridge.BaseService.OnStop() // Always call the overridden method.

	// cancel subscription if any
	bridge.cancelSubscription()

	// cancel header process
	bridge.cancelHeaderProcess()
}

func (bridge *Bridge) StartPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := bridge.MaticClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				bridge.HeaderChannel <- header
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (bridge *Bridge) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			bridge.Logger.Error("Error while subscribing new blocks", "error", err)
			bridge.Stop()

			// cancel subscription
			bridge.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func startBridge() {
	bridge := NewBridge()
	bridge.Start()

	// go routine to catch signal
	catchSignal := make(chan os.Signal, 1)
	signal.Notify(catchSignal, os.Interrupt)
	go func() {
		// sig is a ^C, handle it
		for sig := range catchSignal {
			// print sig
			bridge.Logger.Debug("Captured and topping profiler and exiting", "sig", sig)

			// stop
			bridge.Stop()

			// exit
			os.Exit(1)
		}
	}()

	// wait for bridge to quiet
	bridge.Wait()
}
