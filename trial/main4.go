package main

import (
	"github.com/basecoin/contracts"
)

func main() {
	//client, err := ethclient.Dial("https://kovan.infura.io")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//stakeManagerAddress := "8b28d78eb59c323867c43b4ab8d06e0f1efa1573"
	//stakeManagerInstance, err := StakeManager.NewContracts(common.HexToAddress(stakeManagerAddress), client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//last, _ := stakeManagerInstance.LastValidatorIndex(nil)
	//fmt.Println("The last validator index is %v", last)
	//validator, _ := stakeManagerInstance.Validators(nil, big.NewInt(int64(0)))
	//fmt.Println("The validator is %v", validator)
	contract.SelectProposer()
}
