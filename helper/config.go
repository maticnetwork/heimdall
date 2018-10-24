package helper

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/contracts/validatorset"
	logger "github.com/maticnetwork/heimdall/libs"
)

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)

	// initialize heimdall
	initHeimdall()
}

// Configuration represents heimdall config
type Configuration struct {
	MainRPCUrl          string `mapstructure:"main_rpcurl"`
	MaticRPCUrl         string `mapstructure:"matic_rpcurl"`
	StakeManagerAddress string `mapstructure:"stakemanager_address"`
	RootchainAddress    string `mapstructure:"rootchain_address"`
	ValidatorSetAddress string `mapstructure:"validatorset_address"`
	ValidatorFilePVPath string `mapstructure:"priv_validator_path"`

	// Tendermint endpoint
	TendermintEndpoint string `mapstructure:"tendermint_endpoint"`
}

var conf Configuration

// MainChainClient stores eth client for Main chain Network
var MainChainClient *ethclient.Client

// MaticClient stores eth client for Matic Network
var MaticClient *ethclient.Client

// Logger stores global logger object
var Logger logger.Logger

func initHeimdall() {
	heimdallViper := viper.New()
	heimdallViper.SetConfigName("heimdall-config")         // name of config file (without extension)
	heimdallViper.AddConfigPath("$HOME/.heimdalld/config") // call multiple times to add many search paths
	heimdallViper.AddConfigPath("$HOME/.heimdalld")        // call multiple times to add many search paths

	err := heimdallViper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}
	if err = heimdallViper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	// setup eth client
	if MainChainClient, err = ethclient.Dial(GetConfig().MainRPCUrl); err != nil {
		log.Fatal(err)
	}

	if MaticClient, err = ethclient.Dial(GetConfig().MainRPCUrl); err != nil {
		log.Fatal(err)
	}
	Logger = logger.NewMainLogger(logger.NewSyncWriter(os.Stdout))
}

func GetConfig() Configuration {
	return conf
}

func GetValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(GetConfig().ValidatorSetAddress), client)
	if err != nil {
		log.Fatal(err)
	}

	return validatorSetInstance
}
