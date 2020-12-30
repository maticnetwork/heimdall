package cli

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/common"
	ethcrypto "github.com/maticnetwork/bor/crypto"

	// "github.com/maticnetwork/heimdall/bridge/setu/util"

	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

// ForeignEventName is used in ValidatorJoinTxCmd
const ForeignEventName = "Staked"


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

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
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

			// get PubKey string
			pubkeyStr, _ := cmd.Flags().GetString(FlagSignerPubkey)
			if pubkeyStr == "" {
				return fmt.Errorf("pubkey is required")
			}

			// convert PubKey to bytes
			compressedPubkeyBytes := common.FromHex(pubkeyStr)

			ecdsaPubkey, err := ethcrypto.DecompressPubkey(compressedPubkeyBytes)
			if err != nil {
				return err
			}
			pubkeyBytes := ethcrypto.FromECDSAPub(ecdsaPubkey)

			if len(pubkeyBytes) != 65 {
				return fmt.Errorf("invalid public key length")
			}
			pubkey := hmTypes.NewPubKey(pubkeyBytes)

			// total stake amount
			amountStr, _ := cmd.Flags().GetString(FlagAmount)
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return fmt.Errorf("invalid stake amount")
			}

			// Get contractCaller ref
			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			// TODO uncomment this when integrating chainmanager
			// chainmanagerParams, err := util.GetChainmanagerParams(cliCtx)
			// if err != nil {
			// 	return err
			// }

			// get main tx receipt
			// NOTE: Use 'chainmanagerParams.MainchainTxConfirmations'. Now it is hard coded.
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(txhash).EthHash(), 6)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			abiObject := &contractCallerObj.StakingInfoABI
			event := new(stakinginfo.StakinginfoStaked)
			var logIndex uint64
			found := false
			for _, vLog := range receipt.Logs {
				topic := vLog.Topics[0].Bytes()
				selectedEvent := helper.EventByID(abiObject, topic)
				if selectedEvent != nil && selectedEvent.Name == ForeignEventName {
					if err := helper.UnpackLog(abiObject, event, ForeignEventName, vLog); err != nil {
						return err
					}

					logIndex = uint64(vLog.Index)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("Invalid tx for validator join")
			}

			if !bytes.Equal(event.SignerPubkey, pubkey.Bytes()[1:]) {
				return fmt.Errorf("Public key mismatch with event log")
			}

			activationEpoch, _ := cmd.Flags().GetUint64(FlagActivationEpoch)
			blockNumber, _ := cmd.Flags().GetUint64(FlagBlockNumber)

			// msg new ValidatorJion message
			msg := types.NewMsgValidatorJoin(
				proposer,
				event.ValidatorId.Uint64(),
				activationEpoch,
				amount,
				pubkey,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				blockNumber,
				event.Nonce.Uint64(),
			)

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagSignerPubkey, "", "--signer-pubkey=<signer pubkey here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().String(FlagAmount, "0", "--amount=<amount>")
	cmd.Flags().Uint64(FlagActivationEpoch, 0, "--activation-epoch=<activation-epoch>")

	_ = cmd.MarkFlagRequired(FlagBlockNumber)
	_ = cmd.MarkFlagRequired(FlagActivationEpoch)
	_ = cmd.MarkFlagRequired(FlagAmount)
	_ = cmd.MarkFlagRequired(FlagSignerPubkey)
	_ = cmd.MarkFlagRequired(FlagTxHash)

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

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)
			if err != nil {
				panic(err)
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// get validatorID from flags
			ValidatorID, _ := cmd.Flags().GetUint64(FlagValidatorID)
			if ValidatorID == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}

			// get PubKey string
			pubkeyStr, _ := cmd.Flags().GetString(FlagSignerPubkey)
			if pubkeyStr == "" {
				return fmt.Errorf("pubkey is required")
			}

			// convert PubKey to bytes
			compressedPubkeyBytes := common.FromHex(pubkeyStr)

			ecdsaPubkey, err := ethcrypto.DecompressPubkey(compressedPubkeyBytes)
			if err != nil {
				panic(err)
			}
			pubkeyBytes := ethcrypto.FromECDSAPub(ecdsaPubkey)

			if len(pubkeyBytes) != 65 {
				return fmt.Errorf("invalid public key length")
			}
			pubkey := hmTypes.NewPubKey(pubkeyBytes)

			// get txHash from flag
			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)
			blockNumber, _ := cmd.Flags().GetUint64(FlagBlockNumber)
			nonce, _ := cmd.Flags().GetUint64(FlagNonce)

			// draft new SingerUpdate message
			msg := types.NewMsgSignerUpdate(
				proposer,
				ValidatorID,
				pubkey,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				blockNumber,
				nonce,
			)

			// broadcast messages
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().Uint64(FlagValidatorID, 0, "--id=<validator-id>")
	cmd.Flags().String(FlagNewSignerPubkey, "", "--new-pubkey=<new-signer-pubkey>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().Int(FlagNonce, 0, "--nonce=<nonce>")

	cmd.MarkFlagRequired(FlagValidatorID)     //nolint
	cmd.MarkFlagRequired(FlagTxHash)          //nolint
	cmd.MarkFlagRequired(FlagNewSignerPubkey) //nolint
	cmd.MarkFlagRequired(FlagLogIndex)        //nolint
	cmd.MarkFlagRequired(FlagBlockNumber)     //nolint
	cmd.MarkFlagRequired(FlagNonce)           //nolint

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

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)

			if err != nil {
				panic(err)
			}

			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// get validatorID from flag
			validatorID, _ := cmd.Flags().GetUint64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}

			// get txHash from flag
			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			// total stake amount
			amountStr, _ := cmd.Flags().GetString(FlagAmount)
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return errors.New("Invalid new stake amount")
			}

			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)
			blockNumber, _ := cmd.Flags().GetUint64(FlagBlockNumber)
			nonce, _ := cmd.Flags().GetUint64(FlagNonce)

			// draft new StakeUpdate message
			msg := types.NewMsgStakeUpdate(
				proposer,
				validatorID,
				amount,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				blockNumber,
				nonce,
			)

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().Uint64(FlagValidatorID, 0, "--id=<validator-id>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().String(FlagAmount, "", "--amount=<amount>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().Int(FlagNonce, 0, "--nonce=<nonce>")

	cmd.MarkFlagRequired(FlagTxHash)      //nolint
	cmd.MarkFlagRequired(FlagLogIndex)    //nolint
	cmd.MarkFlagRequired(FlagValidatorID) //nolint
	cmd.MarkFlagRequired(FlagBlockNumber) //nolint
	cmd.MarkFlagRequired(FlagAmount)      //nolint
	cmd.MarkFlagRequired(FlagNonce)       //nolint

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

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)

			if err != nil {
				panic(err)
			}
			//proposer := sdk.AccAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// get validatorid from flag
			validatorID, _ := cmd.Flags().GetUint64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}

			// get txHash from flag
			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)
			blockNumber, _ := cmd.Flags().GetUint64(FlagBlockNumber)
			nonce, _ := cmd.Flags().GetUint64(FlagNonce)
			deactivationEpoch, _ := cmd.Flags().GetUint64(FlagDeactivationEpoch)

			// draf new ValidatorExit message
			msg := types.NewMsgValidatorExit(
				proposer,
				validatorID,
				deactivationEpoch,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				blockNumber,
				nonce,
			)

			// broadcast message
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().Uint64(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagDeactivationEpoch, 0, "--deactivation-epoch=<deactivation-epoch>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().Int(FlagNonce, 0, "--nonce=<nonce>")

	cmd.MarkFlagRequired(FlagValidatorID) //nolint
	cmd.MarkFlagRequired(FlagTxHash)      //nolint
	cmd.MarkFlagRequired(FlagLogIndex)    //nolint
	cmd.MarkFlagRequired(FlagBlockNumber) //nolint
	cmd.MarkFlagRequired(FlagNonce)       //nolint

	return cmd
}
