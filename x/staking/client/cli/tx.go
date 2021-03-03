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

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	chainmanagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	stakingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingTxCmd.AddCommand(
		ValidatorJoinTxCmd(),
		StakeUpdateTxCmd(),
		SignerUpdateTxCmd(),
		ValidatorExitTxCmd(),
	)

	return stakingTxCmd
}

// validateAndCompressPubKey validate and compress the pubkey
func validateAndCompressPubKey(pubkeyBytes []byte) ([]byte, error) {
	if len(pubkeyBytes) == helper.UNCOMPRESSED_PUBKEY_SIZE {
		pubkeyBytes = helper.AppendPubkeyPrefix(pubkeyBytes)
	}

	// check if key is uncompressed
	if len(pubkeyBytes) == helper.UNCOMPRESSED_PUBKEY_SIZE_WITH_PREFIX {
		var err error
		pubkeyBytes, err = helper.CompressPubKey(pubkeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Invalid uncompressed pubkey %s", err)
		}
	}

	if len(pubkeyBytes) != helper.COMPRESSED_PUBKEY_SIZE_WITH_PREFIX {
		return nil, fmt.Errorf("Invalid compressed pubkey")
	}

	return pubkeyBytes, nil
}

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

//
//
//

// ValidatorJoinTxCmd send validator join message
func ValidatorJoinTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-join",
		Short: "Join Heimdall as a validator",
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

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}
			stakingManagerAddress, _ := sdk.AccAddressFromHex(chainmanagerParams.ChainParams.StakingManagerAddress)
			event, err := contractCallerObj.DecodeValidatorJoinEvent(
				stakingManagerAddress,
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// convert PubKey to bytes
			pubkeyBytes, err := validateAndCompressPubKey(event.SignerPubkey)
			if err != nil {
				return fmt.Errorf("Invalid uncompressed pubkey %s", err)
			}
			// create new pub key
			pubkey := hmTypes.NewPubKey(pubkeyBytes)

			// msg new ValidatorJion message
			msg := types.NewMsgValidatorJoin(
				proposer,
				event.ValidatorId.Uint64(),
				event.ActivationEpoch.Uint64(),
				sdk.NewIntFromBigInt(event.Amount),
				pubkey,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
				event.Nonce.Uint64(),
			)

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)

	// add common tx flags to cmd
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// SignerUpdateTxCmd send singer update transaction
func SignerUpdateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-update",
		Short: "Update signer for a validator",
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

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			event, err := contractCallerObj.DecodeSignerUpdateEvent(
				common.FromHex(chainmanagerParams.ChainParams.StakingManagerAddress),
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// convert PubKey to bytes
			pubkeyBytes, err := validateAndCompressPubKey(event.SignerPubkey)
			if err != nil {
				return fmt.Errorf("Invalid uncompressed pubkey %s", err)
			}
			// create new pub key
			pubkey := hmTypes.NewPubKey(pubkeyBytes)

			// draft new SingerUpdate message
			msg := types.NewMsgSignerUpdate(
				proposer,
				event.ValidatorId.Uint64(),
				pubkey,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
				event.Nonce.Uint64(),
			)

			// broadcast messages
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)

	// add common tx flags to cmd
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// StakeUpdateTxCmd send stake update transaction
func StakeUpdateTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake-update",
		Short: "Update stake for a validator",
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

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			event, err := contractCallerObj.DecodeValidatorStakeUpdateEvent(
				common.FromHex(chainmanagerParams.ChainParams.StakingManagerAddress),
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// draft new StakeUpdate message
			msg := types.NewMsgStakeUpdate(
				proposer,
				event.ValidatorId.Uint64(),
				sdk.NewIntFromBigInt(event.NewAmount),
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
				event.Nonce.Uint64(),
			)
			if err != nil {
				return err
			}

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)

	// add common tx flags to cmd
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// ValidatorExitTxCmd sends validator exit transaction
func ValidatorExitTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-exit",
		Short: "Exit heimdall as a validator ",
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

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			event, err := contractCallerObj.DecodeValidatorExitEvent(
				common.FromHex(chainmanagerParams.ChainParams.StakingManagerAddress),
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// draft new ValidatorExit message
			msg := types.NewMsgValidatorExit(
				proposer,
				event.ValidatorId.Uint64(),
				event.DeactivationEpoch.Uint64(),
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
				event.Nonce.Uint64(),
			)

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)

	// add common tx flags to cmd
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
