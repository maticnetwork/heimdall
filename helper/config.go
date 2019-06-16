package helper

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	logger "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"

	"math/big"

	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
)

const (
	NodeFlag               = "node"
	WithHeimdallConfigFlag = "with-heimdall-config"
	HomeFlag               = "home"
	FlagClientHome         = "home-client"

	// --- TODO Move these to common client flags
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

	// Variables below to be used while init
	MainRPCUrl                      = "https://ropsten.infura.io"
	MaticRPCUrl                     = "https://testnet.matic.network"
	NoACKWaitTime                   = time.Second * 1800 // Time ack service waits to clear buffer and elect new proposer (1800 seconds ~ 30 mins)
	CheckpointBufferTime            = time.Second * 1000 // Time checkpoint is allowed to stay in buffer (1000 seconds ~ 17 mins)
	DefaultCheckpointerPollInterval = 60 * 1000          // 1 minute in milliseconds
	DefaultSyncerPollInterval       = 30 * 1000          // 0.5 seconds in milliseconds
	DefaultNoACKPollInterval        = 1010 * time.Second
	DefaultCheckpointLength         = 256   // checkpoint number 	 with 0, so length = defaultCheckpointLength -1
	MaxCheckpointLength             = 1024  // max blocks in one checkpoint
	DefaultChildBlockInterval       = 10000 // difference between 2 indexes of header blocks
	ConfirmationBlocks              = 6
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
	MainRPCUrl          string `json:"mainRPCUrl"`          // RPC endpoint for main chain
	MaticRPCUrl         string `json:"maticRPCUrl"`         // RPC endpoint for matic chain
	StakeManagerAddress string `json:"stakeManagerAddress"` // Stake manager address on main chain
	RootchainAddress    string `json:"rootchainAddress"`    // Rootchain contract address on main chain
	ChildBlockInterval  uint64 `json:"childBlockInterval"`  // Difference between header index of 2 child blocks submitted on main chain
	// config related to bridge
	CheckpointerPollInterval int           `json:"checkpointerPollInterval"` // Poll interval for checkpointer service to send new checkpoints or missing ACK
	SyncerPollInterval       int           `json:"syncerPollInterval"`       // Poll interval for syncher service to sync for changes on main chain
	NoACKPollInterval        time.Duration `json:"noackPollInterval"`        // Poll interval for ack service to send no-ack in case of no checkpoints
	// checkpoint length related options
	AvgCheckpointLength uint64 `json:"avgCheckpointLength"` // Average number of blocks checkpoint would contain
	MaxCheckpointLength uint64 `json:"maxCheckpointLength"` // Maximium number of blocks checkpoint would contain
	// wait time related options
	NoACKWaitTime        time.Duration `json:"noackWaitTime"`        // Time ack service waits to clear buffer and elect new proposer
	CheckpointBufferTime time.Duration `json:"checkpointBufferTime"` // Time checkpoint is allowed to stay in buffer

	ConfirmationBlocks uint64 `json:"confirmationBlocks"` // Number of blocks for confirmation
}

var conf Configuration

// MainChainClient stores eth client for Main chain Network
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

	if mainRPCClient, err = rpc.Dial(conf.MainRPCUrl); err != nil {
		Logger.Error("Error while creating matic chain RPC client", "error", err)
		log.Fatal(err)
	}
	mainChainClient = ethclient.NewClient(mainRPCClient)

	if maticRPCClient, err = rpc.Dial(conf.MaticRPCUrl); err != nil {
		Logger.Error("Error while creating matic chain RPC client", "error", err)
		log.Fatal(err)
	}
	maticClient = ethclient.NewClient(maticRPCClient)

	// load pv file, unmarshall and set to privObject
	privVal := privval.LoadFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(configDir, "priv_validator_key.json"))
	cdc.MustUnmarshalBinaryBare(privVal.Key.PrivKey.Bytes(), &privObject)
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
	rootChainInstance, err := rootchain.NewRootchain(GetRootChainAddress(), mainChainClient)
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
	stakeManagerInstance, err := stakemanager.NewStakemanager(GetStakeManagerAddress(), mainChainClient)
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

// GetPubKey returns pub key object
func GetPubKey() secp256k1.PubKeySecp256k1 {
	return pubObject
}

// GetAddress returns address object
func GetAddress() []byte {
	return GetPubKey().Address().Bytes()
}
