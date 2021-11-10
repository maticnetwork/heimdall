package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/helper"
)

const (
	bridgeDBFlag   = "bridge-db"
	borChainIDFlag = "bor-chain-id"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "heimdall-bridge",
	Aliases: []string{"bridge"},
	Short:   "Heimdall bridge deamon",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// initialize tendermint viper config
		initTendermintViperConfig(cmd)
	},
}

// BridgeCommands returns command for bridge service
func BridgeCommands() *cobra.Command {
	return rootCmd
}

// initTendermintViperConfig sets global viper configuration needed to heimdall
func initTendermintViperConfig(cmd *cobra.Command) {
	heimdallNode, _ := cmd.Flags().GetString(helper.HeimdallNodeFlag)
	homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)
	withHeimdallConfigValue, _ := cmd.Flags().GetString(helper.WithHeimdallConfigFlag)
	bridgeDBValue, _ := cmd.Flags().GetString(bridgeDBFlag)
	borChainIDValue, _ := cmd.Flags().GetString(borChainIDFlag)

	// bridge-db directory (default storage)
	if bridgeDBValue == "" {
		bridgeDBValue = filepath.Join(homeValue, "bridge", "storage")
	}

	// set to viper
	viper.Set(helper.HeimdallNodeFlag, heimdallNode)
	viper.Set(helper.HomeFlag, homeValue)
	viper.Set(helper.WithHeimdallConfigFlag, withHeimdallConfigValue)
	viper.Set(bridgeDBFlag, bridgeDBValue)
	viper.Set(borChainIDFlag, borChainIDValue)

	// start heimdall config
	helper.InitHeimdallConfig("")
}

func init() {
	var logger = helper.Logger.With("module", "bridge/cmd/")
	rootCmd.PersistentFlags().StringP(helper.HeimdallNodeFlag, "n", helper.DefaultHeimdallNode, "Node to connect to")
	rootCmd.PersistentFlags().String(helper.HomeFlag, helper.DefaultNodeHome, "directory for config and data")
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

	// bind all flags with viper
	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		logger.Error("init | BindPFlag | rootCmd.Flags", "Error", err)
	}
}
