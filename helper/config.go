package helper

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	logger "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"

	"github.com/maticnetwork/heimdall/file"
	hmTypes "github.com/maticnetwork/heimdall/types"

	cfg "github.com/tendermint/tendermint/config"
	tmTypes "github.com/tendermint/tendermint/types"
)

const (
	TendermintNodeFlag     = "node"
	WithHeimdallConfigFlag = "heimdall-config"
	HomeFlag               = "home"
	FlagClientHome         = "home-client"
	OverwriteGenesisFlag   = "overwrite-genesis"
	RestServerFlag         = "rest-server"
	BridgeFlag             = "bridge"
	LogLevel               = "log_level"
	LogsWriterFileFlag     = "logs_writer_file"
	SeedsFlag              = "seeds"

	MainChain   = "mainnet"
	MumbaiChain = "mumbai"
	LocalChain  = "local"

	// heimdall-config flags
	MainRPCUrlFlag               = "eth_rpc_url"
	BorRPCUrlFlag                = "bor_rpc_url"
	TendermintNodeURLFlag        = "tendermint_rpc_url"
	HeimdallServerURLFlag        = "heimdall_rest_server"
	AmqpURLFlag                  = "amqp_url"
	CheckpointerPollIntervalFlag = "checkpoint_poll_interval"
	SyncerPollIntervalFlag       = "syncer_poll_interval"
	NoACKPollIntervalFlag        = "noack_poll_interval"
	ClerkPollIntervalFlag        = "clerk_poll_interval"
	SpanPollIntervalFlag         = "span_poll_interval"
	MainchainGasLimitFlag        = "main_chain_gas_limit"
	MainchainMaxGasPriceFlag     = "main_chain_max_gas_price"
	NoACKWaitTimeFlag            = "no_ack_wait_time"
	ChainFlag                    = "chain"

	// ---
	// TODO Move these to common client flags
	// BroadcastBlock defines a tx broadcasting mode where the client waits for
	// the tx to be committed in a block.
	BroadcastBlock = "block"

	// BroadcastSync defines a tx broadcasting mode where the client waits for
	// a CheckTx execution response only.
	BroadcastSync = "sync"

	// BroadcastAsync defines a tx broadcasting mode where the client returns
	// immediately.
	BroadcastAsync = "async"
	// --

	// RPC Endpoints
	DefaultMainRPCUrl = "http://localhost:9545"
	DefaultBorRPCUrl  = "http://localhost:8545"

	// RPC Timeouts
	DefaultEthRPCTimeout = 5 * time.Second
	DefaultBorRPCTimeout = 5 * time.Second

	// Services

	// DefaultAmqpURL represents default AMQP url
	DefaultAmqpURL           = "amqp://guest:guest@localhost:5672/"
	DefaultHeimdallServerURL = "http://0.0.0.0:1317"
	DefaultTendermintNodeURL = "http://0.0.0.0:26657"

	NoACKWaitTime = 1800 * time.Second // Time ack service waits to clear buffer and elect new proposer (1800 seconds ~ 30 mins)

	DefaultCheckpointerPollInterval = 5 * time.Minute
	DefaultSyncerPollInterval       = 1 * time.Minute
	DefaultNoACKPollInterval        = 1010 * time.Second
	DefaultClerkPollInterval        = 10 * time.Second
	DefaultSpanPollInterval         = 1 * time.Minute
	DefaultEnableSH                 = false
	DefaultSHStateSyncedInterval    = 15 * time.Minute
	DefaultSHStakeUpdateInterval    = 3 * time.Hour
	DefaultSHMaxDepthDuration       = time.Hour

	DefaultMainchainGasLimit = uint64(5000000)

	DefaultMainchainMaxGasPrice = 400000000000 // 400 Gwei

	DefaultBorChainID = "15001"

	DefaultLogsType = "json"
	DefaultChain    = MainChain

	DefaultTendermintNode = "tcp://localhost:26657"

	DefaultMainnetSeeds = "f4f605d60b8ffaaf15240564e58a81103510631c@159.203.9.164:26656,4fb1bc820088764a564d4f66bba1963d47d82329@44.232.55.71:26656,2eadba4be3ce47ac8db0a3538cb923b57b41c927@35.199.4.13:26656,3b23b20017a6f348d329c102ddc0088f0a10a444@35.221.13.28:26656,25f5f65a09c56e9f1d2d90618aa70cd358aa68da@35.230.116.151:26656"

	DefaultTestnetSeeds = "4cd60c1d76e44b05f7dfd8bab3f447b119e87042@54.147.31.250:26656,b18bbe1f3d8576f4b73d9b18976e71c65e839149@34.226.134.117:26656"

	secretFilePerm = 0600

	// Legacy value - DO NOT CHANGE
	// Maximum allowed event record data size
	LegacyMaxStateSyncSize = 100000

	// New max state sync size after hardfork
	MaxStateSyncSize = 30000

	// Default Open Collector Endpoint
	DefaultOpenCollectorEndpoint = "localhost:4317"
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.heimdallcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.heimdalld")
	MinBalance      = big.NewInt(100000000000000000) // aka 0.1 Ether
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.PrivKeyAminoName, nil)

	Logger = logger.NewTMLogger(logger.NewSyncWriter(os.Stdout))
}

