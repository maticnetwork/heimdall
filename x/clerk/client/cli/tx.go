package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	chainmanagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
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
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			chainmanagerParams, err := getChainmanagerParams(clientCtx)
			if err != nil {
				return err
			}

			// Get contractCaller ref
			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				return err
			}
			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposerAddrStr = strings.ToLower(proposerAddrStr)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)

			if err != nil {
				return fmt.Errorf("invalid proposer address: %v", err)
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// get txHash
			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash is required")
			}

			// parse log index
			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)

			// bor chain id
			borChainID, err := cmd.Flags().GetString(FlagBorChainId)
			if err != nil {
				return err
			}
			if borChainID == "" {
				return fmt.Errorf("BorChainID cannot be empty")
			}

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}
			stateSenderAddress, _ := sdk.AccAddressFromHex(chainmanagerParams.ChainParams.StateSenderAddress)
			event, err := contractCallerObj.DecodeStateSyncedEvent(
				stateSenderAddress,
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// create new state record
			msg := types.NewMsgEventRecord(
				proposer,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
				event.Id.Uint64(),
				sdk.AccAddress(event.ContractAddress.Bytes()),
				event.Data,
				borChainID,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<tx-hash>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=<bor-chain-id>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)
	_ = cmd.MarkFlagRequired(FlagBorChainId)

	// add common tx flags to cmd
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

//
// Get chainmanager params
//

// Fetch chain manager params
func getChainmanagerParams(clientCtx client.Context) (*chainmanagerTypes.Params, error) {
	// create query client
	queryClient := chainmanagerTypes.NewQueryClient(clientCtx)
	req := &chainmanagerTypes.QueryParamsRequest{}
	res, err := queryClient.Params(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.GetParams(), nil
}
