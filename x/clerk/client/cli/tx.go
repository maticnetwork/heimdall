package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmTypesCommon "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/clerk/types"
)

// NewTxCmd returns the transaction commands for this module
// governance ModuleClient is slightly different from other ModuleClients in that
// it contains a slice of "proposal" child commands. These commands are respective
// to proposal type handlers that are implemented in other modules but are mounted
// under the governance CLI (eg. parameter change proposals).
func NewTxCmd() *cobra.Command {
	clerkTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	clerkTxCmd.AddCommand(
		CreateNewStateRecord(),
	)

	return clerkTxCmd
}

// CreateNewStateRecord send checkpoint transaction
func CreateNewStateRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "new state record",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// bor chain id
			borChainID, err := cmd.Flags().GetString(FlagBorChainId)
			if err != nil {
				return err
			}
			if borChainID == "" {
				return fmt.Errorf("BorChainID cannot be empty")
			}

			// get proposer
			proposerCmdStr, err := cmd.Flags().GetString(FlagProposerAddress)
			if err != nil {
				return err
			}
			proposerCmdStr = strings.ToLower(proposerCmdStr)
			proposer, err := sdk.AccAddressFromHex(proposerCmdStr)
			if err != nil {
				return err
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// tx hash
			txHashStr, err := cmd.Flags().GetString(FlagTxHash)
			if err != nil {
				return err
			}
			if txHashStr == "" {
				return fmt.Errorf("tx hash cannot be empty")
			}

			// tx hash
			recordIDStr, err := cmd.Flags().GetString(FlagRecordID)
			if err != nil {
				return err
			}
			if recordIDStr == "" {
				return fmt.Errorf("record id cannot be empty")
			}

			recordID, err := strconv.ParseUint(recordIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("record id cannot be empty")
			}

			// get contract Addr
			contractAddrCmdStr, err := cmd.Flags().GetString(FlagContractAddress)
			if err != nil {
				return err
			}
			contractAddrCmdStr = strings.ToLower(contractAddrCmdStr)
			contractAddr, err := sdk.AccAddressFromHex(contractAddrCmdStr)
			if err != nil {
				return err
			}
			if contractAddr.Empty() {
				return fmt.Errorf("contract Address cannot be empty")
			}

			// log index
			logIndexStr, err := cmd.Flags().GetString(FlagLogIndex)
			if err != nil {
				return err
			}
			if logIndexStr == "" {
				return fmt.Errorf("log index cannot be empty")
			}

			logIndex, err := strconv.ParseUint(logIndexStr, 10, 64)
			if err != nil {
				return fmt.Errorf("log index cannot be parsed")
			}

			// log index
			dataStr, err := cmd.Flags().GetString(FlagData)
			if err != nil {
				return err
			}
			if dataStr == "" {
				return fmt.Errorf("data cannot be empty")
			}

			data := hmTypes.HexToHexBytes(dataStr)
			if dataStr == "" {
				return fmt.Errorf("data should be hex string")
			}

			flagBlockNumber, err := cmd.Flags().GetUint64(FlagBlockNumber)
			if err != nil {
				return err
			}

			// create new state record
			msg := types.NewMsgEventRecord(
				proposer,
				hmTypesCommon.HexToHeimdallHash(txHashStr),
				logIndex,
				flagBlockNumber,
				recordID,
				contractAddr,
				data,
				borChainID,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<tx-hash>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.Flags().String(FlagRecordID, "", "--id=<record-id>")
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=<bor-chain-id>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().String(FlagContractAddress, "", "--contract-addr=<contract-addr>")
	cmd.Flags().String(FlagData, "", "--data=<data>")

	return cmd
}