// Configuration represents heimdall config
type Configuration struct {
	EthRPCUrl        string `mapstructure:"eth_rpc_url"`        // RPC endpoint for main chain
	BorRPCUrl        string `mapstructure:"bor_rpc_url"`        // RPC endpoint for bor chain
	TendermintRPCUrl string `mapstructure:"tendermint_rpc_url"` // tendemint node url

	EthRPCTimeout time.Duration `mapstructure:"eth_rpc_timeout"` // timeout for eth rpc
	BorRPCTimeout time.Duration `mapstructure:"bor_rpc_timeout"` // timeout for bor rpc

	AmqpURL           string `mapstructure:"amqp_url"`             // amqp url
	HeimdallServerURL string `mapstructure:"heimdall_rest_server"` // heimdall server url

	MainchainGasLimit uint64 `mapstructure:"main_chain_gas_limit"` // gas limit to mainchain transaction. eg....submit checkpoint.

	MainchainMaxGasPrice int64 `mapstructure:"main_chain_max_gas_price"` // max gas price to mainchain transaction. eg....submit checkpoint.

	// config related to bridge
	CheckpointerPollInterval time.Duration `mapstructure:"checkpoint_poll_interval"` // Poll interval for checkpointer service to send new checkpoints or missing ACK
	SyncerPollInterval       time.Duration `mapstructure:"syncer_poll_interval"`     // Poll interval for syncher service to sync for changes on main chain
	NoACKPollInterval        time.Duration `mapstructure:"noack_poll_interval"`      // Poll interval for ack service to send no-ack in case of no checkpoints
	ClerkPollInterval        time.Duration `mapstructure:"clerk_poll_interval"`
	SpanPollInterval         time.Duration `mapstructure:"span_poll_interval"`
	EnableSH                 bool          `mapstructure:"enable_self_heal"`         // Enable self healing
	SHStateSyncedInterval    time.Duration `mapstructure:"sh_state_synced_interval"` // Interval to self-heal StateSynced events if missing
	SHStakeUpdateInterval    time.Duration `mapstructure:"sh_stake_update_interval"` // Interval to self-heal StakeUpdate events if missing
	SHMaxDepthDuration       time.Duration `mapstructure:"sh_max_depth_duration"`    // Max duration that allows to suggest self-healing is not needed

	// wait time related options
	NoACKWaitTime time.Duration `mapstructure:"no_ack_wait_time"` // Time ack service waits to clear buffer and elect new proposer

	// Log related options
	LogsType       string `mapstructure:"logs_type"`        // if true, enable logging in json format
	LogsWriterFile string `mapstructure:"logs_writer_file"` // if given, Logs will be written to this file else os.Stdout

	// current chain - newSelectionAlgoHeight depends on this
	Chain string `mapstructure:"chain"`
}

var conf Configuration

// MainChainClient stores eth clie nt for Main chain Network
var mainChainClient *ethclient.Client
var mainRPCClient *rpc.Client

