package helper

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/heimdall/contracts/validatorset"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"log"
	"os"
)

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)

	initHeimdall()
}

type Configuration struct {
	validatorFilePVPath string
	stakeManagerAddress string
	rootchainAddress    string
	validatorSetAddress string
	dialKovan           string
	dialMatic           string
}

var configuration Configuration

func initHeimdall() {
	file, _ := os.Open("../config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func GetConfig() (configration Configuration) {
	return configration
}

func getValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(GetConfig().validatorSetAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return validatorSetInstance
}

var kovanClient, _ = ethclient.Dial(GetConfig().dialKovan)
var maticClient, _ = ethclient.Dial(GetConfig().dialMatic)
