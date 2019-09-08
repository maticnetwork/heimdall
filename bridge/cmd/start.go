package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"
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
		_queueConnector := pier.NewQueueConnector(helper.GetConfig().AmqpURL, "hq", "bq", "cq")
		_httpClient := httpClient.NewHTTP(helper.GetConfig().TendermintNodeURL, "/websocket")

		// create codec
		cdc := app.MakeCodec()

		services := [...]common.Service{
			pier.NewCheckpointer(cdc, _queueConnector, _httpClient),
			pier.NewSyncer(cdc, _queueConnector, _httpClient),
			pier.NewAckService(cdc, _queueConnector, _httpClient),
			pier.NewSpanService(cdc, _queueConnector, _httpClient),

			// queue consumer server
			pier.NewConsumerService(cdc, _queueConnector),
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

func init() {
	rootCmd.AddCommand(startCmd)
}
