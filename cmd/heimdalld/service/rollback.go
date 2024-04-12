package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
)

func rollbackCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback cosmos-sdk and tendermint state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when Tendermint has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also roll back to height n - 1. No blocks are removed, so upon
restarting Tendermint the transactions in block n will be re-executed against the
application.
`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			db, err := sdk.NewLevelDB("application", config.DBDir())
			if err != nil {
				return err
			}

			height, hash, err := commands.RollbackState(config)

			if err != nil {
				return fmt.Errorf("failed to rollback tendermint state: %w", err)
			}
			// rollback the multistore
			cms := rootmulti.NewStore(db)
			cms.RollbackToVersion(height)
			fmt.Printf("Rolled back state to height %d and hash %X", height, hash)
			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "Node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "Client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "Genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int(stakingcli.FlagValidatorID, 1, "--id=<validator ID here>, if left blank will be assigned 1")

	return cmd
}
