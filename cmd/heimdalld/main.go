package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	//"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	cfg "github.com/tendermint/tendermint/config"
	//tmcli "github.com/tendermint/tendermint/libs/cli"
	//"github.com/tendermint/tendermint/p2p"
	//"path/filepath"

	//"github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	gaiaAppInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	//"github.com/cosmos/cosmos-sdk/x/auth"
	//"github.com/tendermint/tendermint/crypto"
	//cmn "github.com/tendermint/tendermint/libs/common"
	//"io/ioutil"
	//"path"
	//"sort"
	"github.com/cosmos/cosmos-sdk/client"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"path/filepath"

	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	"io/ioutil"
	"sort"
	"strings"
	"github.com/tendermint/tendermint/p2p"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "heimdalld",
		Short:             "Heimdall Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// add new persistent flag for heimdall-config
	rootCmd.PersistentFlags().String(
		helper.WithHeimdallConfigFlag,
		"",
		"Heimdall config file path (default <home>/config/heimdall-config.json)",
	)

	// bind with-heimdall-config config with root cmd
	viper.BindPFlag(
		helper.WithHeimdallConfigFlag,
		rootCmd.Flags().Lookup(helper.WithHeimdallConfigFlag),
	)

	// add custom root command
	rootCmd.AddCommand(newAccountCmd())
	// cosmos server commands
	server.AddCommands(
		ctx,
		cdc,
		rootCmd,
		server.DefaultAppInit,
		server.AppCreator(newApp),
		server.AppExporter(exportAppStateAndTMValidators),
	)

	//rootCmd.AddCommand(InitCmd(ctx, cdc, DefaultAppInit))
	rootCmd.AddCommand(InitBasecoinCmd(ctx, cdc, server.DefaultAppInit))
	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdalld"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	// init heimdall config
	helper.InitHeimdallConfig()

	// create new heimdall app
	return app.NewHeimdallApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewHeimdallApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}

func newAccountCmd() *cobra.Command {
	type Account struct {
		Address string `json:"address"`
		PrivKey string `json:"private_key"`
		PubKey  string `json:"public_key"`
	}

	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's private key and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig()

			// get private and public keys
			privObject := helper.GetPrivKey()
			pubObject := helper.GetPubKey()

			account := &Account{
				Address: "0x" + hex.EncodeToString(pubObject.Address().Bytes()),
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
			}

			b, err := json.Marshal(&account)
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

const (
	flagWithTxs      = "with-txs"
	flagOverwrite    = "overwrite"
	flagClientHome   = "home-client"
	flagOverwriteKey = "overwrite-key"
	flagSkipGenesis  = "skip-genesis"
	flagMoniker      = "moniker"
)

type initConfig struct {
	ChainID      string
	GenTxsDir    string
	Name         string
	NodeID       string
	ClientHome   string
	WithTxs      bool
	Overwrite    bool
	OverwriteKey bool
	ValPubKey    crypto.PubKey
}

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	AppMessage json.RawMessage `json:"app_message"`
}

func displayInfo(cdc *codec.Codec, info printInfo) error {
	out, err := codec.MarshalJSONIndent(cdc, info)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", string(out))
	return nil
}

