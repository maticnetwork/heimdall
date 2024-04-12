package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	tmTypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	ethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	bridgeCmd "github.com/maticnetwork/heimdall/bridge/cmd"
	"github.com/maticnetwork/heimdall/helper"
	restServer "github.com/maticnetwork/heimdall/server"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/version"
)

var logger = helper.Logger.With("module", "cmd/heimdalld")

var (
	flagNodeDirPrefix    = "node-dir-prefix"
	flagNumValidators    = "v"
	flagNumNonValidators = "n"
	flagOutputDir        = "output-dir"
	flagNodeDaemonHome   = "node-daemon-home"
	flagNodeCliHome      = "node-cli-home"
	flagNodeHostPrefix   = "node-host-prefix"
)

// Tendermint full-node start flags
const (
	flagAddress      = "address"
	flagTraceStore   = "trace-store"
	flagPruning      = "pruning"
	flagCPUProfile   = "cpu-profile"
	FlagMinGasPrices = "minimum-gas-prices"
	FlagHaltHeight   = "halt-height"
	FlagHaltTime     = "halt-time"
)

// Open Collector Flags
var (
	FlagOpenTracing           = "open-tracing"
	FlagOpenCollectorEndpoint = "open-collector-endpoint"
)

const (
	nodeDirPerm = 0755
)

var ZeroIntString = big.NewInt(0).String()

var hApp *app.HeimdallApp

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address string `json:"address,omitempty" yaml:"address"`
	PrivKey string `json:"priv_key,omitempty" yaml:"priv_key"`
	PubKey  string `json:"pub_key,omitempty" yaml:"pub_key"`
}

// GetSignerInfo returns signer information
func GetSignerInfo(pub crypto.PubKey, priv []byte, cdc *codec.Codec) ValidatorAccountFormatter {
	var privObject secp256k1.PrivKeySecp256k1

	cdc.MustUnmarshalBinaryBare(priv, &privObject)

	return ValidatorAccountFormatter{
		Address: ethCommon.BytesToAddress(pub.Address().Bytes()).String(),
		PubKey:  CryptoKeyToPubkey(pub).String(),
		PrivKey: "0x" + hex.EncodeToString(privObject[:]),
	}
}

func GetHeimdallApp() *app.HeimdallApp {
	tickCount := 0

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for hApp == nil {
		select { //nolint
		case <-tick.C:
			logger.Info("Waiting for heimdall app to be initialized")

			if hApp != nil {
				logger.Info("Heimdall app initialized")
				return hApp
			}

			tickCount++
			if tickCount > 10 {
				panic("Heimdall app not initialized")
			}
		}
	}

	return hApp
}

func NewHeimdallService(pCtx context.Context, args []string) {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	shutdownCtx, stop := signal.NotifyContext(pCtx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	rootCmd := &cobra.Command{
		Use:               "heimdalld",
		Short:             "Heimdall Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// adding heimdall configuration flags to root command
	helper.DecorateWithHeimdallFlags(rootCmd, viper.GetViper(), logger, "main")
	helper.DecorateWithTendermintFlags(rootCmd, viper.GetViper(), logger, "main")

	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	rootCmd.AddCommand(heimdallStart(shutdownCtx, ctx, getNewApp(ctx), cdc)) // New Heimdall start command

	tendermintCmd.AddCommand(
		server.ShowNodeIDCmd(ctx),
		server.ShowValidatorCmd(ctx),
		server.ShowAddressCmd(ctx),
		server.VersionCmd(ctx),
	)

	rootCmd.AddCommand(server.UnsafeResetAllCmd(ctx))
	rootCmd.AddCommand(flags.LineBreak)
	rootCmd.AddCommand(tendermintCmd)
	rootCmd.AddCommand(server.ExportCmd(ctx, cdc, exportAppStateAndTMValidators))
	rootCmd.AddCommand(flags.LineBreak)
	rootCmd.AddCommand(version.Cmd) // Using heimdall version, not Cosmos SDK version
	// End of block

	rootCmd.AddCommand(showAccountCmd())
	rootCmd.AddCommand(showPrivateKeyCmd())
	rootCmd.AddCommand(restServer.ServeCommands(shutdownCtx, cdc, restServer.RegisterRoutes))
	rootCmd.AddCommand(bridgeCmd.BridgeCommands(viper.GetViper(), logger, "main"))
	rootCmd.AddCommand(VerifyGenesis(ctx, cdc))
	rootCmd.AddCommand(initCmd(ctx, cdc))
	rootCmd.AddCommand(testnetCmd(ctx, cdc))

	// rollback cmd
	rootCmd.AddCommand(rollbackCmd(ctx))

	if args != nil && len(args) > 0 { //nolint
		rootCmd.SetArgs(args)
	}

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "HD", os.ExpandEnv("/var/lib/heimdall"))
	if err := executor.Execute(); err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func getNewApp(serverCtx *server.Context) func(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	return func(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
		// init heimdall config
		helper.InitHeimdallConfig("")
		helper.UpdateTendermintConfig(serverCtx.Config, viper.GetViper())
		// create new heimdall app
		hApp = app.NewHeimdallApp(logger, db,
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString(flagPruning))),
			baseapp.SetHaltHeight(viper.GetUint64(FlagHaltHeight)),
			baseapp.SetHaltTime(viper.GetUint64(FlagHaltTime)))

		return hApp
	}
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmTypes.GenesisValidator, error) {
	bapp := app.NewHeimdallApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}

