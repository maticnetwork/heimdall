package helper

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/heimdall/contracts/validatorset"
	"log"
	"os"
)

type Configuration struct {
	validatorFilePVPath string
	stakeManagerAddress string
	rootchainAddress    string
	validatorSetAddress string
}

func GetConfig() (configration Configuration) {
	file, _ := os.Open("../config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configration
}
func getValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(GetConfig().validatorSetAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return validatorSetInstance
}