// get cmd to initialize all files for tendermint and application
// nolint
func InitCmd(ctx *server.Context, cdc *codec.Codec, appInit AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long: `Initialize validators's and node's configuration files.

Note that only node's configuration files will be written if the flag --skip-genesis is
enabled, and the genesis file will not be generated.
`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			name := viper.GetString(client.FlagName)
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}
			nodeID, valPubKey, err := gaiaAppInit.InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			if viper.GetString(flagMoniker) != "" {
				config.Moniker = viper.GetString(flagMoniker)
			}
			if config.Moniker == "" && name != "" {
				config.Moniker = name
			}
			toPrint := printInfo{
				ChainID: chainID,
				Moniker: config.Moniker,
				NodeID:  nodeID,
			}
			if viper.GetBool(flagSkipGenesis) {
				cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
				return displayInfo(cdc, toPrint)
			}

			initCfg := initConfig{
				ChainID:      chainID,
				GenTxsDir:    filepath.Join(config.RootDir, "config", "gentx"),
				Name:         name,
				NodeID:       nodeID,
				ClientHome:   viper.GetString(flagClientHome),
				WithTxs:      viper.GetBool(flagWithTxs),
				Overwrite:    viper.GetBool(flagOverwrite),
				OverwriteKey: viper.GetBool(flagOverwriteKey),
				ValPubKey:    valPubKey,
			}
			_, err = initWithConfig(cdc, config, initCfg, appInit)
			// print out some key information
			if err != nil {
				return err
			}

			//toPrint.AppMessage = json.Marshal("")
			return displayInfo(cdc, toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(flagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool(flagWithTxs, false, "apply existing genesis transactions from [--home]/config/gentx/")
	cmd.Flags().String(client.FlagName, "", "name of private key with which to sign the gentx")
	cmd.Flags().String(flagMoniker, "", "overrides --name flag and set the validator's moniker to a different value; ignored if it runs without the --with-txs flag")
	cmd.Flags().String(flagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().Bool(flagOverwriteKey, false, "overwrite client's key")
	cmd.Flags().Bool(flagSkipGenesis, false, "do not create genesis.json")
	return cmd
}


// get cmd to initialize all files for tendermint and application
// nolint: errcheck
func InitBasecoinCmd(ctx *server.Context, cdc *codec.Codec, appInit server.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
			}

			nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
			if err != nil {
				return err
			}
			nodeID := string(nodeKey.ID())

			pk := gaiaInit.ReadOrCreatePrivValidator(config.PrivValidatorFile())
			genTx, appMessage, validator, err := server.SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := appInit.AppGenState(cdc, []json.RawMessage{genTx})
			if err != nil {
				return err
			}
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return gaiaInit.WriteGenesisFile(config.GenesisFile(), chainID, []tmtypes.GenesisValidator{validator}, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(client.FlagName, "", "validator's moniker")
	cmd.MarkFlagRequired(client.FlagName)
	return cmd
}


func initWithConfig(cdc *codec.Codec, config *cfg.Config, initCfg initConfig, appInit AppInit) (
	appMessage json.RawMessage, err error) {
	genFile := config.GenesisFile()
	if !initCfg.Overwrite && common.FileExists(genFile) {
		err = fmt.Errorf("genesis.json file already exists: %v", genFile)
		return
	}

	// process genesis transactions, else create default genesis.json
	var appGenTxs []auth.StdTx
	var persistentPeers string
	var genTxs []json.RawMessage
	var appState json.RawMessage
	var jsonRawTx json.RawMessage
	chainID := initCfg.ChainID

	if initCfg.WithTxs {
		_, appGenTxs, persistentPeers, err = CollectStdTxs(config.Moniker, initCfg.GenTxsDir, cdc)
		if err != nil {
			return
		}
		genTxs = make([]json.RawMessage, len(appGenTxs))
		config.P2P.PersistentPeers = persistentPeers
		for i, stdTx := range appGenTxs {
			jsonRawTx, err = cdc.MarshalJSON(stdTx)
			if err != nil {
				return
			}
			genTxs[i] = jsonRawTx
		}
	} else {
		//var keyPass, secret string
		//var addr sdk.AccAddress
		//var signedTx auth.StdTx
		//var ip string

		if initCfg.Name == "" {
			err = errors.New("must specify validator's moniker (--name)")
			return
		}

		config.Moniker = initCfg.Name

	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	appState, err = appInit.AppGenState(cdc, genTxs)
	if err != nil {
		return
	}

	err = gaiaAppInit.WriteGenesisFile(genFile, chainID, nil, appState)

	return
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns the list of validators,
// appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(moniker string, genTxsDir string, cdc *codec.Codec) (
	validators []tmtypes.GenesisValidator, appGenTxs []auth.StdTx, persistentPeers string, err error) {
	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return
	}

	var addresses []string
	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		jsonRawTx, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		var genStdTx auth.StdTx
		err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx)
		if err != nil {
			return
		}
		appGenTxs = append(appGenTxs, genStdTx)

		nodeAddr := genStdTx.GetMemo()
		if len(nodeAddr) == 0 {
			err = fmt.Errorf("couldn't find node's address in %s", fo.Name())
			return
		}

		msgs := genStdTx.GetMsgs()
		if len(msgs) != 1 {
			err = errors.New("each genesis transaction must provide a single genesis message")
			return
		}

		//msg := msgs[0].(stake.MsgCreateValidator)
		//validators = append(validators, tmtypes.GenesisValidator{
		//	PubKey: msg.PubKey,
		//	Power:  freeFermionVal,
		//	Name:   msg.Description.Moniker,
		//})
		//
		//// exclude itself from persistent peers
		//if msg.Description.Moniker != moniker {
		//	addresses = append(addresses, nodeAddr)
		//}
	}

	sort.Strings(addresses)
	persistentPeers = strings.Join(addresses, ",")

	return
}

// SimpleGenTx is a simple genesis tx
type SimpleGenTx struct {
	Addr sdk.AccAddress `json:"addr"`
}

// create the genesis app state
func SimpleAppGenState(cdc *codec.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	var tx SimpleGenTx
	//err = cdc.UnmarshalJSON(appGenTxs[0], &tx)
	//if err != nil {
	//	return
	//}

	appState = json.RawMessage(fmt.Sprintf(`{
  "accounts": [{
    "address": "%s",
    "coins": [
      {
        
      }
    ]
  }]
}`, tx.Addr))
	return
}

var DefaultAppInit = AppInit{
	AppGenState: SimpleAppGenState,
}

type AppInit struct {
	// AppGenState creates the core parameters initialization. It takes in a
	// pubkey meant to represent the pubkey of the validator of this machine.
	AppGenState func(cdc *codec.Codec, appGenTx []json.RawMessage) (appState json.RawMessage, err error)
}
