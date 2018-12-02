package cmd

import (
	"os"
	"os/signal"
	"sync"

	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/common"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		services := [...]common.Service{
			pier.NewMaticCheckpointer(),
			pier.NewChainSyncer(),
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

				// exit
				os.Exit(1)
			}
		}()

		// strt all processes
		for _, service := range services {
			go func(serv common.Service) {
				defer wg.Done()
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