// MaticClient stores eth/rpc client for Matic Network
var maticClient *ethclient.Client
var maticRPCClient *rpc.Client

// private key object
var privObject secp256k1.PrivKeySecp256k1

var pubObject secp256k1.PubKeySecp256k1

// Logger stores global logger object
var Logger logger.Logger

// GenesisDoc contains the genesis file
var GenesisDoc tmTypes.GenesisDoc

var newSelectionAlgoHeight int64 = 0

var spanOverrideHeight int64 = 0

type ChainManagerAddressMigration struct {
	MaticTokenAddress     hmTypes.HeimdallAddress
	RootChainAddress      hmTypes.HeimdallAddress
	StakingManagerAddress hmTypes.HeimdallAddress
	SlashManagerAddress   hmTypes.HeimdallAddress
	StakingInfoAddress    hmTypes.HeimdallAddress
	StateSenderAddress    hmTypes.HeimdallAddress
}

var chainManagerAddressMigrations = map[string]map[int64]ChainManagerAddressMigration{
	MainChain:   {},
	MumbaiChain: {},
	"default":   {},
}

// Contracts
// var RootChain types.Contract
// var DepositManager types.Contract

// InitHeimdallConfig initializes with viper config (from heimdall configuration)
func InitHeimdallConfig(homeDir string) {
	if strings.Compare(homeDir, "") == 0 {
		// get home dir from viper
		homeDir = viper.GetString(HomeFlag)
	}

	// get heimdall config filepath from viper/cobra flag
	heimdallConfigFileFromFlag := viper.GetString(WithHeimdallConfigFlag)

	// init heimdall with changed config files
	InitHeimdallConfigWith(homeDir, heimdallConfigFileFromFlag)
}

