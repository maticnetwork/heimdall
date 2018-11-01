package main

import (
	"context"
	"encoding/hex"
	"fmt"
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
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTx "github.com/maticnetwork/heimdall/checkpoint/rest"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	tendermintProposerBridge = "TendermintProposerBridge"
	defaultPollInterval      = 5000
	defaultCheckpointLength  = 256
)

var rootCmd = &cobra.Command{
	Use:   "heimdall-bridge",
	Short: "Heimdall bridge deamon",
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

	start := big.NewInt(0)
	end := big.NewInt(0)

	// add 1 if lastCheckpointEnd > 0
	if lastCheckpointEnd.Sign() > 0 {
		start = start.Add(lastCheckpointEnd, big.NewInt(1))
	}

	diff := big.NewInt(0)
	diff = diff.Sub(newHeader.Number, start)

	// process if diff > 0 (positive)
	if diff.Sign() > 0 {
		if diff.Uint64() >= defaultCheckpointLength {
			end = end.Add(start, big.NewInt(defaultCheckpointLength-1))
			bridge.Logger.Debug("start - end >= defaultCheckpointLength", "latest", newHeader.Number, "start", start, "end", end, "defaultCheckpointLength", defaultCheckpointLength)
		} else {
			bridge.Logger.Debug("start - end < defaultCheckpointLength", "latest", newHeader.Number, "start", start, "defaultCheckpointLength", defaultCheckpointLength)
			// TODO wait for last checkpoint. If checkpoint time > 10 min create checkpoint with remaining blocks
		}
	}

	if end.Sign() < 0 {
		return
	}

	// Get root hash
	root := checkpoint.GetHeaders(start.Uint64(), end.Uint64())
	bridge.Logger.Info("New checkpoint header created", "start", start, "end", end, "root", root)

	// TODO submit checkcoint
	txBytes, err := checkpointTx.CreateTxBytes(checkpointTx.EpochCheckpoint{
		RootHash:   root,
		StartBlock: start.Uint64(),
		EndBlock:   end.Uint64(),
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

func main() {
	rootCmd.PersistentFlags().StringP(helper.NodeFlag, "n", "tcp://localhost:26657", "Node to connect to")
	rootCmd.PersistentFlags().String(helper.HomeFlag, os.ExpandEnv("$HOME/.heimdalld"), "directory for config and data")
	rootCmd.PersistentFlags().String(
		helper.WithHeimdallConfigFlag,
		"",
		"Heimdall config file path (default <home>/config/heimdall-config.json)",
	)

	// bind all flags with viper
	viper.BindPFlags(rootCmd.Flags())

	// add start command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start bridge server",
		Run: func(cmd *cobra.Command, args []string) {
			tendermintNode, _ := cmd.Flags().GetString(helper.NodeFlag)
			homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)
			withHeimdallConfigValue, _ := cmd.Flags().GetString(helper.WithHeimdallConfigFlag)

			// set to viper
			viper.Set(helper.NodeFlag, tendermintNode)
			viper.Set(helper.HomeFlag, homeValue)
			viper.Set(helper.WithHeimdallConfigFlag, withHeimdallConfigValue)

			// start heimdall config
			helper.InitHeimdallConfig()

			// start bridge
			startBridge()
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
