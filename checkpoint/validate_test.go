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
	start := uint64(0)
	end:=uint64(300)
	result, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		fmt.Println("error", err)
	} else {
		fmt.Println("rootHash generated ", hex.EncodeToString(result))
		fmt.Println("validating roothash ", checkpoint.ValidateCheckpoint(start, end, common.BytesToHash(result)))
	}
}
