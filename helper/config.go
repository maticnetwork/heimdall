package helper

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/heimdall/file"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	logger "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"

	tmTypes "github.com/tendermint/tendermint/types"
)

const (
	NodeFlag               = "node"
	WithHeimdallConfigFlag = "with-heimdall-config"
	HomeFlag               = "home"
	FlagClientHome         = "home-client"

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

	DefaultMainRPCUrl = "http://localhost:9545"
	DefaultBorRPCUrl  = "http://localhost:8545"

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

	DefaultMainchainGasLimit = uint64(5000000)

	DefaultMainchainMaxGasPrice = 400000000000 // 400 Gwei

	DefaultBorChainID string = "15001"

	secretFilePerm = 0600

	// Legacy value - DO NOT CHANGE
	// Maximum allowed event record data size
	LegacyMaxStateSyncSize = 100000

	// New max state sync size after hardfork
	MaxStateSyncSize = 30000
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

	// wait time related options
	NoACKWaitTime time.Duration `mapstructure:"no_ack_wait_time"` // Time ack service waits to clear buffer and elect new proposer
}

var conf Configuration

// MainChainClient stores eth clie nt for Main chain Network
var mainChainClient *ethclient.Client
var mainRPCClient *rpc.Client

// MaticClient stores eth/rpc client for Matic Network
var maticClient *ethclient.Client
var maticRPCClient *rpc.Client

var maticEthClient *eth.EthAPIBackend

// private key object
var privObject secp256k1.PrivKeySecp256k1

var pubObject secp256k1.PubKeySecp256k1

// Logger stores global logger object
var Logger logger.Logger

// GenesisDoc contains the genesis file
var GenesisDoc tmTypes.GenesisDoc

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
	heimdallConfigFilePath := viper.GetString(WithHeimdallConfigFlag)

	// init heimdall with changed config files
	InitHeimdallConfigWith(homeDir, heimdallConfigFilePath)
}

// InitHeimdallConfigWith initializes passed heimdall/tendermint config files
func InitHeimdallConfigWith(homeDir string, heimdallConfigFilePath string) {
	if strings.Compare(homeDir, "") == 0 {
		return
	}

	if strings.Compare(conf.BorRPCUrl, "") != 0 {
		return
	}

	configDir := filepath.Join(homeDir, "config")

	heimdallViper := viper.New()
	heimdallViper.SetEnvPrefix("HEIMDALL")
	heimdallViper.AutomaticEnv()
	if heimdallConfigFilePath == "" {
		heimdallViper.SetConfigName("heimdall-config") // name of config file (without extension)
		heimdallViper.AddConfigPath(configDir)         // call multiple times to add many search paths
	} else {
		heimdallViper.SetConfigFile(heimdallConfigFilePath) // set config file explicitly
	}

	err := heimdallViper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	if err = heimdallViper.UnmarshalExact(&conf); err != nil {
		log.Fatalln("Unable to unmarshall config", "Error", err)
	}

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
}

// GetDefaultHeimdallConfig returns configration with default params
func GetDefaultHeimdallConfig() Configuration {
	return Configuration{
		EthRPCUrl:        DefaultMainRPCUrl,
		BorRPCUrl:        DefaultBorRPCUrl,
		TendermintRPCUrl: DefaultTendermintNodeURL,

		AmqpURL:           DefaultAmqpURL,
		HeimdallServerURL: DefaultHeimdallServerURL,

		MainchainGasLimit: DefaultMainchainGasLimit,

		MainchainMaxGasPrice: DefaultMainchainMaxGasPrice,

		CheckpointerPollInterval: DefaultCheckpointerPollInterval,
		SyncerPollInterval:       DefaultSyncerPollInterval,
		NoACKPollInterval:        DefaultNoACKPollInterval,
		ClerkPollInterval:        DefaultClerkPollInterval,
		SpanPollInterval:         DefaultSpanPollInterval,

		NoACKWaitTime: NoACKWaitTime,
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

// GetMaticEthClient returns matic's Eth client
func GetMaticEthClient() *eth.EthAPIBackend {
	return maticEthClient
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
