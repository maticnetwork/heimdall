package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	WaitDuration = 1 * time.Minute
)

// GetStartCmd returns the start command to start bridge
func GetStartCmd() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start bridge server",
		Run: func(cmd *cobra.Command, args []string) {
			logger := pier.Logger.With("module", "bridge")

			// create codec
			cdc := app.MakeCodec()
			app.MakePulp()
			// queue connector & http client
			_queueConnector := pier.NewQueueConnector(cdc, helper.GetConfig().AmqpURL)
			_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintNodeURL, "/websocket")

			// selected services to start
			services := SelectedServices(cdc, _httpClient, _queueConnector)
			if len(services) == 0 {
				panic(fmt.Sprintf("No services selected to start. select services using --all or --only flag"))
			}

			// sync group
			var wg sync.WaitGroup

			// go routine to catch signal
			catchSignal := make(chan os.Signal, 1)
			signal.Notify(catchSignal, os.Interrupt)
			go func() {
				// sig is a ^C, handle it
				for range catchSignal {
					// stop processes
					for _, service := range services {
						service.Stop()
					}

					// stop http client
					_httpClient.Stop()

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
				if pier.IsSynced(cliCtx) {
					logger.Info("Node upto date, starting bridge services")
					break
				} else {
					logger.Info("Waiting for heimdall node to be fully synced")
				}
				time.Sleep(WaitDuration)
			}

			// strt all processes
			for _, service := range services {
				go func(serv common.Service) {
					defer wg.Done()
					// TODO handle error while starting service
					serv.Start()
					<-serv.Quit()
				}(service)
			}
			// wait for all processes
			wg.Add(len(services))
			wg.Wait()
		}}
	startCmd.Flags().Bool("all", false, "start all bridge services")
	viper.BindPFlag("all", startCmd.Flags().Lookup("all"))

	startCmd.Flags().StringSlice("only", []string{}, "comma separated bridge services to start")
	viper.BindPFlag("only", startCmd.Flags().Lookup("only"))
	return startCmd
}

// SelectedServices will select services to start based on set flags --all, --only
func SelectedServices(cdc *codec.Codec, _httpClient *httpClient.HTTP, _queueConnector *pier.QueueConnector) []common.Service {
	services := []common.Service{
		pier.NewConsumerService(cdc, _queueConnector),
	}

	startAll := viper.GetBool("all")
	onlyServices := viper.GetStringSlice("only")

	if startAll {
		services = append(services,
			pier.NewCheckpointer(cdc, _queueConnector, _httpClient),
			pier.NewSyncer(cdc, _queueConnector, _httpClient),
			pier.NewAckService(cdc, _queueConnector, _httpClient),
			pier.NewSpanService(cdc, _queueConnector, _httpClient),
			pier.NewClerkService(cdc, _queueConnector, _httpClient),
		)
	} else {
		for _, service := range onlyServices {
			switch service {
			case "checkpoint":
				services = append(services, pier.NewCheckpointer(cdc, _queueConnector, _httpClient))
			case "syncer":
				services = append(services, pier.NewSyncer(cdc, _queueConnector, _httpClient))
			case "ack":
				services = append(services, pier.NewAckService(cdc, _queueConnector, _httpClient))
			case "span":
				services = append(services, pier.NewSpanService(cdc, _queueConnector, _httpClient))
			case "clerk":
				services = append(services, pier.NewClerkService(cdc, _queueConnector, _httpClient))
			}
		}
	}
	return services
}

func init() {
	rootCmd.AddCommand(GetStartCmd())
}