// InitHeimdallConfigWith initializes passed heimdall/tendermint config files
func InitHeimdallConfigWith(homeDir string, heimdallConfigFileFromFLag string) {
	if strings.Compare(homeDir, "") == 0 {
		return
	}

	if strings.Compare(conf.BorRPCUrl, "") != 0 {
		return
	}

	// read configuration from the standard configuration file
	configDir := filepath.Join(homeDir, "config")
	heimdallViper := viper.New()
	heimdallViper.SetEnvPrefix("HEIMDALL")
	heimdallViper.AutomaticEnv()

	if heimdallConfigFileFromFLag == "" {
		heimdallViper.SetConfigName("heimdall-config") // name of config file (without extension)
		heimdallViper.AddConfigPath(configDir)         // call multiple times to add many search paths
	} else {
		heimdallViper.SetConfigFile(heimdallConfigFileFromFLag) // set config file explicitly
	}

	// Handle errors reading the config file
	if err := heimdallViper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	// unmarshal configuration from the standard configuration file
	if err := heimdallViper.UnmarshalExact(&conf); err != nil {
		log.Fatalln("Unable to unmarshall config", "Error", err)
	}

	//  if there is a file with overrides submitted via flags => read it an merge it with the alreadey read standard configuration
	if heimdallConfigFileFromFLag != "" {
		heimdallViperFromFlag := viper.New()
		heimdallViperFromFlag.SetConfigFile(heimdallConfigFileFromFLag) // set flag config file explicitly

		err := heimdallViperFromFlag.ReadInConfig()
		if err != nil { // Handle errors reading the config file sybmitted as a flag
			log.Fatalln("Unable to read config file submitted via flag", "Error", err)
		}

		var confFromFlag Configuration
		// unmarshal configuration from the configuration file submited as a flag
		if err = heimdallViperFromFlag.UnmarshalExact(&confFromFlag); err != nil {
			log.Fatalln("Unable to unmarshall config file submitted via flag", "Error", err)
		}

		conf.Merge(&confFromFlag)
	}

	// update configuration data with submitted flags
	if err := conf.UpdateWithFlags(viper.GetViper(), Logger); err != nil {
		log.Fatalln("Unable to read flag values. Check log for details.", "Error", err)
	}

	// perform check for json logging
	if conf.LogsType == "json" {
		Logger = logger.NewTMJSONLogger(logger.NewSyncWriter(GetLogsWriter(conf.LogsWriterFile)))
	} else {
		// default fallback
		Logger = logger.NewTMLogger(logger.NewSyncWriter(GetLogsWriter(conf.LogsWriterFile)))
	}

	// perform checks for timeout
	if conf.EthRPCTimeout == 0 {
		// fallback to default
		Logger.Debug("Invalid ETH RPC timeout provided, falling back to default value", "timeout", DefaultEthRPCTimeout)
		conf.EthRPCTimeout = DefaultEthRPCTimeout
	}

	if conf.BorRPCTimeout == 0 {
		// fallback to default
		Logger.Debug("Invalid BOR RPC timeout provided, falling back to default value", "timeout", DefaultBorRPCTimeout)
		conf.BorRPCTimeout = DefaultBorRPCTimeout
	}

	if conf.SHStateSyncedInterval == 0 {
		// fallback to default
		Logger.Debug("Invalid self-healing StateSynced interval provided, falling back to default value", "interval", DefaultSHStateSyncedInterval)
		conf.SHStateSyncedInterval = DefaultSHStateSyncedInterval
	}

	if conf.SHStakeUpdateInterval == 0 {
		// fallback to default
		Logger.Debug("Invalid self-healing StakeUpdate interval provided, falling back to default value", "interval", DefaultSHStakeUpdateInterval)
		conf.SHStakeUpdateInterval = DefaultSHStakeUpdateInterval
	}

	if conf.SHMaxDepthDuration == 0 {
		// fallback to default
		Logger.Debug("Invalid self-healing max depth duration provided, falling back to default value", "duration", DefaultSHMaxDepthDuration)
		conf.SHMaxDepthDuration = DefaultSHMaxDepthDuration
	}

	var err error
	if mainRPCClient, err = rpc.Dial(conf.EthRPCUrl); err != nil {
		log.Fatalln("Unable to dial via ethClient", "URL=", conf.EthRPCUrl, "chain=eth", "Error", err)
	}

	mainChainClient = ethclient.NewClient(mainRPCClient)

	if maticRPCClient, err = rpc.Dial(conf.BorRPCUrl); err != nil {
		log.Fatal(err)
	}

	maticClient = ethclient.NewClient(maticRPCClient)
	// Loading genesis doc
	genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(configDir, "genesis.json"))
	if err != nil {
		log.Fatal(err)
	}

	GenesisDoc = *genDoc

	// load pv file, unmarshall and set to privObject
	err = file.PermCheck(file.Rootify("priv_validator_key.json", configDir), secretFilePerm)
	if err != nil {
		Logger.Error(err.Error())
	}

	privVal := privval.LoadFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(configDir, "priv_validator_key.json"))
	cdc.MustUnmarshalBinaryBare(privVal.Key.PrivKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)

	switch conf.Chain {
	case MainChain:
		newSelectionAlgoHeight = 375300
		spanOverrideHeight = 8664000
	case MumbaiChain:
		newSelectionAlgoHeight = 282500
		spanOverrideHeight = 10205000
	default:
		newSelectionAlgoHeight = 0
		spanOverrideHeight = 0
	}
}

// GetDefaultHeimdallConfig returns configration with default params
func GetDefaultHeimdallConfig() Configuration {
	return Configuration{
		EthRPCUrl:        DefaultMainRPCUrl,
		BorRPCUrl:        DefaultBorRPCUrl,
		TendermintRPCUrl: DefaultTendermintNodeURL,

		EthRPCTimeout: DefaultEthRPCTimeout,
		BorRPCTimeout: DefaultBorRPCTimeout,

		AmqpURL:           DefaultAmqpURL,
		HeimdallServerURL: DefaultHeimdallServerURL,

		MainchainGasLimit: DefaultMainchainGasLimit,

		MainchainMaxGasPrice: DefaultMainchainMaxGasPrice,

		CheckpointerPollInterval: DefaultCheckpointerPollInterval,
		SyncerPollInterval:       DefaultSyncerPollInterval,
		NoACKPollInterval:        DefaultNoACKPollInterval,
		ClerkPollInterval:        DefaultClerkPollInterval,
		SpanPollInterval:         DefaultSpanPollInterval,
		EnableSH:                 DefaultEnableSH,
		SHStateSyncedInterval:    DefaultSHStateSyncedInterval,
		SHStakeUpdateInterval:    DefaultSHStakeUpdateInterval,
		SHMaxDepthDuration:       DefaultSHMaxDepthDuration,

		NoACKWaitTime: NoACKWaitTime,

		LogsType:       DefaultLogsType,
		Chain:          DefaultChain,
		LogsWriterFile: "", // default to stdout
	}
}

