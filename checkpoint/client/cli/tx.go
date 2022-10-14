package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	types "github.com/maticnetwork/heimdall/checkpoint/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "checkpoint/client/cli")

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Checkpoint transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			SendCheckpointTx(cdc),
			SendCheckpointACKTx(cdc),
			SendCheckpointNoACKTx(cdc),
			SendCheckpointAdjust(cdc),
		)...,
	)
	return txCmd
}

// SendCheckpointAdjust adjusts previous checkpoint transaction
func SendCheckpointAdjust(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-adjust",
		Short: "adjusts previous checkpoint transaction according to ethereum chain (details to be provided for checkpoint on ethereum chain)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// header index
			headerIndexStr := viper.GetString(FlagHeaderNumber)
			if headerIndexStr == "" {
				return fmt.Errorf("header number cannot be empty")
			}
			headerIndex, err := strconv.ParseUint(headerIndexStr, 10, 64)
			if err != nil {
				return err
			}

			// get from
			from := helper.GetFromAddress(cliCtx)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				return fmt.Errorf("proposer cannot be empty")
			}

			//	start block
			startBlockStr := viper.GetString(FlagStartBlock)
			if startBlockStr == "" {
				return fmt.Errorf("start block cannot be empty")
			}
			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			//	end block
			endBlockStr := viper.GetString(FlagEndBlock)
			if endBlockStr == "" {
				return fmt.Errorf("end block cannot be empty")
			}
			endBlock, err := strconv.ParseUint(endBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// root hash
			rootHashStr := viper.GetString(FlagRootHash)
			if rootHashStr == "" {
				return fmt.Errorf("root hash cannot be empty")
			}

			msg := types.NewMsgCheckpointAdjust(
				headerIndex,
				startBlock,
				endBlock,
				proposer,
				from,
				hmTypes.HexToHeimdallHash(rootHashStr),
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(FlagHeaderNumber, "", "--header=<header-index>")
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.Flags().String(FlagEndBlock, "", "--end-block=<end-block-number>")
	cmd.Flags().StringP(FlagRootHash, "r", "", "--root-hash=<root-hash>")

	cmd.MarkFlagRequired(FlagHeaderNumber)
	cmd.MarkFlagRequired(FlagRootHash)
	cmd.MarkFlagRequired(FlagProposerAddress)
	cmd.MarkFlagRequired(FlagStartBlock)
	cmd.MarkFlagRequired(FlagEndBlock)

	return cmd
}

// SendCheckpointTx send checkpoint transaction
func SendCheckpointTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// bor chain id
			borChainID := viper.GetString(FlagBorChainID)
			if borChainID == "" {
				return fmt.Errorf("bor chain id cannot be empty")
			}

			if viper.GetBool(FlagAutoConfigure) {
				var checkpointProposer hmTypes.Validator
				proposerBytes, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", types.StakingQuerierRoute, types.QueryCurrentProposer))
				if err != nil {
					return err
				}

				if err := json.Unmarshal(proposerBytes, &checkpointProposer); err != nil {
					return err
				}

				if !bytes.Equal(checkpointProposer.Signer.Bytes(), helper.GetAddress()) {
					return fmt.Errorf("Please wait for your turn to propose checkpoint. Checkpoint proposer:%v", checkpointProposer.String())
				}

				// create bor chain id params
				borChainIDParams := types.NewQueryBorChainID(borChainID)
				bz, err := cliCtx.Codec.MarshalJSON(borChainIDParams)
				if err != nil {
					return err
				}

				// fetch msg checkpoint
				result, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextCheckpoint), bz)
				if err != nil {
					return err
				}

				// unmarsall the checkpoint msg
				var newCheckpointMsg types.MsgCheckpoint
				if err := json.Unmarshal(result, &newCheckpointMsg); err != nil {
					return err
				}

				// broadcast this checkpoint
				return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{newCheckpointMsg})
			}

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			//	start block
			startBlockStr := viper.GetString(FlagStartBlock)
			if startBlockStr == "" {
				return fmt.Errorf("start block cannot be empty")
			}
			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			//	end block
			endBlockStr := viper.GetString(FlagEndBlock)
			if endBlockStr == "" {
				return fmt.Errorf("end block cannot be empty")
			}
			endBlock, err := strconv.ParseUint(endBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// root hash
			rootHashStr := viper.GetString(FlagRootHash)
			if rootHashStr == "" {
				return fmt.Errorf("root hash cannot be empty")
			}

			// Account Root Hash
			accountRootHashStr := viper.GetString(FlagAccountRootHash)
			if accountRootHashStr == "" {
				return fmt.Errorf("account root hash cannot be empty")
			}

			msg := types.NewMsgCheckpointBlock(
				proposer,
				startBlock,
				endBlock,
				hmTypes.HexToHeimdallHash(rootHashStr),
				hmTypes.HexToHeimdallHash(accountRootHashStr),
				borChainID,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.Flags().String(FlagEndBlock, "", "--end-block=<end-block-number>")
	cmd.Flags().StringP(FlagRootHash, "r", "", "--root-hash=<root-hash>")
	cmd.Flags().String(FlagAccountRootHash, "", "--account-root=<account-root>")
	cmd.Flags().String(FlagBorChainID, "", "--bor-chain-id=<bor-chain-id>")
	cmd.Flags().Bool(FlagAutoConfigure, false, "--auto-configure=true/false")

	cmd.MarkFlagRequired(FlagRootHash)
	cmd.MarkFlagRequired(FlagAccountRootHash)
	cmd.MarkFlagRequired(FlagBorChainID)

	return cmd
}

// SendCheckpointACKTx send checkpoint ack transaction
func SendCheckpointACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-ack",
		Short: "send acknowledgement for checkpoint in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			headerBlockStr := viper.GetString(FlagHeaderNumber)
			if headerBlockStr == "" {
				return fmt.Errorf("header number cannot be empty")
			}

			headerBlock, err := strconv.ParseUint(headerBlockStr, 10, 64)
			if err != nil {
				return err
			}

			txHashStr := viper.GetString(FlagCheckpointTxHash)
			if txHashStr == "" {
				return fmt.Errorf("checkpoint tx hash cannot be empty")
			}

			txHash := hmTypes.BytesToHeimdallHash(common.FromHex(txHashStr))

			//
			// Get header details
			//

			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			chainmanagerParams, err := util.GetChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(txHash.EthHash(), chainmanagerParams.MainchainTxConfirmations)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			// decode new header block event
			res, err := contractCallerObj.DecodeNewHeaderBlockEvent(
				chainmanagerParams.ChainParams.RootChainAddress.EthAddress(),
				receipt,
				uint64(viper.GetInt64(FlagCheckpointLogIndex)),
			)
			if err != nil {
				return errors.New("Invalid transaction for header block")
			}

			// draft new checkpoint no-ack msg
			msg := types.NewMsgCheckpointAck(
				proposer, // ack tx sender
				headerBlock,
				hmTypes.BytesToHeimdallAddress(res.Proposer.Bytes()),
				res.Start.Uint64(),
				res.End.Uint64(),
				res.Root,
				txHash,
				uint64(viper.GetInt64(FlagCheckpointLogIndex)),
			)

			// msg
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagHeaderNumber, "", "--header=<header-index>")
	cmd.Flags().StringP(FlagCheckpointTxHash, "t", "", "--txhash=<checkpoint-txhash>")
	cmd.Flags().String(FlagCheckpointLogIndex, "", "--log-index=<log-index>")

	if err := cmd.MarkFlagRequired(FlagHeaderNumber); err != nil {
		logger.Error("SendCheckpointACKTx | MarkFlagRequired | FlagHeaderNumber", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagCheckpointTxHash); err != nil {
		logger.Error("SendCheckpointACKTx | MarkFlagRequired | FlagCheckpointTxHash", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagCheckpointLogIndex); err != nil {
		logger.Error("SendCheckpointACKTx | MarkFlagRequired | FlagCheckpointLogIndex", "Error", err)
	}

	return cmd
}

// SendCheckpointNoACKTx send no-ack transaction
func SendCheckpointNoACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-noack",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// create new checkpoint no-ack
			msg := types.NewMsgCheckpointNoAck(
				proposer,
			)

			// broadcast messages
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	return cmd
}
