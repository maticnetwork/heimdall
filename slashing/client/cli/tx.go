package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "staking/client/cli")

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	slashingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Slashing transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	slashingTxCmd.AddCommand(flags.PostCommands(
		GetCmdUnjail(cdc),
		GetCmdTick(cdc),
		GetCmdTickAck(cdc),
	)...)

	return slashingTxCmd
}

func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed",
		Long: `unjail a jailed validator:

$ <appcli> tx slashing unjail --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			validator := viper.GetInt64(FlagValidatorID)
			if validator == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}

			txHash := viper.GetString(FlagTxHash)
			if txHash == "" {
				return fmt.Errorf("transaction hash is required")
			}

			msg := types.NewMsgUnjail(
				proposer,
				uint64(validator),
				hmTypes.HexToHeimdallHash(txHash),
				uint64(viper.GetInt64(FlagLogIndex)),
				viper.GetUint64(FlagBlockNumber),
			)

			// broadcast messages
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.MarkFlagRequired(FlagTxHash)
	if err := cmd.MarkFlagRequired(FlagBlockNumber); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagBlockNumber", "Error", err)
	}
	return cmd
}

func GetCmdTick(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "tick",
		Short: "send slash tick when total slashedamount exceeds limit",
		Long:  "<appcli>",

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			slashInfoBytes := viper.GetString(FlagSlashInfoBytes)
			if slashInfoBytes == "" {
				return fmt.Errorf("slashinfo bytes has to be supplied")
			}

			msg := types.NewMsgTick(
				viper.GetUint64(FlagTickID),
				proposer,
				hmTypes.HexToHexBytes(slashInfoBytes),
			)

			// braodcast messages
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagSlashInfoBytes, "", "--slashinfo-bytes=<slashinfo-bytes>")
	cmd.Flags().Uint64(FlagTickID, 1, "--tick-id=<tick-id>")

	cmd.MarkFlagRequired(FlagSlashInfoBytes)
	cmd.MarkFlagRequired(FlagTickID)

	return cmd
}

func GetCmdTickAck(cdc *codec.Codec) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "tick-ack",
		Short: "send tick ack",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			txHash := viper.GetString(FlagTxHash)
			if txHash == "" {
				return fmt.Errorf("transaction hash is required")
			}

			msg := types.NewMsgTickAck(
				proposer,
				viper.GetUint64(FlagTickID),
				viper.GetUint64(FlagAmount),
				hmTypes.HexToHeimdallHash(txHash),
				uint64(viper.GetInt64(FlagLogIndex)),
				viper.GetUint64(FlagBlockNumber),
			)

			// broadcast messages
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagAmount, 0, "--amount=<amount>")
	cmd.Flags().Uint64(FlagTickID, 1, "--tick-id=<tick-id>")

	if err := cmd.MarkFlagRequired(FlagBlockNumber); err != nil {
		logger.Error("SendTickAckTx | MarkFlagRequired | FlagBlockNumber", "Error", err)
	}
	cmd.MarkFlagRequired(FlagTxHash)
	cmd.MarkFlagRequired(FlagLogIndex)
	cmd.MarkFlagRequired(FlagAmount)
	cmd.MarkFlagRequired(FlagTickID)

	return cmd
}
