package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/helper"
	logger "github.com/tendermint/tendermint/libs/log"
)

const (
	bridgeDBFlag   = "bridge-db"
	borChainIDFlag = "bor-chain-id"
	logsTypeFlag   = "logs-type"
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
func BridgeCommands(v *viper.Viper, loggerInstance logger.Logger, caller string) *cobra.Command {
	DecorateWithBridgeRootFlags(rootCmd, v, loggerInstance, caller)
	return rootCmd
}

// function is called when bridge flags needs to be added to command
func DecorateWithBridgeRootFlags(cmd *cobra.Command, v *viper.Viper, loggerInstance logger.Logger, caller string) {
	cmd.PersistentFlags().StringP(helper.TendermintNodeFlag, "n", helper.DefaultTendermintNode, "Node to connect to")
	if err := v.BindPFlag(helper.TendermintNodeFlag, cmd.PersistentFlags().Lookup(helper.TendermintNodeFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, helper.TendermintNodeFlag), "Error", err)
	}

	cmd.PersistentFlags().String(helper.HomeFlag, helper.DefaultNodeHome, "directory for config and data")
	if err := v.BindPFlag(helper.HomeFlag, cmd.PersistentFlags().Lookup(helper.HomeFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, helper.HomeFlag), "Error", err)
	}

	// bridge storage db
	cmd.PersistentFlags().String(
		bridgeDBFlag,
		"",
		"Bridge db path (default <home>/bridge/storage)",
	)
	if err := v.BindPFlag(bridgeDBFlag, cmd.PersistentFlags().Lookup(bridgeDBFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, bridgeDBFlag), "Error", err)
	}

	// bridge chain id
	cmd.PersistentFlags().String(
		borChainIDFlag,
		helper.DefaultBorChainID,
		"Bor chain id",
	)

	// bridge logging type
	cmd.PersistentFlags().String(
		logsTypeFlag,
		helper.DefaultLogsType,
		"Use json logger",
	)

	if err := v.BindPFlag(borChainIDFlag, cmd.PersistentFlags().Lookup(borChainIDFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, borChainIDFlag), "Error", err)
	}
}

// function is called to set appropriate bridge db path
func AdjustBridgeDBValue(cmd *cobra.Command, v *viper.Viper) {
	bridgeDBValue, _ := cmd.Flags().GetString(bridgeDBFlag)
	homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)

	// bridge-db directory (default storage)
	if bridgeDBValue == "" {
		bridgeDBValue = filepath.Join(homeValue, "bridge", "storage")
	}

	v.Set(bridgeDBFlag, bridgeDBValue)
}

// initTendermintViperConfig sets global viper configuration needed to heimdall
func initTendermintViperConfig(cmd *cobra.Command) {

	logsTypeValue, _ := cmd.Flags().GetString(logsTypeFlag)
	viper.Set(logsTypeFlag, logsTypeValue)

	// set appropriate bridge DB
	AdjustBridgeDBValue(cmd, viper.GetViper())

	// start heimdall config
	helper.InitHeimdallConfig("")
}
