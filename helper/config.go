package helper

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/types/common"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/maticnetwork/bor/eth"
	"github.com/maticnetwork/bor/ethclient"
	"github.com/maticnetwork/bor/rpc"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	logger "github.com/tendermint/tendermint/libs/log"
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

	DefaultBorChainID string = "15001"

	// secretFilePerm = 0600
)

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.heimdallcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.heimdalld")
	MinBalance      = big.NewInt(100000000000000000) // aka 0.1 Ether
)

var cdc = amino.NewCodec()

func init() {
	//interfaceRegistery := codectypes.NewInterfaceRegistry()

	//interfaceRegistery.RegisterInterface("pubKey", (*secp256k1.PubKey{})(nil))
	//cdc.RegisterConcrete(secp256k1.PubKey{}, secp256k1.PubKeyName, nil)
	//cdc.RegisterConcrete(secp256k1.PrivKey{}, secp256k1.PrivKeyName, nil)
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

// MainChainClient stores eth client for Main chain Network
var mainChainClient *ethclient.Client
var mainRPCClient *rpc.Client

// MaticClient stores eth/rpc client for Matic Network
var maticClient *ethclient.Client
var maticRPCClient *rpc.Client

var maticEthClient *eth.EthAPIBackend

// private key object
var FilePV *privval.FilePV

// Logger stores global logger object
var Logger logger.Logger

// GenesisDoc contains the genesis file
var GenesisDoc tmTypes.GenesisDoc

// Contracts
// var RootChain types.Contract
// var DepositManager types.Contract

// InitHeimdallConfig initializes passed heimdall/tendermint config files
func InitHeimdallConfig() error {
	rootDir := viper.GetString(flags.FlagHome)
	configDir := filepath.Join(rootDir, "config")

	heimdallConfigFilePath := filepath.Join(configDir, "heimdall-config.toml")
	if _, err := os.Stat(heimdallConfigFilePath); os.IsNotExist(err) {
		hc := GetDefaultHeimdallConfig()
		WriteConfigFile(heimdallConfigFilePath, &hc)
	}

	// create new viper and
	configViper := viper.New()
	configViper.SetConfigType("toml")
	configViper.SetConfigName("heimdall-config")
	configViper.AddConfigPath(configDir) // set config file explicitly

	err := configViper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		return err
	}

	if err = configViper.UnmarshalExact(&conf); err != nil {
		return fmt.Errorf("unable to unmarshall config %v", err)
	}

	if mainRPCClient, err = rpc.Dial(conf.EthRPCUrl); err != nil {
		return fmt.Errorf("unable to dial via ethClient. URL=%s, chain=eth, error=%v", conf.EthRPCUrl, err)
	}

	mainChainClient = ethclient.NewClient(mainRPCClient)
	if maticRPCClient, err = rpc.Dial(conf.BorRPCUrl); err != nil {
		return err
	}

	maticClient = ethclient.NewClient(maticRPCClient)

	// Loading genesis doc
	genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(configDir, "genesis.json"))
	if err != nil {
		return nil
	}
	GenesisDoc = *genDoc

	var dataDir = filepath.Join(rootDir, "data")

	FilePV = privval.LoadFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(dataDir, "priv_validator_state.json"))

	return nil
}

// GetDefaultHeimdallConfig returns configuration with default params
func GetDefaultHeimdallConfig() Configuration {
	return Configuration{
		EthRPCUrl:        DefaultMainRPCUrl,
		BorRPCUrl:        DefaultBorRPCUrl,
		TendermintRPCUrl: DefaultTendermintNodeURL,

		AmqpURL:           DefaultAmqpURL,
		HeimdallServerURL: DefaultHeimdallServerURL,

		MainchainGasLimit: DefaultMainchainGasLimit,

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
func GetPrivKey() secp256k1.PrivKey {
	return FilePV.Key.PrivKey.Bytes()
}

func GetPubKeyForCosmos() cryptotypes.PubKey {
	return hmCommon.CosmosCryptoPubKey(GetPubKey())
}

// GetECDSAPrivKey return ecdsa private key
//func GetECDSAPrivKey() *ecdsa.PrivateKey {
//	// get priv key
//	pkObject := GetPrivKey()
//
//	// create ecdsa private key
//	ecdsaPrivateKey, _ := ethCrypto.ToECDSA(pkObject[:])
//	return ecdsaPrivateKey
//}

// GetPubKey returns pub key object
func GetPubKey() secp256k1.PubKey {
	return FilePV.Key.PubKey.Bytes()
}

//func GetCryptoPrivKey() cryptotypes.PrivKey {
//	return FilePV.Key.PrivKey.
//}

// GetAddress returns address object
func GetAddress() []byte {
	addr, _ := sdk.AccAddressFromHex(viper.GetString("account-address"))
	return addr
}

// GetAddressStr returns address string object
func GetAddressStr() string {
	addr, _ := sdk.AccAddressFromHex(viper.GetString("account-address"))
	return addr.String()
}
