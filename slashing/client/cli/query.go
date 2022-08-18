package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group slashing queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the slashing module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// slashingQueryCmd query command
	slashingQueryCmd.AddCommand(
		client.GetCommands(
			// GetCmdQuerySigningInfo(cdc),
			GetCmdQueryParams(cdc),
			GetSigningInfo(cdc),
			GetSigningInfos(cdc),
			GetLatestSlashInfo(cdc),
			GetLatestSlashingInfos(cdc),
			GetTickSlashingInfos(cdc),
			GetLatestSlashInfoBytes(cdc),
			GetTickCount(cdc),
			IsOldTx(cdc),
		)...,
	)
	return slashingQueryCmd

}

/* // GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "signing-info [validator-id]",
		Short: "Query a validator's signing information",
		Long: strings.TrimSpace(`Use a validators' id to find the signing-info for that validator:

$ <appcli> query slashing signing-info {valID}
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validatorID := viper.GetUint64(FlagValidatorID)

			if validatorID == 0 {
				return fmt.Errorf("validator ID is required")
			}

			key := types.GetValidatorSigningInfoKey(hmTypes.NewValidatorID(validatorID).Bytes())

			res, _, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("validator %s not found in slashing store", validatorID)
			}

			var signingInfo hmTypes.ValidatorSigningInfo
			signingInfo, err = hmTypes.UnmarshallValSigningInfo(types.ModuleCdc, res)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(signingInfo)
		},
	}
} */

//Give signing info by id
func GetSigningInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-info",
		Short: "show signing-info by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			id := viper.GetUint64(FlagId)

			params := types.NewQuerySigningInfoParams(hmTypes.ValidatorID(id))

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfo)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagId, "", "--id=<id here>")

	if err := cmd.MarkFlagRequired(FlagId); err != nil {
		logger.Error("GetSigningInfo | MarkFlagRequired | FlagId", "Error", err)
	}

	return cmd
}

//Give signing info by id
func GetSigningInfos(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-infos",
		Short: "show signing-info by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(FlagPage)

			limit := viper.GetInt(FlagLimit)

			params := types.NewQuerySigningInfosParams(page, limit)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfos)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		logger.Error("GetSigningInfos | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		logger.Error("GetSigningInfos | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}

func GetLatestSlashInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashing-info",
		Short: "show latest slash info by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			id := viper.GetUint64(FlagId)

			params := types.NewQuerySlashingInfoParams(hmTypes.ValidatorID(id))

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfo)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagId, "", "--id=<id here>")

	if err := cmd.MarkFlagRequired(FlagId); err != nil {
		logger.Error("GetLatestSlashInfo | MarkFlagRequired | FlagId", "Error", err)
	}

	return cmd
}

//Give signing info by id
func GetLatestSlashingInfos(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slashing-infos",
		Short: "show slashing infos",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(FlagPage)

			limit := viper.GetInt(FlagLimit)

			params := types.NewQuerySlashingInfosParams(page, limit)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfos)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		logger.Error("GetLatestSlashingInfos | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		logger.Error("GetLatestSlashingInfos | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}

//Give tick slash infos
func GetTickSlashingInfos(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tick-slash-infos",
		Short: "show tick slash infos",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(FlagPage)

			limit := viper.GetInt(FlagLimit)

			params := types.NewQueryTickSlashingInfosParams(page, limit)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTickSlashingInfos)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		logger.Error("GetTickSlashingInfos | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		logger.Error("GetTickSlashingInfos | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}

func GetLatestSlashInfoBytes(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-slash-info-bytes",
		Short: "Give the latest slash info bytes",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfoBytes), nil)

			if err != nil {
				return err
			}

			// error if no slashInfoBytes found
			if len(res) == 0 {
				fmt.Printf("no slashing Bytes found")
				return nil
			}

			var slashInfoBytes = hmTypes.BytesToHexBytes(res)

			result, err := json.Marshal(&slashInfoBytes)
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil

		},
	}

	return cmd
}

// GetCmdQueryParams implements a command to fetch slashing parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current slashing parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query genesis parameters for the slashing module:

$ <appcli> query slashing params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
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
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySlashingSequenceParams(txHash, logIndex))
			if err != nil {
				return err
			}

			seqNo, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingSequence), queryParams)
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

func GetTickCount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tick-count",
		Short: "Give the tick-count",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTickCount), nil)

			if err != nil {
				return err
			}

			// error if no slashInfoBytes found
			if len(res) == 0 {
				fmt.Printf("no slashing Bytes found")
				return nil
			}

			var tickCount uint64
			if err := json.Unmarshal(res, &tickCount); err != nil {
				return err
			}

			result, err := json.Marshal(&tickCount)
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil

		},
	}

	return cmd
}
