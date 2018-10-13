package main

import (
	"encoding/hex"
	"fmt"
	"github.com/basecoin/checkpoint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {

	msg := checkpoint.NewMsgCheckpointBlock(uint64(198), uint64(877), common.BytesToHash([]byte("lol")))

	tx := checkpoint.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		fmt.Printf("Error generating TXBYtes %v", err)
	}
	fmt.Printf("The tx bytes are %v ", hex.EncodeToString(txBytes))

	// RLP decodes the txBytes to a BaseTx
	txDecoder(txBytes)

}
func txDecoder(txBytes []byte) {
	var tx = checkpoint.BaseTx{}

	err := rlp.DecodeBytes(txBytes, &tx)
	if err != nil {
		fmt.Printf("Error decoding %v", err)
	}
	fmt.Printf("After decoding data found is %v", tx.Msg.Proposer.String())
}
