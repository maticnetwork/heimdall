package main

import (
	"encoding/hex"
	"fmt"
	"github.com/basecoin/checkpoint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

type obj struct {
	First  []byte
	Second []byte
}

func main() {

	msg := checkpoint.NewMsgCheckpointBlock(uint64(198), uint64(877), common.BytesToHash([]byte("lol")), "0x84f8a67E4d16bb05aBCa3d154091566921e0B5e9")

	tx := checkpoint.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		fmt.Printf("Error generating TXBYtes %v", err)
	}
	fmt.Printf("The tx bytes are %v ", hex.EncodeToString(txBytes))

	//var new obj
	//new.First = []byte("hello")
	//new.Second = []byte("world")
	//newbytes, err := rlp.EncodeToBytes(new)
	//if err != nil {
	//	fmt.Printf("Error generating TXBYtes %v", err)
	//}
	//fmt.Printf("encoded obj is %v", hex.EncodeToString(newbytes))

	// RLP decodes the txBytes to a BaseTx
	//txDecoder(txBytes)
	//res, _ := hex.DecodeString("F840F83E94FA9BF0CBA703174B2717CFEA0359F7E5E1519837834838678348386AA0D494377D4439A844214B565E1C211EA7154CA300B98E3C296F19FC9ADA36DB33")
	//hash := tmhash.Sum(res)
	//fmt.Printf("after encryp %v", hex.EncodeToString(hash))

	lol := tmhash.Sum([]byte("f840f83e94fa9bf0cba703174b2717cfea0359f7e5e1519837837e27cb837e27cea00b6e0e1df9a3c7d50b4d967272abd586b8a07dead457b5489a32023166443da9"))
	fmt.Printf("After hash result %v   ", hex.EncodeToString(lol))

}
func txDecoder(txBytes []byte) {
	var tx = checkpoint.BaseTx{}

	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		fmt.Printf("Error decoding %v", err)
	}
	fmt.Printf("After decoding data found is %v", tx.Msg.Proposer.String())
}
