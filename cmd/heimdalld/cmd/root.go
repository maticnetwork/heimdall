package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tendermint/tendermint/libs/cli"
	tmjson "github.com/tendermint/tendermint/libs/json"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/bor/console/prompt"

	"github.com/pborman/uuid"

	"github.com/maticnetwork/bor/accounts/keystore"

	ethCommon "github.com/maticnetwork/bor/common"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"

	ivalcommon "github.com/cosmos/iavl/common"
	bcrypto "github.com/maticnetwork/bor/crypto"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/app/params"
	"github.com/maticnetwork/heimdall/file"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sdkutil"
	"github.com/maticnetwork/heimdall/types/common"
)

var ZeroIntString = big.NewInt(0).String()

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome)

	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Heimdall Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := server.InterceptConfigsPreRunHandler(cmd); err != nil {
				return err
			}

			ctx := server.GetServerContextFromCmd(cmd)

			// bind flags for environment variables
			bindFlags(cmd, ctx.Viper)

			// load heimdall config
			if err := helper.InitHeimdallConfig(); err != nil {
				return err
			}

			return client.SetCmdClientContextHandler(initClientCtx, cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

// Execute executes the root command.
func Execute(rootCmd *cobra.Command) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, server.NewDefaultContext())

	executor := tmcli.PrepareBaseCmd(rootCmd, "HEIMDALL", app.DefaultNodeHome)
	return executor.ExecuteContext(ctx)
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	// Initialize sdk config with prefix and necessary configurations
	sdkutil.InitSDKConfig()
	authclient.Codec = encodingConfig.Marshaler

	_, legacyCodec := app.MakeCodecs()
	ctx := server.NewDefaultContext()
	// serverCtx := server.GetServerContextFromCmd(rootCmd)

	debugCmd := debug.Cmd()
	debugCmd.AddCommand(HeimdallPubkeyCmd())

	rootCmd.AddCommand(
		// TODO use default (cosmos-sdk) initCmd instead of custom
		// genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		initCmd(ctx, legacyCodec, app.ModuleBasics, app.DefaultNodeHome),
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		// keyring related commands
		keys.Commands(app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		// genutilcli.MigrateGenesisCmd(), // No need for migrate cmd for now
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics, encodingConfig.TxConfig),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debugCmd,
		showAccountCmd(),
		showPrivateKeyCmd(),
		generateKeystoreCmd(),
		generateValidatorKey(),
		convertAddressToHexCmd(),
		convertHexToAddressCmd(),
		exportCmd(ctx),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, createSimappAndExport, addModuleInitFlags)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	// crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
		flags.LineBreak,
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	return app.NewHeimdallApp(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		app.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

func createSimappAndExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailWhiteList []string,
	appOpts servertypes.AppOptions,
) (servertypes.ExportedApp, error) {
	encCfg := app.MakeEncodingConfig() // Ideally, we would reuse the one created by NewRootCmd.
	encCfg.Marshaler = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	var a app.App
	if height != -1 {
		a = app.NewHeimdallApp(logger, db, traceStore, false, map[int64]bool{}, "", uint(1), encCfg)

		if err := a.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		a = app.NewHeimdallApp(logger, db, traceStore, true, map[int64]bool{}, "", uint(1), encCfg)
	}

	return a.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// WriteGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
// nolint: unparam
func writeGenesisFile(genesisTime time.Time, genesisFile, chainID string, appState json.RawMessage) error {
	genDoc := tmtypes.GenesisDoc{
		GenesisTime: genesisTime,
		ChainID:     chainID,
		AppState:    appState,
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genesisFile)
}

// InitializeNodeValidatorFiles initializes node and priv validator files
func InitializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, valPubKey crypto.PubKey, priv crypto.PrivKey, err error,
) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := tmos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return "", nil, nil, err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := tmos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return "", nil, nil, err
	}

	filePV := privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)

	valPrivKey := filePV.Key.PrivKey

	// tmValPubKey, err := filePV.GetPubKey()
	// if err != nil {
	// 	return "", nil, nil, err
	// }

	// valPubKey, err = secp256k1.FromTmSecp256k1(tmValPubKey)
	// if err != nil {
	// 	return "", nil, nil, err
	// }
	valPubKey, _ = filePV.GetPubKey()
	return nodeID, valPubKey, valPrivKey, nil
}

