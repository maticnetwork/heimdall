package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/listener"
	"github.com/maticnetwork/heimdall/bridge/setu/processor"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	waitDuration = 1 * time.Minute
	logLevel     = "log_level"
)

// GetStartCmd returns the start command to start bridge
func GetStartCmd() *cobra.Command {
	var logger = helper.Logger.With("module", "bridge/cmd/")
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start bridge server",
		Run: func(cmd *cobra.Command, args []string) {

			// create codec
			cdc := app.MakeCodec()
			// queue connector & http client
			_queueConnector := queue.NewQueueConnector(helper.GetConfig().AmqpURL)
			_queueConnector.StartWorker()

			_txBroadcaster := broadcaster.NewTxBroadcaster(cdc)
			_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintRPCUrl, "/websocket")

			// selected services to start
			services := []common.Service{}
			services = append(services,
				listener.NewListenerService(cdc, _queueConnector, _httpClient),
				processor.NewProcessorService(cdc, _queueConnector, _httpClient, _txBroadcaster),
			)

			// sync group
			var wg sync.WaitGroup

			// go routine to catch signal
			catchSignal := make(chan os.Signal, 1)
			signal.Notify(catchSignal, os.Interrupt, syscall.SIGTERM)
			go func() {
				// sig is a ^C, handle it
				for range catchSignal {
					// stop processes
					logger.Info("Received stop signal - Stopping all services")
					for _, service := range services {
						if err := service.Stop(); err != nil {
							logger.Error("GetStartCmd | service.Stop", "Error", err)
						}
					}

					// stop http client
					if err := _httpClient.Stop(); err != nil {
						logger.Error("GetStartCmd | _httpClient.Stop", "Error", err)
					}

					// stop db instance
					util.CloseBridgeDBInstance()

					// exit
					os.Exit(1)
				}
			}()

			// Start http client
			err := _httpClient.Start()
			if err != nil {
				panic(fmt.Sprintf("Error connecting to server %v", err))
			}

			// cli context
			cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
			cliCtx.BroadcastMode = client.BroadcastAsync
			cliCtx.TrustNode = true

			// start bridge services only when node fully synced
			for {
				if !util.IsCatchingUp(cliCtx) {
					logger.Info("Node upto date, starting bridge services")
					break
				} else {
					logger.Info("Waiting for heimdall to be synced")
				}
				time.Sleep(waitDuration)
			}

			// strt all processes
			for _, service := range services {
				go func(serv common.Service) {
					defer wg.Done()
					// TODO handle error while starting service
					if err := serv.Start(); err != nil {
						logger.Error("GetStartCmd | serv.Start", "Error", err)
					}
					<-serv.Quit()
				}(service)
			}
			// wait for all processes
			wg.Add(len(services))
			wg.Wait()
		}}

	// log level
	startCmd.Flags().String(logLevel, "info", "Log level for bridge")
	if err := viper.BindPFlag(logLevel, startCmd.Flags().Lookup(logLevel)); err != nil {
		logger.Error("GetStartCmd | BindPFlag | logLevel", "Error", err)
	}

	startCmd.Flags().Bool("all", false, "start all bridge services")
	if err := viper.BindPFlag("all", startCmd.Flags().Lookup("all")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | all", "Error", err)
	}

	startCmd.Flags().StringSlice("only", []string{}, "comma separated bridge services to start")
	if err := viper.BindPFlag("only", startCmd.Flags().Lookup("only")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | only", "Error", err)
	}
	return startCmd
}

func init() {
	rootCmd.AddCommand(GetStartCmd())
}
