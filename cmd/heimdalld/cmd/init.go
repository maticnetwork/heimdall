package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	borTypes "github.com/maticnetwork/heimdall/x/bor/types"

	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/iavl/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	hmcommon "github.com/maticnetwork/heimdall/types/common"
	stakingtypes "github.com/maticnetwork/heimdall/x/staking/types"
)

// InitCmd initialises files required to start heimdall
func initCmd(ctx *server.Context, amino *codec.LegacyAmino, mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			cdc := clientCtx.JSONMarshaler

			config.SetRoot(clientCtx.HomeDir)
			// config.SetRoot(viper.GetString(cli.HomeFlag))
			// TODO : change default node home to flag
			// config.SetRoot(app.DefaultNodeHome)
			// create chain id
			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			nodeID, valPubKey, _, err := InitializeNodeValidatorFiles(config, "")
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(genutilcli.FlagOverwrite)

			if !overwrite && tmos.FileExists(genFile) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// get pubkey
			newPubkey := hmcommon.NewPubKey(valPubKey.Bytes())

			// create validator account
			validator := hmtypes.NewValidator(
				hmtypes.NewValidatorID(uint64(1)),
				0,
				0,
				1,
				1,
				newPubkey,
				valPubKey.Address().Bytes(),
			)

			// validators and validator set
			validators := []*hmtypes.Validator{validator}
			validatorSet := hmtypes.NewValidatorSet(validators)

			// signer address
			signer, _ := sdk.AccAddressFromHex(validator.Signer)
			// create dividend account for validator
			dividendAccount := hmtypes.NewDividendAccount(signer, ZeroIntString)
			// dividend accounts
			dividendAccounts := []*hmtypes.DividendAccount{&dividendAccount}

			// create validator signing info
			valSigningInfo := hmtypes.NewValidatorSigningInfo(validator.ID, 0, 0, 0)
			valSigningInfoMap := make(map[string]hmtypes.ValidatorSigningInfo)
			valSigningInfoMap[valSigningInfo.ValID.String()] = valSigningInfo

			// create genesis state
			// appStateBytes := app.NewDefaultGenesisState()
			appState := mbm.DefaultGenesis(cdc)
			// authState.Accounts = accounts
			// appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(&authState)

			genesisAccount := getGenesisAccount(signer.Bytes(), newPubkey)

			//
			// auth state change
			//
			authGenState := authtypes.GetGenesisStateFromAppState(authclient.Codec, appState)
			accounts, err := authtypes.UnpackAccounts(authGenState.Accounts)
			if err != nil {
				return fmt.Errorf("failed to unpack accounts: %w", err)
			}
			accounts = append(accounts, genesisAccount)
			genAccs, err := authtypes.PackAccounts(accounts)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}
			authGenState.Accounts = genAccs
			// TODO - check this
			appState[authtypes.ModuleName] = authclient.Codec.MustMarshalJSON(&authGenState)

			//
			// staking state change
			//
			_, err = stakingtypes.SetGenesisStateToAppState(authclient.Codec, appState, validators, validatorSet)
			if err != nil {
				return err
			}

			// // slashing state change
			//appStateBytes, err = slashingTypes.SetGenesisStateToAppState(appStateBytes, valSigningInfoMap)
			//if err != nil {
			//	return err
			//}

			// bor state change
			appState, err = borTypes.SetGenesisStateToAppState(appState, *validatorSet)
			if err != nil {
				return err
			}

			// topup state change
			appState, err = topupTypes.SetGenesisStateToAppState(appState, dividendAccounts)
			if err != nil {
				return err
			}

			// app state json
			appStateJSON, err := json.MarshalIndent(appState, "", " ")
			if err != nil {
				return errors.Wrap(err, "Failed to marshall default genesis state")
			}

			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

			if err := writeGenesisFile(tmtime.Now(), config.GenesisFile(), chainID, appStateJSON); err != nil {
				return err
			}

			// print info
			return displayInfo(newPrintInfo(config.Moniker, chainID, nodeID, "", appStateJSON))
		},
	}

	cmd.Flags().String(cli.HomeFlag, app.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().BoolP(genutilcli.FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	return cmd
}

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", string(sdk.MustSortJSON(out)))
	return err
}

func getGenesisAccount(address []byte, pk []byte) authtypes.GenesisAccount {
	acc := authtypes.NewBaseAccountWithAddress(address)
	pkObject := hmcommon.CosmosCryptoPubKey(pk)
	obj, err := codectypes.NewAnyWithValue(pkObject)
	if err != nil {
		panic(err)
	}
	acc.PubKey = obj
	return acc
}
