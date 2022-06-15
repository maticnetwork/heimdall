package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/version"
)

const (
	bridgeDBFlag   = "bridge-db"
	borChainIDFlag = "bor-chain-id"
	logsTypeFlag   = "logs-type"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "heimdall-bridge",
	Short: "Heimdall bridge deamon",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use != version.Cmd.Use {
			// initialize tendermint viper config
			InitTendermintViperConfig(cmd)
		}
	},
}

// InitTendermintViperConfig sets global viper configuration needed to heimdall
func InitTendermintViperConfig(cmd *cobra.Command) {
	tendermintNode, _ := cmd.Flags().GetString(helper.NodeFlag)
	homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)
	withHeimdallConfigValue, _ := cmd.Flags().GetString(helper.WithHeimdallConfigFlag)
	bridgeDBValue, _ := cmd.Flags().GetString(bridgeDBFlag)
	borChainIDValue, _ := cmd.Flags().GetString(borChainIDFlag)
	logsTypeValue, _ := cmd.Flags().GetString(logsTypeFlag)

	// bridge-db directory (default storage)
	if bridgeDBValue == "" {
		bridgeDBValue = filepath.Join(homeValue, "bridge", "storage")
	}

	// set to viper
	viper.Set(helper.NodeFlag, tendermintNode)
	viper.Set(helper.HomeFlag, homeValue)
	viper.Set(helper.WithHeimdallConfigFlag, withHeimdallConfigValue)
	viper.Set(bridgeDBFlag, bridgeDBValue)
	viper.Set(borChainIDFlag, borChainIDValue)
	viper.Set(logsTypeFlag, logsTypeValue)

	// start heimdall config
	helper.InitHeimdallConfig("")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(version.Cmd)
	rootCmd.PersistentFlags().StringP(helper.NodeFlag, "n", "tcp://localhost:26657", "Node to connect to")
	rootCmd.PersistentFlags().String(helper.HomeFlag, os.ExpandEnv("$HOME/.heimdalld"), "directory for config and data")
	rootCmd.PersistentFlags().String(
		helper.WithHeimdallConfigFlag,
		"",
		"Heimdall config file path (default <home>/config/heimdall-config.json)",
	)
	// bridge storage db
	rootCmd.PersistentFlags().String(
		bridgeDBFlag,
		"",
		"Bridge db path (default <home>/bridge/storage)",
	)
	// bridge chain id
	rootCmd.PersistentFlags().String(
		borChainIDFlag,
		helper.DefaultBorChainID,
		"Bor chain id",
	)
	// bridge logging type
	rootCmd.PersistentFlags().String(
		logsTypeFlag,
		helper.DefaultLogsType,
		"Use json logger",
	)

	var logger = helper.Logger.With("module", "bridge/cmd/")
	// bind all flags with viper
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		logger.Error("init | BindPFlag | rootCmd.Flags", "Error", err)
	}
}
