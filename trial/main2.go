package main
import (
"encoding/hex"
"fmt"
"github.com/ethereum/go-ethereum/crypto/sha3"
)

func main() {
hash := sha3.NewKeccak256()

var buf []byte
//hash.Write([]byte{0xcc})
fmt.Printf("cc after decode %v",decodeHex("cc"))
hash.Write(decodeHex("cc"))
buf = hash.Sum(buf)

fmt.Println(hex.EncodeToString(buf))
//expected := "EEAD6DBFC7340A56CAEDC044696A168870549A6A7F6F56961E84A54BD9970B8A"
}

func decodeHex(s string) []byte {
b, err := hex.DecodeString(s)
if err != nil {
panic(err)
}

return b
}