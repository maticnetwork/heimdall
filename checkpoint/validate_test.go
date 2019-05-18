package checkpoint

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/stretchr/testify/require"
)

func TestFetchHeaders(t *testing.T) {
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	start := uint64(0)
	end:=uint64(300)
	result, err := checkpoint.GetHeaders(start, end)
	require.Empty(t, err, "Unable to fetch headers, Error:%v", err)
	ok:= checkpoint.ValidateCheckpoint(start, end, common.BytesToHash(result))
	require.Equal(t,true,ok,"Root hash should match ")
}
