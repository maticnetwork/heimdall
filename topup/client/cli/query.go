package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/client"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/topup/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group topup queries under a subcommand
	topupQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the topup module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// topup query command
	topupQueryCmd.AddCommand(
		GetSequence(),
	)

	return topupQueryCmd
}

// GetSequence validator information via id or address
func GetSequence() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get sequence from txhash and logindex",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			logIndex := viper.GetUint64(FlagLogIndex)
			txHashStr := viper.GetString(FlagTxHash)
			if txHashStr == "" {
				return fmt.Errorf("LogIndex and transaction hash required")
			}

			var queryParams []byte
			var err error
			var t string = ""
			if txHashStr != "" {
				queryParams, err = clientCtx.Codec.MarshalJSON(types.NewQuerySequenceParams(txHashStr, logIndex))
				if err != nil {
					return err
				}
				t = types.QuerySequence
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, t), queryParams)
			if err != nil {
				fmt.Println("No topup exists")
				return nil
			}

			fmt.Println("Success. Topup exists with sequence:", string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		cliLogger.Error("GetSequence | MarkFlagRequired | FlagTxHash", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		cliLogger.Error("GetSequence | MarkFlagRequired | FlagLogIndex", "Error", err)
	}
	return cmd
}
