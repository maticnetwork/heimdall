package checkpoint

import (
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
)

func TestFetchHeaders(t *testing.T) {
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	result, err := checkpoint.GetHeaders(0, 10000)
	if err != nil {
		fmt.Println("error", err)
	} else {
		fmt.Println("rootHash generated ", hex.EncodeToString(result))
		fmt.Println("validating roothash ", checkpoint.ValidateCheckpoint(0, 10000, common.BytesToHash(result)))
	}
}