// GetConfig returns cached configuration object
func GetConfig() Configuration {
	return conf
}

func GetGenesisDoc() tmTypes.GenesisDoc {
	return GenesisDoc
}

// TEST PURPOSE ONLY
// SetTestConfig sets test configuration
func SetTestConfig(_conf Configuration) {
	conf = _conf
}

//
// Get main/matic clients
//

// GetMainChainRPCClient returns main chain RPC client
func GetMainChainRPCClient() *rpc.Client {
	return mainRPCClient
}

// GetMainClient returns main chain's eth client
func GetMainClient() *ethclient.Client {
	return mainChainClient
}

// GetMaticClient returns matic's eth client
func GetMaticClient() *ethclient.Client {
	return maticClient
}

// GetMaticRPCClient returns matic's RPC client
func GetMaticRPCClient() *rpc.Client {
	return maticRPCClient
}

// GetPrivKey returns priv key object
func GetPrivKey() secp256k1.PrivKeySecp256k1 {
	return privObject
}

// GetECDSAPrivKey return ecdsa private key
func GetECDSAPrivKey() *ecdsa.PrivateKey {
	// get priv key
	pkObject := GetPrivKey()

	// create ecdsa private key
	ecdsaPrivateKey, _ := ethCrypto.ToECDSA(pkObject[:])

	return ecdsaPrivateKey
}

// GetPubKey returns pub key object
func GetPubKey() secp256k1.PubKeySecp256k1 {
	return pubObject
}

// GetAddress returns address object
func GetAddress() []byte {
	return GetPubKey().Address().Bytes()
}

// GetValidChains returns all the valid chains
func GetValidChains() []string {
	return []string{"mainnet", "mumbai", "local"}
}

// GetNewSelectionAlgoHeight returns newSelectionAlgoHeight
func GetNewSelectionAlgoHeight() int64 {
	return newSelectionAlgoHeight
}

// GetSpanOverrideHeight returns spanOverrideHeight
func GetSpanOverrideHeight() int64 {
	return spanOverrideHeight
}

func GetChainManagerAddressMigration(blockNum int64) (ChainManagerAddressMigration, bool) {
	chainMigration := chainManagerAddressMigrations[conf.Chain]
	if chainMigration == nil {
		chainMigration = chainManagerAddressMigrations["default"]
	}

	result, found := chainMigration[blockNum]

	return result, found
}

