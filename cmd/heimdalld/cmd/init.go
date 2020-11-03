package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/iavl/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	stakingcli "github.com/maticnetwork/heimdall/x/staking/client/cli"

	authTypes "github.com/maticnetwork/heimdall/x/auth/types"
	// borTypes "github.com/maticnetwork/heimdall/x/bor/types"
	stakingTypes "github.com/maticnetwork/heimdall/x/staking/types"
	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
	commonTypes "github.com/maticnetwork/heimdall/types/common"
)

// InitCmd initialises files required to start heimdall
func initCmd(ctx *server.Context, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			fmt.Println("custom init")
			// create chain id
			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			validatorID := viper.GetInt64(stakingcli.FlagValidatorID)
			nodeID, valPubKey, _, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			WriteDefaultHeimdallConfig(filepath.Join(config.RootDir, "config/heimdall-config.toml"), helper.GetDefaultHeimdallConfig())

			// get pubkey
			newPubkey := CryptoKeyToPubkey(valPubKey)

			// create validator account
			validator := hmTypes.NewValidator(hmTypes.NewValidatorID(uint64(validatorID)),
				0, 0, 1, 1, newPubkey,
				commonTypes.BytesToHeimdallAddress(valPubKey.Address().Bytes()))

			// create dividend account for validator
			dividendAccount := hmTypes.NewDividendAccount(validator.Signer, ZeroIntString)

			vals := []*hmTypes.Validator{validator}
			validatorSet := hmTypes.NewValidatorSet(vals)

			dividendAccounts := []hmTypes.DividendAccount{dividendAccount}

			// create validator signing info
			valSigningInfo := hmTypes.NewValidatorSigningInfo(validator.ID, 0, 0, 0)
			valSigningInfoMap := make(map[string]hmTypes.ValidatorSigningInfo)
			valSigningInfoMap[valSigningInfo.ValID.String()] = valSigningInfo

			// create genesis state
			appStateBytes := app.NewDefaultGenesisState()

			// auth state change
			appStateBytes, err = authTypes.SetGenesisStateToAppState(
				appStateBytes,
				[]authTypes.GenesisAccount{getGenesisAccount(validator.Signer.Bytes())},
			)
			if err != nil {
				return err
			}

			// staking state change
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, vals, validatorSet)
			if err != nil {
				return err
			}

			// slashing state change
			// appStateBytes, err = slashingTypes.SetGenesisStateToAppState(appStateBytes, valSigningInfoMap)
			// if err != nil {
			// 	return err
			// }

			// bor state change
			// appStateBytes, err = borTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet)
			// if err != nil {
			// 	return err
			// }

			// topup state change
			appStateBytes, err = topupTypes.SetGenesisStateToAppState(appStateBytes, dividendAccounts)
			if err != nil {
				return err
			}

			// app state json
			appStateJSON, err := json.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID string `json:"chain_id"`
				NodeID  string `json:"node_id"`
			}{
				chainID,
				nodeID,
			}

			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return writeGenesisFile(tmtime.Now(), config.GenesisFile(), chainID, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int(stakingcli.FlagValidatorID, 1, "--id=<validator ID here>, if left blank will be assigned 1")
	return cmd
}
