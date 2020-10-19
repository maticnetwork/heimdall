package main

import (
	"os"

	"github.com/maticnetwork/heimdall/cmd/heimdalld/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