// DecorateWithHeimdallFlags adds persistent flags for heimdall-config and bind flags with command
func DecorateWithHeimdallFlags(cmd *cobra.Command, v *viper.Viper, loggerInstance logger.Logger, caller string) {
	// add with-heimdall-config flag
	cmd.PersistentFlags().String(
		WithHeimdallConfigFlag,
		"",
		"Override of Heimdall config file (default <home>/config/heimdall-config.json)",
	)

	if err := v.BindPFlag(WithHeimdallConfigFlag, cmd.PersistentFlags().Lookup(WithHeimdallConfigFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, WithHeimdallConfigFlag), "Error", err)
	}

	// add MainRPCUrlFlag flag
	cmd.PersistentFlags().String(
		MainRPCUrlFlag,
		"",
		"Set RPC endpoint for ethereum chain",
	)

	if err := v.BindPFlag(MainRPCUrlFlag, cmd.PersistentFlags().Lookup(MainRPCUrlFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainRPCUrlFlag), "Error", err)
	}

	// add BorRPCUrlFlag flag
	cmd.PersistentFlags().String(
		BorRPCUrlFlag,
		"",
		"Set RPC endpoint for bor chain",
	)

	if err := v.BindPFlag(BorRPCUrlFlag, cmd.PersistentFlags().Lookup(BorRPCUrlFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, BorRPCUrlFlag), "Error", err)
	}

	// add TendermintNodeURLFlag flag
	cmd.PersistentFlags().String(
		TendermintNodeURLFlag,
		"",
		"Set RPC endpoint for tendermint",
	)

	if err := v.BindPFlag(TendermintNodeURLFlag, cmd.PersistentFlags().Lookup(TendermintNodeURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, TendermintNodeURLFlag), "Error", err)
	}

	// add HeimdallServerURLFlag flag
	cmd.PersistentFlags().String(
		HeimdallServerURLFlag,
		"",
		"Set Heimdall REST server endpoint",
	)

	if err := v.BindPFlag(HeimdallServerURLFlag, cmd.PersistentFlags().Lookup(HeimdallServerURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, HeimdallServerURLFlag), "Error", err)
	}

	// add AmqpURLFlag flag
	cmd.PersistentFlags().String(
		AmqpURLFlag,
		"",
		"Set AMQP endpoint",
	)

	if err := v.BindPFlag(AmqpURLFlag, cmd.PersistentFlags().Lookup(AmqpURLFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, AmqpURLFlag), "Error", err)
	}

	// add CheckpointerPollIntervalFlag flag
	cmd.PersistentFlags().String(
		CheckpointerPollIntervalFlag,
		"",
		"Set check point pull interval",
	)

	if err := v.BindPFlag(CheckpointerPollIntervalFlag, cmd.PersistentFlags().Lookup(CheckpointerPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, CheckpointerPollIntervalFlag), "Error", err)
	}

	// add SyncerPollIntervalFlag flag
	cmd.PersistentFlags().String(
		SyncerPollIntervalFlag,
		"",
		"Set syncer pull interval",
	)

	if err := v.BindPFlag(SyncerPollIntervalFlag, cmd.PersistentFlags().Lookup(SyncerPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, SyncerPollIntervalFlag), "Error", err)
	}

	// add NoACKPollIntervalFlag flag
	cmd.PersistentFlags().String(
		NoACKPollIntervalFlag,
		"",
		"Set no acknowledge pull interval",
	)

	if err := v.BindPFlag(NoACKPollIntervalFlag, cmd.PersistentFlags().Lookup(NoACKPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, NoACKPollIntervalFlag), "Error", err)
	}

	// add ClerkPollIntervalFlag flag
	cmd.PersistentFlags().String(
		ClerkPollIntervalFlag,
		"",
		"Set clerk pull interval",
	)

	if err := v.BindPFlag(ClerkPollIntervalFlag, cmd.PersistentFlags().Lookup(ClerkPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ClerkPollIntervalFlag), "Error", err)
	}

	// add SpanPollIntervalFlag flag
	cmd.PersistentFlags().String(
		SpanPollIntervalFlag,
		"",
		"Set span pull interval",
	)

	if err := v.BindPFlag(SpanPollIntervalFlag, cmd.PersistentFlags().Lookup(SpanPollIntervalFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, SpanPollIntervalFlag), "Error", err)
	}

	// add MainchainGasLimitFlag flag
	cmd.PersistentFlags().Uint64(
		MainchainGasLimitFlag,
		0,
		"Set main chain gas limti",
	)

	if err := v.BindPFlag(MainchainGasLimitFlag, cmd.PersistentFlags().Lookup(MainchainGasLimitFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainchainGasLimitFlag), "Error", err)
	}

	// add MainchainMaxGasPriceFlag flag
	cmd.PersistentFlags().Int64(
		MainchainMaxGasPriceFlag,
		0,
		"Set main chain max gas limti",
	)

	if err := v.BindPFlag(MainchainMaxGasPriceFlag, cmd.PersistentFlags().Lookup(MainchainMaxGasPriceFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, MainchainMaxGasPriceFlag), "Error", err)
	}

	// add NoACKWaitTimeFlag flag
	cmd.PersistentFlags().String(
		NoACKWaitTimeFlag,
		"",
		"Set time ack service waits to clear buffer and elect new proposer",
	)

	if err := v.BindPFlag(NoACKWaitTimeFlag, cmd.PersistentFlags().Lookup(NoACKWaitTimeFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, NoACKWaitTimeFlag), "Error", err)
	}

	// add chain flag
	cmd.PersistentFlags().String(
		ChainFlag,
		"",
		fmt.Sprintf("Set one of the chains: [%s]", strings.Join(GetValidChains(), ",")),
	)

	if err := v.BindPFlag(ChainFlag, cmd.PersistentFlags().Lookup(ChainFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, ChainFlag), "Error", err)
	}

	// add logsWriterFile flag
	cmd.PersistentFlags().String(
		LogsWriterFileFlag,
		"",
		"Set logs writer file, Default is os.Stdout",
	)

	if err := v.BindPFlag(LogsWriterFileFlag, cmd.PersistentFlags().Lookup(LogsWriterFileFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", caller, LogsWriterFileFlag), "Error", err)
	}
}

