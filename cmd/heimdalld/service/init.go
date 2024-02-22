package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/helper"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

type initHeimdallConfig struct {
	clientHome  string
	chainID     string
	validatorID int64
	chain       string
	forceInit   bool
}

func heimdallInit(_ *server.Context, cdc *codec.Codec, initConfig *initHeimdallConfig, config *cfg.Config) error {
	conf := helper.GetDefaultHeimdallConfig()
	conf.Chain = initConfig.chain
	WriteDefaultHeimdallConfig(filepath.Join(config.RootDir, "config/heimdall-config.toml"), conf)

	nodeID, valPubKey, _, err := InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	// do not execute init if forceInit is false and genesis.json already exists (or we do not have permission to write to file)
	writeGenesis := initConfig.forceInit

	if !writeGenesis {
		// When not forcing, check if genesis file exists
		_, err := os.Stat(config.GenesisFile())
		if err != nil && errors.Is(err, os.ErrNotExist) {
			logger.Info(fmt.Sprintf("Genesis file %v not found, writing genesis file\n", config.GenesisFile()))

			writeGenesis = true
		} else if err == nil {
			logger.Info(fmt.Sprintf("Found genesis file %v, skipping writing genesis file\n", config.GenesisFile()))
		} else {
			logger.Error(fmt.Sprintf("Error checking if genesis file %v exists: %v\n", config.GenesisFile(), err))
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("Force writing genesis file to %v\n", config.GenesisFile()))
	}

	if writeGenesis {
		genesisCreated, err := helper.WriteGenesisFile(initConfig.chain, config.GenesisFile(), cdc)
		if err != nil {
			return err
		} else if genesisCreated {
			return nil
		}
	} else {
		return nil
	}

	// create chain id
	chainID := initConfig.chainID
	if chainID == "" {
		chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
	}

	// get pubkey
	newPubkey := CryptoKeyToPubkey(valPubKey)

	// create validator account
	validator := hmTypes.NewValidator(hmTypes.NewValidatorID(uint64(initConfig.validatorID)),
		0, 0, 1, 1, newPubkey,
		hmTypes.BytesToHeimdallAddress(valPubKey.Address().Bytes()))

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
	appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, vals, *validatorSet)
	if err != nil {
		return err
	}

	// slashing state change
	appStateBytes, err = slashingTypes.SetGenesisStateToAppState(appStateBytes, valSigningInfoMap)
	if err != nil {
		return err
	}

	// bor state change
	appStateBytes, err = borTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet)
	if err != nil {
		return err
	}

	// topup state change
	appStateBytes, err = topupTypes.SetGenesisStateToAppState(appStateBytes, dividendAccounts)
	if err != nil {
		return err
	}

	// app state json
	appStateJSON, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(appStateBytes)
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
}

// InitCmd initialises files required to start heimdall
func initCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			initConfig := &initHeimdallConfig{
				chainID:     viper.GetString(client.FlagChainID),
				chain:       viper.GetString(helper.ChainFlag),
				validatorID: viper.GetInt64(stakingcli.FlagValidatorID),
				clientHome:  viper.GetString(helper.FlagClientHome),
				forceInit:   viper.GetBool(helper.OverwriteGenesisFlag),
			}
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			return heimdallInit(ctx, cdc, initConfig, config)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "Node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "Client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "Genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int(stakingcli.FlagValidatorID, 1, "--id=<validator ID here>, if left blank will be assigned 1")
	cmd.Flags().Bool(helper.OverwriteGenesisFlag, false, "Overwrite the genesis.json file if it exists")

	return cmd
}
