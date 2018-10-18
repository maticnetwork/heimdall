package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)
}
func main() {

	//var x secp256k1.PubKeySecp256k1
	//k, _ := hex.DecodeString("041FE1CDE7D9D8C9182AC967EC8362262216FF8A10061F0DE0F1472F9E45F965D0909DE527E18C7BFB9FCD42335E60FB6E18367A4DC37F1A7FC3265C7241597973")
	//copy(x[:], k[:])
	//fmt.Printf("slice 1 %v  and 2 %v  ", x[:32], x[32:])
	//fmt.Println("%v", x)
	//fmt.Println("Address is :%v", x.Address())
	privVal := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")

	// retrieve private key
	var pkObject secp256k1.PrivKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privVal.PrivKey.Bytes(), &pkObject)

	// create ecdsa private key

	// from address
	fromAddress := common.BytesToAddress(privVal.Address)
	fmt.Println("public key %v and from address %v", privVal.PubKey, fromAddress.String())

	fmt.Printf(" part 1 %v , part 2 %v", hex.EncodeToString(privVal.PubKey.Bytes()[:32]), hex.EncodeToString(privVal.PubKey.Bytes()[32:]))
}
