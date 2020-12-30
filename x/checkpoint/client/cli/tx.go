package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	checkpointTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	checkpointTxCmd.AddCommand(
		CheckpointTxCmd(),
	)

	return checkpointTxCmd
}

// CheckpointTxCmd send validator join message
func CheckpointTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			// clientCtx := client.GetClientContextFromCmd(cmd)
			// clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())

			// // bor chain id
			// borChainID := viper.GetString(FlagBorChainID)
			// if borChainID == "" {
			// 	return fmt.Errorf("bor chain id cannot be empty")
			// }

			// if viper.GetBool(FlagAutoConfigure) {
			// 	var checkpointProposer hmTypes.Validator
			// 	proposerBytes, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.StakingQuerierRoute, types.QueryCurrentProposer))
			// 	if err != nil {
			// 		return err
			// 	}

			// 	if err := json.Unmarshal(proposerBytes, &checkpointProposer); err != nil {
			// 		return err
			// 	}

			// 	if !bytes.Equal([]byte(checkpointProposer.Signer), helper.GetAddress()) {
			// 		return fmt.Errorf("Please wait for your turn to propose checkpoint. Checkpoint proposer:%v", checkpointProposer.String())
			// 	}

			// 	// create bor chain id params
			// 	borChainIDParams := types.NewQueryBorChainID(borChainID)
			// 	bz, err := clientCtx.cdc.MarshalJSON(borChainIDParams)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	// fetch msg checkpoint
			// 	result, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextCheckpoint), bz)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	// unmarsall the checkpoint msg
			// 	var newCheckpointMsg types.MsgCheckpoint
			// 	if err := json.Unmarshal(result, &newCheckpointMsg); err != nil {
			// 		return err
			// 	}

			// 	// broadcast this checkpoint
			// 	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &newCheckpointMsg)
			// }

			// // get proposer
			// proposer := sdk.AccAddress([]byte(viper.GetString(FlagProposerAddress)))
			// if proposer.Empty() {
			// 	proposer = helper.GetFromAddress(clientCtx)
			// }

			// //	start block

			// startBlockStr := viper.GetString(FlagStartBlock)
			// if startBlockStr == "" {
			// 	return fmt.Errorf("start block cannot be empty")
			// }

			// startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			// if err != nil {
			// 	return err
			// }

			// //	end block

			// endBlockStr := viper.GetString(FlagEndBlock)
			// if endBlockStr == "" {
			// 	return fmt.Errorf("end block cannot be empty")
			// }

			// endBlock, err := strconv.ParseUint(endBlockStr, 10, 64)
			// if err != nil {
			// 	return err
			// }

			// // root hash

			// rootHashStr := viper.GetString(FlagRootHash)
			// if rootHashStr == "" {
			// 	return fmt.Errorf("root hash cannot be empty")
			// }

			// // Account Root Hash
			// accountRootHashStr := viper.GetString(FlagAccountRootHash)
			// if accountRootHashStr == "" {
			// 	return fmt.Errorf("account root hash cannot be empty")
			// }

			// msg := types.NewMsgCheckpointBlock(
			// 	proposer,
			// 	startBlock,
			// 	endBlock,
			// 	hmcommonTypes.HexToHeimdallHash(rootHashStr),
			// 	hmcommonTypes.HexToHeimdallHash(accountRootHashStr),
			// 	borChainID,
			// )

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.Flags().String(FlagEndBlock, "", "--end-block=<end-block-number>")
	cmd.Flags().StringP(FlagRootHash, "r", "", "--root-hash=<root-hash>")
	cmd.Flags().String(FlagAccountRootHash, "", "--account-root=<account-root>")
	cmd.Flags().String(FlagBorChainID, "", "--bor-chain-id=<bor-chain-id>")
	cmd.Flags().Bool(FlagAutoConfigure, false, "--auto-configure=true/false")

	_ = cmd.MarkFlagRequired(FlagRootHash)
	_ = cmd.MarkFlagRequired(FlagAccountRootHash)
	_ = cmd.MarkFlagRequired(FlagBorChainID)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