func CryptoKeyToPubkey(key crypto.PubKey) common.PubKey {
	validatorPublicKey := helper.GetPubObjects(key)
	return common.NewPubKey(validatorPublicKey[:])
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		_ = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", "HEIMDALL", strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))))
		_ = v.BindPFlag(f.Name, f)

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address string `json:"address,omitempty" yaml:"address"`
	PrivKey string `json:"priv_key,omitempty" yaml:"priv_key"`
	PubKey  string `json:"pub_key,omitempty" yaml:"pub_key"`
}

func showAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's address and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// get public keys
			pubObject := helper.GetPubKey()

			account := &ValidatorAccountFormatter{
				Address: ethCommon.BytesToAddress(pubObject.Address().Bytes()).String(),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
			}

			b, err := json.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

func showPrivateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-privatekey",
		Short: "Print the account's private key",
		Run: func(cmd *cobra.Command, args []string) {
			// get private and public keys
			privObject := helper.GetPrivKey()

			account := &ValidatorAccountFormatter{
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
			}

			b, err := json.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

// exportCmd a state dump file
func exportCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-heimdall",
		Short: "Export genesis file with state-dump",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// create chain id
			chainID := viper.GetString(flags.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", ivalcommon.RandStr(6))
			}

			dataDir := path.Join(viper.GetString(cli.HomeFlag), "data")
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			db, err := sdk.NewLevelDB("application", dataDir)
			if err != nil {
				panic(err)
			}

			happ := app.NewHeimdallApp(logger, db, nil, true, map[int64]bool{}, viper.GetString(cli.HomeFlag), 5, app.MakeEncodingConfig())
			appState, err := happ.ExportAppStateAndValidators(false, []string{})

			if err != nil {
				panic(err)
			}

			err = writeGenesisFile(tmtime.Now(), file.Rootify("config/dump-genesis.json", config.RootDir), chainID, appState.AppState)
			if err == nil {
				fmt.Println("New genesis json file created:", file.Rootify("config/dump-genesis.json", config.RootDir))
			}
			return err
		},
	}
	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	return cmd
}

// generateKeystore generate keystore file from private key
func generateKeystoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-keystore <private-key>",
		Short: "Generates keystore file using private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := strings.ReplaceAll(args[0], "0x", "")
			pk, err := bcrypto.HexToECDSA(s)
			if err != nil {
				return err
			}

			id := uuid.NewRandom()
			key := &keystore.Key{
				Id:         id,
				Address:    bcrypto.PubkeyToAddress(pk.PublicKey),
				PrivateKey: pk,
			}

			passphrase, err := promptPassphrase(true)
			if err != nil {
				return err
			}

			keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
			if err != nil {
				return err
			}

			// Then write the new keyfile in place of the old one.
			if err := ioutil.WriteFile(keyFileName(key.Address), keyjson, 0600); err != nil {
				return err
			}
			return nil
		},
	}
}

// generateValidatorKey generate validator key
func generateValidatorKey() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-validatorkey <private-key>",
		Short: "Generate validator key file using private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := strings.ReplaceAll(args[0], "0x", "")
			ds, err := hex.DecodeString(s)
			if err != nil {
				return err
			}

			// set private object
			var privObject secp256k1.PrivKey
			copy(privObject[:], ds)

			// node key
			nodeKey := privval.FilePVKey{
				Address: privObject.PubKey().Address(),
				PubKey:  privObject.PubKey(),
				PrivKey: privObject,
			}

			jsonBytes, err := tmjson.MarshalIndent(nodeKey, "", "  ")
			if err != nil {
				return err
			}
			err = ioutil.WriteFile("priv_validator_key.json", jsonBytes, 0600)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

// convertAddressToHexCmd convert address to hex
func convertAddressToHexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "address-to-hex [address]",
		Short: "Convert address to hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			fmt.Println("Hex:", ethCommon.BytesToAddress(key).String())
			return nil
		},
	}
}

// convertHexToAddressCmd converts hex to address
func convertHexToAddressCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hex-to-address [hex]",
		Short: "Convert hex to address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := ethCommon.HexToAddress(args[0])
			fmt.Println("Address:", sdk.AccAddress(address.Bytes()).String())
			return nil
		},
	}
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address hex>
func keyFileName(keyAddr ethCommon.Address) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), hex.EncodeToString(keyAddr[:]))
}

func toISO8601(t time.Time) string {
	var tz string
	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

// promptPassphrase prompts the user for a passphrase.  Set confirmation to true
// to require the user to confirm the passphrase.
func promptPassphrase(confirmation bool) (string, error) {
	passphrase, err := prompt.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		return "", err
	}

	if confirmation {
		confirm, err := prompt.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			return "", err
		}

		if passphrase != confirm {
			return "", errors.New("Passphrases do not match")
		}
	}

	return passphrase, nil
}
