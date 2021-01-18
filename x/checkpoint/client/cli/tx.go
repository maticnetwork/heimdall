package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
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
		// CheckpointACKTxCmd(),
		CheckpointNoACKTxCmd(),
	)

	return checkpointTxCmd
}

// CheckpointTxCmd send validator join message
func CheckpointTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// bor chain id
			borChainID, err := cmd.Flags().GetString(FlagBorChainID)
			if err != nil {
				return err
			}

			if borChainID == "" {
				return fmt.Errorf("bor chain id cannot be empty")
			}

			isAutoConfigure, err := cmd.Flags().GetBool(FlagAutoConfigure)
			if err != nil {
				return err
			}

			if isAutoConfigure {
				var checkpointProposer hmTypes.Validator
				proposerBytes, _, err := clientCtx.Query(fmt.Sprintf("custom/%s/%s", types.StakingQuerierRoute, types.QueryCurrentProposer))
				if err != nil {
					return err
				}

				if err := json.Unmarshal(proposerBytes, &checkpointProposer); err != nil {
					return err
				}

				if !bytes.Equal([]byte(checkpointProposer.Signer), helper.GetAddress()) {
					return fmt.Errorf("Please wait for your turn to propose checkpoint. Checkpoint proposer:%v", checkpointProposer.String())
				}

				// create bor chain id params
				borChainIDParams := types.NewQueryBorChainID(borChainID)
				bz, err := clientCtx.JSONMarshaler.MarshalJSON(&borChainIDParams)
				if err != nil {
					return err
				}

				// fetch msg checkpoint
				result, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextCheckpoint), bz)
				if err != nil {
					return err
				}

				// unmarsall the checkpoint msg
				var newCheckpointMsg types.MsgCheckpoint
				if err := json.Unmarshal(result, &newCheckpointMsg); err != nil {
					return err
				}

				// broadcast this checkpoint
				return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &newCheckpointMsg)
			}

			// get proposer
			proposerAddressStr, err := cmd.Flags().GetString(FlagProposerAddress)
			if err != nil {
				return err
			}

			proposer := sdk.AccAddress([]byte(proposerAddressStr))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			//	start block

			startBlockStr, err := cmd.Flags().GetString(FlagStartBlock)
			if err != nil {
				return err
			}
			if startBlockStr == "" {
				return fmt.Errorf("start block cannot be empty")
			}

			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			//	end block

			endBlockStr, err := cmd.Flags().GetString(FlagEndBlock)
			if err != nil {
				return err
			}
			if endBlockStr == "" {
				return fmt.Errorf("end block cannot be empty")
			}

			endBlock, err := strconv.ParseUint(endBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// root hash

			rootHashStr, err := cmd.Flags().GetString(FlagRootHash)
			if err != nil {
				return err
			}
			if rootHashStr == "" {
				return fmt.Errorf("root hash cannot be empty")
			}

			// Account Root Hash
			accountRootHashStr, err := cmd.Flags().GetString(FlagAccountRootHash)
			if err != nil {
				return err
			}
			if accountRootHashStr == "" {
				return fmt.Errorf("account root hash cannot be empty")
			}

			msg := types.NewMsgCheckpointBlock(
				proposer,
				startBlock,
				endBlock,
				hmCommonTypes.HexToHeimdallHash(rootHashStr),
				hmCommonTypes.HexToHeimdallHash(accountRootHashStr),
				borChainID,
			)

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

// // CheckpointACKTxCmd send checkpoint ack transaction
// func CheckpointACKTxCmd() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "send-ack",
// 		Short: "send acknowledgement for checkpoint in buffer",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx := client.GetClientContextFromCmd(cmd)
// 			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
// 			if err != nil {
// 				return err
// 			}

// 			// get proposer
// 			proposer := sdk.AccAddress([]byte(cmd.Flags().GetString(FlagProposerAddress)))
// 			if proposer.Empty() {
// 				proposer = helper.GetFromAddress(clientCtx)
// 			}

// 			headerBlockStr := cmd.Flags().GetString(FlagHeaderNumber)
// 			if headerBlockStr == "" {
// 				return fmt.Errorf("header number cannot be empty")
// 			}

// 			headerBlock, err := strconv.ParseUint(headerBlockStr, 10, 64)
// 			if err != nil {
// 				return err
// 			}

// 			txHashStr := cmd.Flags().GetString(FlagCheckpointTxHash)
// 			if txHashStr == "" {
// 				return fmt.Errorf("checkpoint tx hash cannot be empty")
// 			}

// 			txHash := hmCommonTypes.BytesToHeimdallHash(common.FromHex(txHashStr))

// 			//
// 			// Get header details
// 			//

// 			contractCallerObj, err := helper.NewContractCaller()
// 			if err != nil {
// 				return err
// 			}

// 			chainmanagerParams, err := util.GetChainmanagerParams(cliCtx)
// 			if err != nil {
// 				return err
// 			}

// 			// get main tx receipt
// 			receipt, err := contractCallerObj.GetConfirmedTxReceipt(txHash.EthHash(), chainmanagerParams.MainchainTxConfirmations)
// 			if err != nil || receipt == nil {
// 				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
// 			}

// 			// decode new header block event
// 			res, err := contractCallerObj.DecodeNewHeaderBlockEvent(
// 				chainmanagerParams.ChainParams.RootChainAddress.EthAddress(),
// 				receipt,
// 				uint64(cmd.Flags().GetInt64(FlagCheckpointLogIndex)),
// 			)
// 			if err != nil {
// 				return errors.New("Invalid transaction for header block")
// 			}

// 			// draft new checkpoint no-ack msg
// 			msg := types.NewMsgCheckpointAck(
// 				proposer, // ack tx sender
// 				headerBlock,
// 				sdk.AccAddress(res.Proposer.Bytes()),
// 				res.Start.Uint64(),
// 				res.End.Uint64(),
// 				res.Root,
// 				txHash,
// 				uint64(cmd.Flags().GetInt64(FlagCheckpointLogIndex)),
// 			)

// 			// broadcast messages
// 			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)

// 		},
// 	}

// 	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
// 	cmd.Flags().String(FlagHeaderNumber, "", "--header=<header-index>")
// 	cmd.Flags().StringP(FlagCheckpointTxHash, "t", "", "--txhash=<checkpoint-txhash>")
// 	cmd.Flags().String(FlagCheckpointLogIndex, "", "--log-index=<log-index>")

// 	_ = cmd.MarkFlagRequired(FlagHeaderNumber)
// 	_ = cmd.MarkFlagRequired(FlagCheckpointTxHash)
// 	_ = cmd.MarkFlagRequired(FlagCheckpointLogIndex)

// 	flags.AddTxFlagsToCmd(cmd)

// 	return cmd
// }

// CheckpointNoACKTxCmd send no-ack transaction
func CheckpointNoACKTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-noack",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// get proposer

			proposerAddressStr, err := cmd.Flags().GetString(FlagProposerAddress)
			if err != nil {
				return err
			}

			proposer := sdk.AccAddress([]byte(proposerAddressStr))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}
			// create new checkpoint no-ack
			msg := types.NewMsgCheckpointNoAck(
				proposer,
			)

			// broadcast messages
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)

		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