func (c *Configuration) UpdateWithFlags(v *viper.Viper, loggerInstance logger.Logger) error {
	const logErrMsg = "Unable to read flag."

	// get endpoint for ethereum chain from viper/cobra
	stringConfgValue := v.GetString(MainRPCUrlFlag)
	if stringConfgValue != "" {
		c.EthRPCUrl = stringConfgValue
	}

	// get endpoint for bor chain from viper/cobra
	stringConfgValue = v.GetString(BorRPCUrlFlag)
	if stringConfgValue != "" {
		c.BorRPCUrl = stringConfgValue
	}

	// get endpoint for tendermint from viper/cobra
	stringConfgValue = v.GetString(TendermintNodeURLFlag)
	if stringConfgValue != "" {
		c.TendermintRPCUrl = stringConfgValue
	}

	// get endpoint for tendermint from viper/cobra
	stringConfgValue = v.GetString(AmqpURLFlag)
	if stringConfgValue != "" {
		c.AmqpURL = stringConfgValue
	}

	// get Heimdall REST server endpoint from viper/cobra
	stringConfgValue = v.GetString(HeimdallServerURLFlag)
	if stringConfgValue != "" {
		c.HeimdallServerURL = stringConfgValue
	}

	// need this error for parsing Duration values
	var err error

	// get check point pull interval from viper/cobra
	stringConfgValue = v.GetString(CheckpointerPollIntervalFlag)
	if stringConfgValue != "" {
		if c.CheckpointerPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", CheckpointerPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get syncer pull interval from viper/cobra
	stringConfgValue = v.GetString(SyncerPollIntervalFlag)
	if stringConfgValue != "" {
		if c.SyncerPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", SyncerPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get poll interval for ack service to send no-ack in case of no checkpoints from viper/cobra
	stringConfgValue = v.GetString(NoACKPollIntervalFlag)
	if stringConfgValue != "" {
		if c.NoACKPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", NoACKPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get clerk poll interval from viper/cobra
	stringConfgValue = v.GetString(ClerkPollIntervalFlag)
	if stringConfgValue != "" {
		if c.ClerkPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", ClerkPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get span poll interval from viper/cobra
	stringConfgValue = v.GetString(SpanPollIntervalFlag)
	if stringConfgValue != "" {
		if c.SpanPollInterval, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", SpanPollIntervalFlag, "Error", err)
			return err
		}
	}

	// get time that ack service waits to clear buffer and elect new proposer from viper/cobra
	stringConfgValue = v.GetString(NoACKWaitTimeFlag)
	if stringConfgValue != "" {
		if c.NoACKWaitTime, err = time.ParseDuration(stringConfgValue); err != nil {
			loggerInstance.Error(logErrMsg, "Flag", NoACKWaitTimeFlag, "Error", err)
			return err
		}
	}

	// get mainchain gas limit from viper/cobra
	uint64ConfgValue := v.GetUint64(MainchainGasLimitFlag)
	if uint64ConfgValue != 0 {
		c.MainchainGasLimit = uint64ConfgValue
	}

	// get mainchain max gas price from viper/cobra. if it is greater then  zero => set it as configuration parameter
	int64ConfgValue := v.GetInt64(MainchainMaxGasPriceFlag)
	if int64ConfgValue > 0 {
		c.MainchainMaxGasPrice = int64ConfgValue
	}

	// get chain from viper/cobra flag
	stringConfgValue = v.GetString(ChainFlag)
	if stringConfgValue != "" {
		c.Chain = stringConfgValue
	}

	stringConfgValue = v.GetString(LogsWriterFileFlag)
	if stringConfgValue != "" {
		c.LogsWriterFile = stringConfgValue
	}

	return nil
}

func (c *Configuration) Merge(cc *Configuration) {
	if cc.EthRPCUrl != "" {
		c.EthRPCUrl = cc.EthRPCUrl
	}

	if cc.BorRPCUrl != "" {
		c.BorRPCUrl = cc.BorRPCUrl
	}

	if cc.TendermintRPCUrl != "" {
		c.TendermintRPCUrl = cc.TendermintRPCUrl
	}

	if cc.AmqpURL != "" {
		c.AmqpURL = cc.AmqpURL
	}

	if cc.HeimdallServerURL != "" {
		c.HeimdallServerURL = cc.HeimdallServerURL
	}

	if cc.MainchainGasLimit != 0 {
		c.MainchainGasLimit = cc.MainchainGasLimit
	}

	if cc.MainchainMaxGasPrice != 0 {
		c.MainchainMaxGasPrice = cc.MainchainMaxGasPrice
	}

	if cc.CheckpointerPollInterval != 0 {
		c.CheckpointerPollInterval = cc.CheckpointerPollInterval
	}

	if cc.SyncerPollInterval != 0 {
		c.SyncerPollInterval = cc.SyncerPollInterval
	}

	if cc.NoACKPollInterval != 0 {
		c.NoACKPollInterval = cc.NoACKPollInterval
	}

	if cc.ClerkPollInterval != 0 {
		c.ClerkPollInterval = cc.ClerkPollInterval
	}

	if cc.SpanPollInterval != 0 {
		c.SpanPollInterval = cc.SpanPollInterval
	}

	if cc.NoACKWaitTime != 0 {
		c.NoACKWaitTime = cc.NoACKWaitTime
	}

	if cc.Chain != "" {
		c.Chain = cc.Chain
	}

	if cc.LogsWriterFile != "" {
		c.LogsWriterFile = cc.LogsWriterFile
	}
}

// DecorateWithTendermintFlags creates tendermint flags for desired command and bind them to viper
func DecorateWithTendermintFlags(cmd *cobra.Command, v *viper.Viper, loggerInstance logger.Logger, message string) {
	// add seeds flag
	cmd.PersistentFlags().String(
		SeedsFlag,
		"",
		"Override seeds",
	)

	if err := v.BindPFlag(SeedsFlag, cmd.PersistentFlags().Lookup(SeedsFlag)); err != nil {
		loggerInstance.Error(fmt.Sprintf("%v | BindPFlag | %v", message, SeedsFlag), "Error", err)
	}
}

// UpdateTendermintConfig updates tenedermint config with flags and default values if needed
func UpdateTendermintConfig(tendermintConfig *cfg.Config, v *viper.Viper) {
	// update tendermintConfig.P2P.Seeds
	seedsFlagValue := v.GetString(SeedsFlag)
	if seedsFlagValue != "" {
		tendermintConfig.P2P.Seeds = seedsFlagValue
	}

	if tendermintConfig.P2P.Seeds == "" {
		switch conf.Chain {
		case MainChain:
			tendermintConfig.P2P.Seeds = DefaultMainnetSeeds
		case MumbaiChain:
			tendermintConfig.P2P.Seeds = DefaultTestnetSeeds
		}
	}
}

func GetLogsWriter(logsWriterFile string) io.Writer {
	if logsWriterFile != "" {
		logWriter, err := os.OpenFile(logsWriterFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log writer file: %v", err)
		}

		return logWriter
	} else {
		return os.Stdout
	}
}
