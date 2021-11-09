package main

import (
	"fmt"
	"os"

	"github.com/maticnetwork/heimdall/bridge/cmd"
)

func main() {
	rootCmd := cmd.BridgeCommands()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
