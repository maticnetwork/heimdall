package main

import (
	"os"

	cmd "github.com/matiknetwork/heimdall/app/heimdalld/cmd"
)

// In main we call the rootCmd
func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
