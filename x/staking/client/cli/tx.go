package cli

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/bridge/setu/util"
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
		SendValidatorJoinTx(),
		// NewStakeUpdateCmd(),
		// NewSignerUpdateCmd(),
		// NewValidatorExitCmd(),
	)

	return stakingTxCmd
}

// SendValidatorJoinTx send validator join transaction
func SendValidatorJoinTx(cdc *codec.Codec) *cobra.Command {
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
			proposer := hmTypes.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash is required")
			}

			pubkeyStr := viper.GetString(FlagSignerPubkey)
			if pubkeyStr == "" {
				return fmt.Errorf("pubkey is required")
			}

			pubkeyBytes := common.FromHex(pubkeyStr)
			if len(pubkeyBytes) != 65 {
				return fmt.Errorf("Invalid public key length")
			}
			pubkey := hmTypes.NewPubKey(pubkeyBytes)

			// total stake amount
			amount, ok := sdk.NewIntFromString(viper.GetString(FlagAmount))
			if !ok {
				return errors.New("Invalid stake amount")
			}

			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			chainmanagerParams, err := util.GetChainmanagerParams(clientCtx)
			if err != nil {
				return err
			}

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(txhash).EthHash(), chainmanagerParams.MainchainTxConfirmations)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			abiObject := &contractCallerObj.StakingInfoABI
			eventName := "Staked"
			event := new(stakinginfo.StakinginfoStaked)
			var logIndex uint64
			found := false
			for _, vLog := range receipt.Logs {
				topic := vLog.Topics[0].Bytes()
				selectedEvent := helper.EventByID(abiObject, topic)
				if selectedEvent != nil && selectedEvent.Name == eventName {
					if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
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

			// msg
			msg := types.NewMsgValidatorJoin(
				proposer,
				event.ValidatorId.Uint64(),
				viper.GetUint64(FlagActivationEpoch),
				amount,
				pubkey,
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				viper.GetUint64(FlagBlockNumber),
				event.Nonce.Uint64(),
			)

			// broadcast messages
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagSignerPubkey, "", "--signer-pubkey=<signer pubkey here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")
	cmd.Flags().String(FlagAmount, "0", "--amount=<amount>")
	cmd.Flags().Uint64(FlagActivationEpoch, 0, "--activation-epoch=<activation-epoch>")

	if err := cmd.MarkFlagRequired(FlagBlockNumber); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagBlockNumber", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagActivationEpoch); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagActivationEpoch", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagAmount); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagAmount", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagSignerPubkey); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagSignerPubkey", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		logger.Error("SendValidatorJoinTx | MarkFlagRequired | FlagTxHash", "Error", err)
	}
	return cmd
}
