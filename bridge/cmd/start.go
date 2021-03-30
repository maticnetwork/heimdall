package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/maticnetwork/heimdall/bridge/setu/processor"

	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bridge/setu/listener"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/service"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	waitDuration = 1 * time.Minute
	logLevel     = "log_level"
	keyName      = "key-name"
	AccountAddr  = "account-address"
)

// GetStartCmd returns the start command to start bridge
func GetStartCmd() *cobra.Command {
	var logger = helper.Logger.With("module", "bridge/cmd/")
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start bridge server",
		Run: func(cmd *cobra.Command, args []string) {

			nodeKeyName := viper.GetString(keyName)
			fmt.Println("key name  ", nodeKeyName, " ", len(nodeKeyName))
			if nodeKeyName == "" || len(nodeKeyName) == 0 {
				panic(fmt.Sprintf("Validator key name is required"))
			}
			// create codec
			cdc, _ := app.MakeCodecs()
			// encoding
			encoding := app.MakeEncodingConfig()
			// queue connector & http client
			_queueConnector := queue.NewQueueConnector(helper.GetConfig().AmqpURL)
			_queueConnector.StartWorker()

			_httpClient, _ := httpClient.New(helper.GetConfig().TendermintRPCUrl, "/websocket")

			// cli context
			cliCtx := client.Context{}.WithJSONMarshaler(cdc)
			chainID := helper.GetGenesisDoc().ChainID

			acctAddr, _ := sdk.AccAddressFromHex(viper.GetString(AccountAddr))
			cliCtx = cliCtx.WithNodeURI(helper.GetConfig().TendermintRPCUrl).
				WithClient(_httpClient).
				WithAccountRetriever(authtypes.AccountRetriever{}).
				WithInterfaceRegistry(encoding.InterfaceRegistry).
				WithTxConfig(encoding.TxConfig).
				WithFromAddress(acctAddr).
				WithFromName(nodeKeyName).
				WithChainID(chainID).
				WithSkipConfirmation(true)

			cliCtx.BroadcastMode = flags.BroadcastAsync

			_txBroadcaster := broadcaster.NewTxBroadcaster(cliCtx, cdc, cmd.Flags())

			// params context
			_paramsContext := util.NewParamsContext(cliCtx)

			// selected services to start
			var services []service.Service
			services = append(services,
				listener.NewListenerService(cliCtx, cdc, _queueConnector, _httpClient),
				processor.NewProcessorService(cliCtx, cdc, _queueConnector, _httpClient, _txBroadcaster, _paramsContext),
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
					for _, nService := range services {
						if err := nService.Stop(); err != nil {
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

			// start all processes
			for _, nService := range services {
				go func(serv service.Service) {
					defer wg.Done()
					// TODO handle error while starting service
					if err := serv.Start(); err != nil {
						logger.Error("GetStartCmd | serv.Start", "Error", err)
					}
					<-serv.Quit()
				}(nService)
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

	startCmd.Flags().String(AccountAddr, "", "node genesis account address ")
	if err := viper.BindPFlag(AccountAddr, startCmd.Flags().Lookup(AccountAddr)); err != nil {
		logger.Error("GetStartCmd | BindPFlag | "+AccountAddr, "Error", err)
	}

	startCmd.Flags().String(keyName, "", "Validator key name in keyring")
	if err := viper.BindPFlag(keyName, startCmd.Flags().Lookup(keyName)); err != nil {
		logger.Error("GetStartCmd | BindPFlag | "+keyName, "Error", err)
	}

	startCmd.Flags().Bool("all", false, "start all bridge services")
	if err := viper.BindPFlag("all", startCmd.Flags().Lookup("all")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | all", "Error", err)
	}

	startCmd.Flags().StringSlice("only", []string{}, "comma separated bridge services to start")
	if err := viper.BindPFlag("only", startCmd.Flags().Lookup("only")); err != nil {
		logger.Error("GetStartCmd | BindPFlag | only", "Error", err)
	}

	startCmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	if err := viper.BindPFlag(flags.FlagKeyringBackend, startCmd.Flags().Lookup(flags.FlagKeyringBackend)); err != nil {
		logger.Error("GetStartCmd | BindPFlag | "+flags.FlagKeyringBackend, "Error", err)
	}
	return startCmd
}

func init() {
	rootCmd.AddCommand(GetStartCmd())
}
