package helper

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	logger "github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"

	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"

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

	MainRPCUrl  = "https://ropsten.infura.io"
	MaticRPCUrl = "https://testnet2.matic.network"

	// Services

	// DefaultAmqpURL represents default AMQP url
	DefaultAmqpURL           = "amqp://guest:guest@localhost:5672/"
	DefaultHeimdallServerURL = "http://0.0.0.0:1317"
	DefaultTendermintNodeURL = "http://0.0.0.0:26657"

	NoACKWaitTime                   = time.Second * 1800 // Time ack service waits to clear buffer and elect new proposer (1800 seconds ~ 30 mins)
	CheckpointBufferTime            = time.Second * 1000 // Time checkpoint is allowed to stay in buffer (1000 seconds ~ 17 mins)
	DefaultCheckpointerPollInterval = 60 * 1000          // 1 minute in milliseconds
	DefaultSyncerPollInterval       = 30 * 1000          // 0.5 seconds in milliseconds
	DefaultNoACKPollInterval        = 1010 * time.Second
	DefaultCheckpointLength         = 256   // checkpoint number 	 with 0, so length = defaultCheckpointLength -1
	MaxCheckpointLength             = 1024  // max blocks in one checkpoint
	DefaultChildBlockInterval       = 10000 // difference between 2 indexes of header blocks
	ConfirmationBlocks              = 6

	DefaultBorChainID           = 15001
	DefaultStateReceiverAddress = "0000000000000000000000000000000000001001"
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
	MainRPCUrl  string `json:"mainRPCUrl"`  // RPC endpoint for main chain
	MaticRPCUrl string `json:"maticRPCUrl"` // RPC endpoint for matic chain

	AmqpURL           string `json:"amqpURL"`           // amqp url
	HeimdallServerURL string `json:"heimdallServerURL"` // heimdall server url
	TendermintNodeURL string `json:"tendermintNodeURL"` // tendemint noed url

	BorChainID string `json:"borChainID"` // bor chain id

	StakeManagerAddress  string `json:"stakeManagerAddress"`  // Stake manager address on main chain
	RootchainAddress     string `json:"rootchainAddress"`     // Rootchain contract address on main chain
	ValidatorSetAddress  string `json:"validatorSetAddress"`  // Validator Set contract address on bor chain
	StateSenderAddress   string `json:"stateSenderAddress"`   // main
	StateReceiverAddress string `json:"stateReceiverAddress"` // matic

	ChildBlockInterval uint64 `json:"childBlockInterval"` // Difference between header index of 2 child blocks submitted on main chain

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
var GenesisDoc tmTypes.GenesisDoc

// Contracts
// var RootChain types.Contract
// var StakeManager types.Contract
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

	if strings.Compare(conf.MaticRPCUrl, "") != 0 {
		return
	}

	configDir := filepath.Join(homeDir, "config")
	fmt.Println("Initializing tendermint configurations", "configDir", configDir)

	heimdallViper := viper.New()
	if heimdallConfigFilePath == "" {
		heimdallViper.SetConfigName("heimdall-config") // name of config file (without extension)
		heimdallViper.AddConfigPath(configDir)         // call multiple times to add many search paths
		fmt.Println("Loading heimdall configurations", "file", filepath.Join(configDir, "heimdall-config.json"))
	} else {
		heimdallViper.SetConfigFile(heimdallConfigFilePath) // set config file explicitly
		fmt.Println("Loading heimdall configurations", "file", heimdallConfigFilePath)
	}

	err := heimdallViper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}

	if err = heimdallViper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	if mainRPCClient, err = rpc.Dial(conf.MainRPCUrl); err != nil {
		log.Fatal(err)
	}

	mainChainClient = ethclient.NewClient(mainRPCClient)
	if maticRPCClient, err = rpc.Dial(conf.MaticRPCUrl); err != nil {
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
	privVal := privval.LoadFilePV(filepath.Join(configDir, "priv_validator_key.json"), filepath.Join(configDir, "priv_validator_key.json"))
	cdc.MustUnmarshalBinaryBare(privVal.Key.PrivKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)
}

// GetConfig returns cached configuration object
func GetConfig() Configuration {
	return conf
}

func GetGenesisDoc() tmTypes.GenesisDoc {
	return GenesisDoc
}

// func initContracts() error {
// 	rootChainInstance, err := rootchain.NewRootchain(GetRootChainAddress(), mainChainClient)
// 	if err != nil {
// 		return err
// 	}
// 	rootchainABI, err := abi.JSON(strings.NewReader(rootchain.RootchainABI))
// 	if err != nil {
// 		return err
// 	}
// 	RootChain = types.NewContract("rootchain", common.HexToAddress(GetConfig().RootchainAddress), rootchainABI, 0, rootChainInstance)
// }

//
// Root chain
//

// GetRootChainAddress returns RootChain contract address for selected base chain
func GetRootChainAddress() common.Address {
	return common.HexToAddress(GetConfig().RootchainAddress)
}

// GetRootChainInstance returns RootChain contract instance for selected base chain
func GetRootChainInstance() (*rootchain.Rootchain, error) {
	rootChainInstance, err := rootchain.NewRootchain(GetRootChainAddress(), mainChainClient)
	if err != nil {
		fmt.Println("Unable to create root chain instance", "error", err)
	}

	return rootChainInstance, err
}

// GetRootChainABI returns ABI for RootChain contract
func GetRootChainABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(rootchain.RootchainABI))
}

//
// Stake manager
//

// GetStakeManagerAddress returns StakeManager contract address for selected base chain
func GetStakeManagerAddress() common.Address {
	return common.HexToAddress(GetConfig().StakeManagerAddress)
}

// GetStakeManagerInstance returns StakeManager contract instance for selected base chain
func GetStakeManagerInstance() (*stakemanager.Stakemanager, error) {
	return stakemanager.NewStakemanager(GetStakeManagerAddress(), mainChainClient)
}

// GetStakeManagerABI returns ABI for StakeManager contract
func GetStakeManagerABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(stakemanager.StakemanagerABI))
}

//
// Validator set
//

// GetValidatorSetAddress returns Validator set contract address for selected base chain
func GetValidatorSetAddress() common.Address {
	return common.HexToAddress(GetConfig().ValidatorSetAddress)
}

//
// State sender
//

// GetStateSenderAddress returns state sender contract address for selected base chain
func GetStateSenderAddress() common.Address {
	return common.HexToAddress(GetConfig().StateSenderAddress)
}

//
// State sender
//

// GetStateReceiverAddress returns state receiver contract address for selected child chain
func GetStateReceiverAddress() common.Address {
	return common.HexToAddress(GetConfig().StateReceiverAddress)
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