func heimdallStart(shutdownCtx context.Context, ctx *server.Context, appCreator server.AppCreator, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in process.
Starting rest server is provided with the flag --rest-server and starting bridge with
the flag --bridge when starting Tendermint in process.
Pruning options can be provided via the '--pruning' flag. The options are as follows:

syncable: only those states not needed for state syncing will be deleted (keeps last 100 + every 10000th)
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if !strings.HasPrefix(arg, "--") {
					return fmt.Errorf(
						"\tinvalid argument: %s \n\tall flags must start with --",
						arg)
				}
			}
			LogsWriterFile := viper.GetString(helper.LogsWriterFileFlag)
			if LogsWriterFile != "" {
				logWriter := helper.GetLogsWriter(LogsWriterFile)

				logger, err := SetupCtxLogger(logWriter, ctx.Config.LogLevel)
				if err != nil {
					logger.Error("Unable to setup logger", "err", err)
					return err
				}

				ctx.Logger = logger
			}

			ctx.Logger.Info("starting ABCI with Tendermint")

			startRestServer, _ := cmd.Flags().GetBool(helper.RestServerFlag)
			startBridge, _ := cmd.Flags().GetBool(helper.BridgeFlag)

			err := startInProcess(cmd, shutdownCtx, ctx, appCreator, cdc, startRestServer, startBridge)
			return err
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			// bridge binding
			if err := viper.BindPFlag("all", cmd.Flags().Lookup("all")); err != nil {
				logger.Error("GetStartCmd | BindPFlag | all", "Error", err)
			}

			if err := viper.BindPFlag("only", cmd.Flags().Lookup("only")); err != nil {
				logger.Error("GetStartCmd | BindPFlag | only", "Error", err)
			}
		},
	}

	cmd.Flags().Bool(
		helper.RestServerFlag,
		false,
		"Start rest server",
	)

	cmd.Flags().Bool(
		helper.BridgeFlag,
		false,
		"Start bridge service",
	)

	cmd.PersistentFlags().String(helper.LogLevel, ctx.Config.LogLevel, "Log level")

	if err := viper.BindPFlag(helper.LogLevel, cmd.PersistentFlags().Lookup(helper.LogLevel)); err != nil {
		logger.Error("main | BindPFlag | helper.LogLevel", "Error", err)
	}

	// bridge flags =  start flags (all, only) + root bridge cmd flags
	cmd.Flags().Bool("all", false, "Start all bridge services")
	cmd.Flags().StringSlice("only", []string{}, "Comma separated bridge services to start")
	bridgeCmd.DecorateWithBridgeRootFlags(cmd, viper.GetViper(), logger, "main")

	// rest server flags
	restServer.DecorateWithRestFlags(cmd)

	// core flags for the ABCI application
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(
		FlagMinGasPrices, "",
		"Minimum gas prices to accept for transactions; Any fee in a tx must meet this minimum (e.g. 0.01photino;0.0001stake)",
	)
	cmd.Flags().Uint64(FlagHaltHeight, 0, "Height at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().Uint64(FlagHaltTime, 0, "Minimum block time (in Unix seconds) at which to gracefully halt the chain and shutdown the node")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "Client's home directory")

	cmd.Flags().Bool(FlagOpenTracing, false, "Start open tracing")
	cmd.Flags().String(FlagOpenCollectorEndpoint, helper.DefaultOpenCollectorEndpoint, "Default OpenTelemetry Collector Endpoint")

	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)

	return cmd
}

