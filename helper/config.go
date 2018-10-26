package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"log"

	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/validatorSet"
	logger "github.com/maticnetwork/heimdall/libs"
	"os"
	"strings"
)

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)

	Logger = logger.NewMainLogger(logger.NewSyncWriter(os.Stdout))
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

func InitHeimdallConfig() {
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

	if MaticClient, err = ethclient.Dial(GetConfig().MaticRPCUrl); err != nil {
		log.Fatal(err)
	}
}

func GetConfig() Configuration {
	if strings.Compare(conf.MaticRPCUrl, "") == 0 {
		InitHeimdallConfig()
	}

	return conf
}

func GetValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(GetConfig().ValidatorSetAddress), client)
	if err != nil {
		Logger.Error("Unable to create validator set instance", "Error", err, "Client", client)
	}

	return validatorSetInstance
}

func GetStakeManagerInstance(client *ethclient.Client) *stakemanager.Stakemanager {
	stakeManagerInstance, err := stakemanager.NewStakemanager(common.HexToAddress(GetConfig().ValidatorSetAddress), client)
	if err != nil {
		Logger.Error("Unable to create stakemanager instance", "Error", err, "Client", client)
	}

	return stakeManagerInstance
}
