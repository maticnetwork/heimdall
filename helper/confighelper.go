package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/contracts/validatorset"
	"os"

	"github.com/maticnetwork/heimdall/libs"
)

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)

	initHeimdall()

}

type Configuration struct {
	ValidatorFilePVPath string
	StakeManagerAddress string
	RootchainAddress    string
	ValidatorSetAddress string
	DialKovan           string
	DialMatic           string
	RPCUrl              string
}

var conf Configuration

var KovanClient *ethclient.Client
var MaticClient *ethclient.Client

var Logger log.Logger

func initHeimdall() {
	Logger = log.NewMainLogger(log.NewSyncWriter(os.Stdout))

	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("config")                // call multiple times to add many search paths
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		Logger.Error("Error reading Config From Viper %v", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		Logger.Error("Error unmarshalling config %v", err)
	}

	var err error
	// setup eth client
	if KovanClient, err = ethclient.Dial(GetConfig().DialKovan); err != nil {
		Logger.Error("Error creating kovanclient : %v", err)
	}

	if MaticClient, err = ethclient.Dial(GetConfig().DialKovan); err != nil {
		Logger.Error("Error creating maticclient : %v", err)
	}
}

func GetConfig() Configuration {
	return conf
}

func GetValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(GetConfig().ValidatorSetAddress), client)
	if err != nil {
		Logger.Error("Error creating Validator Set Instance : %v", err)
	}
	return validatorSetInstance
}
