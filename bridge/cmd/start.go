package cmd

import (
	"os"
	"os/signal"
	"sync"

	"github.com/spf13/cobra"

	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/tendermint/tendermint/libs/common"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start bridge server",
	Run: func(cmd *cobra.Command, args []string) {
		services := [...]common.BaseService{
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
			for sig := range catchSignal {
				// stop processes
				for service: range services {
					service.Stop()
				}

				// exit
				os.Exit(1)
			}
		}()

		// strt all processes
		for _, service : range services {
			go func(serv) {
				defer wg.Done()
				serv.Start()
				serv.Wait()
			}(service)
		}

		// wait for all processes
		wg.Add(len(services))
		wg.Wait()
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
