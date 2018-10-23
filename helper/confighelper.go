package helper

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/contracts/validatorset"
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
}

var conf Configuration

var KovanClient *ethclient.Client
var MaticClient *ethclient.Client

func initHeimdall() {
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("config")                // call multiple times to add many search paths
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatal(err)
	}

	var err error
	// setup eth client
	if KovanClient, err = ethclient.Dial(GetConfig().DialKovan); err != nil {
		log.Fatal(err)
	}

	if MaticClient, err = ethclient.Dial(GetConfig().DialKovan); err != nil {
		log.Fatal(err)
	}
}

func GetConfig() Configuration {
	return conf
}

func GetValidatorSetInstance(client *ethclient.Client) *validatorSet.ValidatorSet {
	validatorSetInstance, err := validatorSet.NewValidatorSet(common.HexToAddress(GetConfig().ValidatorSetAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return validatorSetInstance
}
