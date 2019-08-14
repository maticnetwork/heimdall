package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        borTypes.ModuleName,
		Short:                      "Bor transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			PostSendProposeSpanTx(cdc),
		)...,
	)
	return txCmd
}

// PostSendProposeSpanTx send propose span transaction
func PostSendProposeSpanTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-span",
		Short: "send propose span tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainID := viper.GetString(FlagBorChainId)
			if chainID == "" {
				return fmt.Errorf("ChainID cannot be empty")
			}

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				return fmt.Errorf("proposer address cannot be empty")
			}

			startBlockStr := viper.GetString(FlagStartBlock)
			if startBlockStr == "" {
				return fmt.Errorf("Start block cannot be empty")
			}

			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			//
			// Query data
			//

			// fetch duration
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", borTypes.QuerierRoute, bor.QueryParams, bor.ParamSpan), nil)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return errors.New("span duration not found")
			}

			var spanDuration uint64
			if err := cliCtx.Codec.UnmarshalJSON(res, &spanDuration); err != nil {
				return err
			}

			//
			// Get validators
			//

			res, err = cliCtx.QueryStore(staking.ACKCountKey, "staking")
			if err != nil {
				return err
			}

			// The query will return empty if there is no data
			if len(res) == 0 {
				return errors.New("No ack key found")
			}

			ackCount, err := strconv.ParseInt(string(res), 10, 64)
			if err != nil {
				return err
			}

			res, err = cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
			if err != nil {
				return err
			}
			// the query will return empty if there is no data
			if len(res) == 0 {
				return errors.New("No current validator set found")
			}

			var _validatorSet types.ValidatorSet
			cdc.UnmarshalBinaryBare(res, &_validatorSet)
			var validators []types.MinimalVal

			for _, val := range _validatorSet.Validators {
				if val.IsCurrentValidator(uint64(ackCount)) {
					// append if validator is current valdiator
					validators = append(validators, (*val).MinimalVal())
				}
			}

			msg := bor.NewMsgProposeSpan(
				proposer,
				startBlock,
				startBlock+spanDuration,
				validators,
				validators,
				chainID,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", helper.GetPubKey().Address().String(), "--proposer=<proposer-address>")
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=<bor-chain-id>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.MarkFlagRequired(FlagProposerAddress)
	cmd.MarkFlagRequired(FlagBorChainId)
	cmd.MarkFlagRequired(FlagStartBlock)

	return cmd
}
