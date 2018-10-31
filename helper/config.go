package helper

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"

	"os"

	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	logger "github.com/tendermint/tendermint/libs/log"
)

const (
	WithHeimdallConfigFlag = "with-heimdall-config"
	HomeFlag               = "home"
)

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)
	Logger = logger.NewTMLogger(logger.NewSyncWriter(os.Stdout))
}

// Configuration represents heimdall config
type Configuration struct {
	MainRPCUrl          string `mapstructure:"main_rpcurl"`
	MaticRPCUrl         string `mapstructure:"matic_rpcurl"`
	StakeManagerAddress string `mapstructure:"stakemanager_address"`
	RootchainAddress    string `mapstructure:"rootchain_address"`

	// Tendermint endpoint
	TendermintEndpoint string `mapstructure:"tendermint_endpoint"`
}

var conf Configuration

// MainChainClient stores eth client for Main chain Network
var mainChainClient *ethclient.Client

// MaticClient stores eth/rpc client for Matic Network
var maticClient *ethclient.Client
var maticRPCClient *rpc.Client

// private key object
var privObject secp256k1.PrivKeySecp256k1
var pubObject secp256k1.PubKeySecp256k1

// Logger stores global logger object
var Logger logger.Logger

// InitHeimdallConfig initializes with viper config (from heimdall configuration)
func InitHeimdallConfig() {
	// get home dir from viper
	homeDir := viper.GetString(HomeFlag)

	// get heimdall config filepath from viper/cobra flag
	heimdallConfigFilePath := viper.GetString(WithHeimdallConfigFlag)

	// init heimdall with changed config files
	InitHeimdallConfigWith(homeDir, heimdallConfigFilePath)
}

// InitHeimdallConfigWith initializes passed heimdall/tendermint config files
func InitHeimdallConfigWith(homeDir string, heimdallConfigFilePath string) {
	if homeDir == "" {
		return
	}

	if strings.Compare(conf.MaticRPCUrl, "") != 0 {
		return
	}

	configDir := filepath.Join(homeDir, "config")
	Logger.Info("Initializing tendermint configurations", "configDir", configDir)

	heimdallViper := viper.New()
	if heimdallConfigFilePath == "" {
		heimdallViper.SetConfigName("heimdall-config") // name of config file (without extension)
		heimdallViper.AddConfigPath(configDir)         // call multiple times to add many search paths
		Logger.Info("Loading heimdall configurations", "file", filepath.Join(configDir, "heimdall-config.json"))
	} else {
		heimdallViper.SetConfigFile(heimdallConfigFilePath) // set config file explicitly
		Logger.Info("Loading heimdall configurations", "file", heimdallConfigFilePath)
	}

	err := heimdallViper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}
	if err = heimdallViper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	rpc.Dial(conf.MainRPCUrl)

	// setup eth client
	if mainChainClient, err = ethclient.Dial(conf.MainRPCUrl); err != nil {
		Logger.Error("Error while creating main chain client", "error", err)
		log.Fatal(err)
	}

	if maticRPCClient, err = rpc.Dial(conf.MaticRPCUrl); err != nil {
		Logger.Error("Error while creating matic chain RPC client", "error", err)
		log.Fatal(err)
	}
	maticClient = ethclient.NewClient(maticRPCClient)

	// load pv file, unmarshall and set to privObject
	privVal := privval.LoadFilePV(filepath.Join(configDir, "priv_validator.json"))
	cdc.MustUnmarshalBinaryBare(privVal.PrivKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)
}

// GetConfig returns cached configuration object
func GetConfig() Configuration {
	return conf
}

//
// Root chain
//

func GetRootChainAddress() common.Address {
	return common.HexToAddress(GetConfig().RootchainAddress)
}

func GetRootChainInstance() (*rootchain.Rootchain, error) {
	rootChainInstance, err := rootchain.NewRootchain(common.HexToAddress(GetConfig().RootchainAddress), mainChainClient)
	if err != nil {
		Logger.Error("Unable to create root chain instance", "error", err)
	}

	return rootChainInstance, err
}

func GetRootChainABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(rootchain.RootchainABI))
}

//
// Stake manager
//

func GetStakeManagerAddress() common.Address {
	return common.HexToAddress(GetConfig().StakeManagerAddress)
}

func GetStakeManagerInstance() (*stakemanager.Stakemanager, error) {
	stakeManagerInstance, err := stakemanager.NewStakemanager(common.HexToAddress(GetConfig().StakeManagerAddress), mainChainClient)
	if err != nil {
		Logger.Error("Unable to create stakemanager instance", "error", err)
	}

	return stakeManagerInstance, err
}

func GetStakeManagerABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(stakemanager.StakemanagerABI))
}

//
// Get main/matic clients
//

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

// GetPubKey returns pub key object
func GetPubKey() secp256k1.PubKeySecp256k1 {
	return pubObject
}
