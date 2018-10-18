package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"reflect"
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
	//contract.SelectProposer()
	pubkey := "04F551822EFD57CA7CCB3657EF22C54FB0349C1AF135D4AC56967D180555C6405139B60C6FCCE2081EBABD70557EB40574BA2866DDFE8284A61DC180FAAB5BCE2F"
	var pubkeyBytes secp256k1.PubKeySecp256k1
	lol, _ := hex.DecodeString(pubkey)
	copy(pubkeyBytes[:], lol)

	fmt.Printf("address obtained is %v", pubkeyBytes.Address().String())
	fmt.Printf("type is %v", reflect.TypeOf(pubkeyBytes))
}