func startOpenTracing(cmd *cobra.Command) (*sdktrace.TracerProvider, *context.Context, error) {
	opentracingEnabled, _ := cmd.Flags().GetBool(FlagOpenTracing)
	if opentracingEnabled {
		openCollectorEndpoint, _ := cmd.Flags().GetString(FlagOpenCollectorEndpoint)
		ctx := context.Background()

		res, err := resource.New(ctx,
			resource.WithAttributes(
				// the service name used to display traces in backends
				semconv.ServiceNameKey.String("heimdall"),
			),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create open telemetry resource for service: %v", err)
		}

		// Set up a trace exporter
		var traceExporter *otlptrace.Exporter

		traceExporterReady := make(chan *otlptrace.Exporter, 1)

		go func() {
			traceExporter, _ := otlptracegrpc.New(
				ctx,
				otlptracegrpc.WithInsecure(),
				otlptracegrpc.WithEndpoint(openCollectorEndpoint),
				otlptracegrpc.WithDialOption(grpc.WithBlock()),
			)
			traceExporterReady <- traceExporter
		}()

		select {
		case traceExporter = <-traceExporterReady:
			fmt.Println("TraceExporter Ready")
		case <-time.After(5 * time.Second):
			fmt.Println("TraceExporter Timed Out in 5 Seconds")
		}

		// Register the trace exporter with a TracerProvider, using a batch
		// span processor to aggregate spans before export.
		if traceExporter == nil {
			return nil, nil, nil
		}

		bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
		tracerProvider := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(res),
			sdktrace.WithSpanProcessor(bsp),
		)
		otel.SetTracerProvider(tracerProvider)

		// set global propagator to tracecontext (the default is no-op).
		otel.SetTextMapPropagator(propagation.TraceContext{})

		return tracerProvider, &ctx, nil
	}

	return nil, nil, nil
}

func startInProcess(cmd *cobra.Command, shutdownCtx context.Context, ctx *server.Context, appCreator server.AppCreator, cdc *codec.Codec, startRestServer bool, startBridge bool) error {
	cfg := ctx.Config
	home := cfg.RootDir
	traceWriterFile := viper.GetString(flagTraceStore)

	// initialize heimdall if needed (do not force!)
	initConfig := &initHeimdallConfig{
		chainID:     "", // chain id should be auto generated if chain flag is not set to mumbai, amoy or mainnet
		chain:       viper.GetString(helper.ChainFlag),
		validatorID: 1, // default id for validator
		clientHome:  viper.GetString(helper.FlagClientHome),
		forceInit:   false,
	}

	if err := heimdallInit(ctx, cdc, initConfig, cfg); err != nil {
		return fmt.Errorf("failed init heimdall: %s", err)
	}

	db, err := openDB(home)
	if err != nil {
		return fmt.Errorf("failed to open DB: %s", err)
	}

	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return fmt.Errorf("failed to open trace writer: %s", err)
	}

	app := appCreator(ctx.Logger, db, traceWriter)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return fmt.Errorf("failed to load or gen node key: %s", err)
	}

	server.UpgradeOldPrivValFile(cfg)

	// create & start tendermint node
	tmNode, err := node.NewNode(
		cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return fmt.Errorf("failed to create new node: %s", err)
	}

	// start Tendermint node here
	if err = tmNode.Start(); err != nil {
		return fmt.Errorf("failed to start Tendermint node: %s", err)
	}

	var cpuProfileCleanup func()

	if cpuProfile := viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)

		if err = pprof.StartCPUProfile(f); err != nil {
			return err
		}

		cpuProfileCleanup = func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}
	}

	tracerProvider, traceCtx, _ := startOpenTracing(cmd)

	// using group context makes sense in case that if one of
	// the processes produces error the rest will go and shutdown
	g, gCtx := errgroup.WithContext(shutdownCtx)
	// start rest
	if startRestServer {
		waitForREST := make(chan struct{})

		g.Go(func() error {
			return restServer.StartRestServer(gCtx, cdc, restServer.RegisterRoutes, waitForREST)
		})

		// hang here for a while, and wait for REST server to start
		<-waitForREST
	}

	// start bridge
	if startBridge {
		bridgeCmd.AdjustBridgeDBValue(cmd, viper.GetViper())
		g.Go(func() error {
			return bridgeCmd.StartBridgeWithCtx(gCtx)
		})
	}

	// stop phase for Tendermint node
	g.Go(func() error {
		// wait here for interrupt signal or
		// until something in the group returns non-nil error
		<-gCtx.Done()
		ctx.Logger.Info("exiting...")

		if tracerProvider != nil {
			// nolint: contextcheck
			if err := tracerProvider.Shutdown(*traceCtx); err == nil {
				ctx.Logger.Info("Shutting Down OpenTelemetry")
			}
		}

		if cpuProfileCleanup != nil {
			cpuProfileCleanup()
		}
		if tmNode.IsRunning() {
			return tmNode.Stop()
		}

		db.Close()

		return nil
	})

	// wait here for all go routines to finish,
	// or something to break
	if err := g.Wait(); err != nil {
		ctx.Logger.Error("Error shutting down services", "Error", err)
		return err
	}

	logger.Info("Heimdall services stopped")

	return nil
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB("application", dataDir)
}

