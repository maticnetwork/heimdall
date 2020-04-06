package checkpoint

import (
	"os"
	"testing"

	"github.com/maticnetwork/bor/common"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
)

func TestFetchHeaders(t *testing.T) {
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	start := uint64(0)
	end := uint64(300)
	result, err := checkpointTypes.GetHeaders(start, end)
	require.Empty(t, err, "Unable to fetch headers, Error:%v", err)
	ok, err := checkpointTypes.ValidateCheckpoint(start, end, types.HeimdallHash(common.BytesToHash(result)))
	require.Equal(t, true, ok, "Root hash should match ")
}
