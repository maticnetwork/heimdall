package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func main()  {

	var x secp256k1.PubKeySecp256k1
	k, _ := hex.DecodeString("041FE1CDE7D9D8C9182AC967EC8362262216FF8A10061F0DE0F1472F9E45F965D0909DE527E18C7BFB9FCD42335E60FB6E18367A4DC37F1A7FC3265C7241597973")
	copy(x[:], k[:])

	fmt.Println("%v", x)
	fmt.Println("Address is :%v",x.Address())

}