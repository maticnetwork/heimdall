package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/maticnetwork/heimdall/helper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		// create codec
		cdc := app.MakeCodec()
		app.MakePulp()

		// queue connector & http client
		_queueConnector := pier.NewQueueConnector(cdc, helper.GetConfig().AmqpURL)
		_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintNodeURL, "/websocket")

		// selected services to start
		services := SelectedServices(cdc, _httpClient, _queueConnector)
		fmt.Println("Services to start - ", services)
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
	},
}

// SelectedServices will select services to start based on set flags --all, --only
func SelectedServices(cdc *codec.Codec, _httpClient *httpClient.HTTP, _queueConnector *pier.QueueConnector) []common.Service {

	services := []common.Service{}

	startAll := viper.GetBool("all")
	onlyServices := viper.GetStringSlice("only")

	if startAll {
		services = []common.Service{
			pier.NewCheckpointer(cdc, _queueConnector, _httpClient),
			pier.NewSyncer(cdc, _queueConnector, _httpClient),
			pier.NewAckService(cdc, _queueConnector, _httpClient),
			pier.NewSpanService(cdc, _queueConnector, _httpClient),
			pier.NewClerkService(cdc, _queueConnector, _httpClient),

			// queue consumer server
			pier.NewConsumerService(cdc, _queueConnector),
		}
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
			case "consumer":
				services = append(services, pier.NewConsumerService(cdc, _queueConnector))
			}
		}
	}

	return services
}

func init() {

	startCmd.Flags().Bool("all", false, "start all bridge services")
	viper.BindPFlag("all", startCmd.Flags().Lookup("all"))

	startCmd.Flags().StringSlice("only", []string{}, "comma separated bridge services to start")
	viper.BindPFlag("only", startCmd.Flags().Lookup("only"))

	rootCmd.AddCommand(startCmd)
}