func openTraceWriter(traceWriterFile string) (io.Writer, error) {
	if traceWriterFile == "" {
		return nil, nil
	}

	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0666,
	)
}

func showAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's address and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig("")

			// get public keys
			pubObject := helper.GetPubKey()

			account := &ValidatorAccountFormatter{
				Address: ethCommon.BytesToAddress(pubObject.Address().Bytes()).String(),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
			}

			b, err := jsoniter.ConfigFastest.MarshalIndent(account, "", "    ")
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
			// init heimdall config
			helper.InitHeimdallConfig("")

			// get private and public keys
			privObject := helper.GetPrivKey()

			account := &ValidatorAccountFormatter{
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
			}

			b, err := jsoniter.ConfigFastest.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

// VerifyGenesis verifies the genesis file and brings it in sync with on-chain contract
func VerifyGenesis(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-genesis",
		Short: "Verify if the genesis matches",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			helper.InitHeimdallConfig("")

			// Loading genesis doc
			genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(config.RootDir, "config/genesis.json"))
			if err != nil {
				return err
			}

			// get genesis state
			var genesisState app.GenesisState
			err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(genDoc.AppState, &genesisState)
			if err != nil {
				return err
			}

			// verify genesis
			for _, b := range app.ModuleBasics {
				m := b.(hmModule.HeimdallModuleBasic)
				if err := m.VerifyGenesis(genesisState); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}

// Total Validators to be included in the testnet
func totalValidators() int {
	numValidators := viper.GetInt(flagNumValidators)
	numNonValidators := viper.GetInt(flagNumNonValidators)

	return numNonValidators + numValidators
}

// get node directory path
func nodeDir(i int) string {
	outDir := viper.GetString(flagOutputDir)
	nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
	nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)

	return filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
}

// hostname of ip of nodes
func hostnameOrIP(i int) string {
	return fmt.Sprintf("%s%d", viper.GetString(flagNodeHostPrefix), i)
}

// populate persistent peers in config
func populatePersistentPeersInConfigAndWriteIt(config *cfg.Config) {
	persistentPeers := make([]string, totalValidators())

	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))

		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return
		}

		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", hostnameOrIP(i), 26656))
	}

	persistentPeersList := strings.Join(persistentPeers, ",")

	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = false

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir(i), "config", "config.toml"), config)
	}
}

func getGenesisAccount(address []byte) authTypes.GenesisAccount {
	acc := authTypes.NewBaseAccountWithAddress(hmTypes.BytesToHeimdallAddress(address))

	genesisBalance, _ := big.NewInt(0).SetString("1000000000000000000000", 10)

	if err := acc.SetCoins(sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: sdk.NewIntFromBigInt(genesisBalance)}}); err != nil {
		logger.Error("getGenesisAccount | SetCoins", "Error", err)
	}

	result, _ := authTypes.NewGenesisAccountI(&acc)

	return result
}

// WriteGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
// nolint: unparam
func writeGenesisFile(genesisTime time.Time, genesisFile, chainID string, appState json.RawMessage) error {
	genDoc := tmTypes.GenesisDoc{
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
		return nodeID, valPubKey, priv, err
	}

	nodeID = string(nodeKey.ID())

	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := common.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, priv, err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := common.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, priv, err
	}

	FilePv := privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	valPubKey = FilePv.GetPubKey()

	return nodeID, valPubKey, FilePv.Key.PrivKey, nil
}

// WriteDefaultHeimdallConfig writes default heimdall config to the given path
func WriteDefaultHeimdallConfig(path string, conf helper.Configuration) {
	// Don't write if config file in path already exists
	if _, err := os.Stat(path); err == nil {
		logger.Info(fmt.Sprintf("Config file %s already exists. Skip writing default heimdall config.", path))
	} else if errors.Is(err, os.ErrNotExist) {
		helper.WriteConfigFile(path, &conf)
	} else {
		logger.Error("Error while checking for config file", "Error", err)
	}
}

func CryptoKeyToPubkey(key crypto.PubKey) hmTypes.PubKey {
	validatorPublicKey := helper.GetPubObjects(key)
	return hmTypes.NewPubKey(validatorPublicKey[:])
}

func SetupCtxLogger(logWriter io.Writer, logLevel string) (log.Logger, error) {
	logger := log.NewTMLogger(log.NewSyncWriter(logWriter))

	logger, err := tmflags.ParseLogLevel(logLevel, logger, cfg.DefaultLogLevel())
	if err != nil {
		return nil, err
	}

	if viper.GetBool(cli.TraceFlag) {
		logger = log.NewTracingLogger(logger)
	}

	logger = logger.With("module", "main")

	return logger, nil
}
