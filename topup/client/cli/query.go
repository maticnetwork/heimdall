package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "topup/client/cli")

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group topup queries under a subcommand
	topupQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the topup module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// topup query command
	topupQueryCmd.AddCommand(
		client.GetCommands(
			GetSequence(cdc),
			IsOldTx(cdc),
			GetDividendAccount(cdc),
			GetDividendAccountRoot(cdc),
			GetAccountProof(cdc),
			GetAccountProofVerify(cdc),
		)...,
	)

	return topupQueryCmd
}

// GetSequence from the txhash and logIndex
func GetSequence(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get sequence from txhash and logindex",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			logIndex := viper.GetUint64(FlagLogIndex)
			txHashStr := viper.GetString(FlagTxHash)
			if txHashStr == "" {
				return fmt.Errorf("LogIndex and transaction hash required")
			}

			var queryParams []byte
			var err error
			var t string = ""
			if txHashStr != "" {
				queryParams, err = cliCtx.Codec.MarshalJSON(types.NewQuerySequenceParams(txHashStr, logIndex))
				if err != nil {
					return err
				}
				t = types.QuerySequence
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, t), queryParams)
			if err != nil {
				fmt.Println("No topup exists")
				// nolint: nilerr
				return nil
			}

			fmt.Println("Success. Topup exists with sequence:", string(res))

			return nil
		},
	}

	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		cliLogger.Error("GetSequence | MarkFlagRequired | FlagTxHash", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		cliLogger.Error("GetSequence | MarkFlagRequired | FlagLogIndex", "Error", err)
	}

	return cmd
}

// Check whether the transaction is old
func IsOldTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-old-tx",
		Short: "Check whether the transaction is old",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// tx hash
			txHash := viper.GetString(FlagTxHash)
			if txHash == "" {
				return fmt.Errorf("tx hash cannot be empty")
			}

			// log index
			logIndex := viper.GetUint64(FlagLogIndex)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySequenceParams(txHash, logIndex))
			if err != nil {
				return err
			}

			seqNo, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySequence), queryParams)
			if err != nil {
				return err
			}

			// error if no tx status found
			if len(seqNo) == 0 {
				fmt.Printf("false")
				return nil
			}

			res := true

			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log index here>")
	cmd.Flags().Uint64(FlagTxHash, 0, "--tx-hash=<tx hash here>")

	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagLogIndex", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagTxHash", "Error", err)
	}

	return cmd
}

func GetDividendAccount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dividend-account",
		Short: "show dividend account via address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorAddressStr := viper.GetString(FlagValidatorAddress)

			userAddress := hmTypes.HexToHeimdallAddress(validatorAddressStr)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryDividendAccountParams(userAddress))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccount), queryParams)
			if err != nil {
				return err
			}

			// error if no dividend account found
			if len(res) == 0 {
				fmt.Printf("Not found")
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")

	if err := cmd.MarkFlagRequired(FlagValidatorAddress); err != nil {
		logger.Error("GetDividendAccount | MarkFlagRequired | FlagValidatorAddress", "Error", err)
	}

	return cmd
}

func GetDividendAccountRoot(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dividend-account-root",
		Short: "show dividend account root",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccountRoot), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				fmt.Printf("Account root not found")
				return nil
			}

			var accountRootHash = hmTypes.BytesToHeimdallHash(res)

			result, err := json.Marshal(&accountRootHash)
			if err != nil {
				return err
			}

			// error if no dividend account found
			if len(res) == 0 {
				fmt.Printf("Not found")
				return nil
			}

			fmt.Println(string(result))
			return nil
		},
	}

	return cmd
}

func GetAccountProof(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-proof",
		Short: "show account proof via address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorAddressStr := viper.GetString(FlagValidatorAddress)

			userAddress := hmTypes.HexToHeimdallAddress(validatorAddressStr)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryAccountProofParams(userAddress))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountProof), queryParams)
			if err != nil {
				return err
			}

			// error if no dividend account found
			if len(res) == 0 {
				fmt.Printf("Not found")
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")

	if err := cmd.MarkFlagRequired(FlagValidatorAddress); err != nil {
		logger.Error("GetAccountProof | MarkFlagRequired | FlagValidatorAddress", "Error", err)
	}

	return cmd
}

func GetAccountProofVerify(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account-proof-verify",
		Short: "show account proof via address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorAddressStr := viper.GetString(FlagValidatorAddress)
			userAddress := hmTypes.HexToHeimdallAddress(validatorAddressStr)

			accountProof := viper.GetString(FlagAccountProof)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryVerifyAccountProofParams(userAddress, accountProof))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVerifyAccountProof), queryParams)
			if err != nil {
				return err
			}

			var accountProofStatus bool
			if err = json.Unmarshal(res, &accountProofStatus); err != nil {
				return err
			}

			res, err = json.Marshal(map[string]interface{}{"result": accountProofStatus})
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")
	cmd.Flags().String(FlagAccountProof, "", "--proof=<proof here>")

	if err := cmd.MarkFlagRequired(FlagValidatorAddress); err != nil {
		logger.Error("GetAccountProofVerify | MarkFlagRequired | FlagValidatorAddress", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagAccountProof); err != nil {
		logger.Error("GetAccountProofVerify | MarkFlagRequired | FlagAccountProof", "Error", err)
	}

	return cmd
}
