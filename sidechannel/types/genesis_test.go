package types_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/maticnetwork/heimdall/sidechannel/simulation"
	"github.com/maticnetwork/heimdall/sidechannel/types"
)

func TestDefaultGenesisState(t *testing.T) {
	genesis := types.DefaultGenesisState()
	require.NotNil(t, genesis, "DefaultGenesisState should not return nil response")
	require.Equal(t, 0, len(genesis.PastCommits), "DefaultGenesisState should have no pre-commits")
}

func TestNewGenesisState(t *testing.T) {
	genesis := types.NewGenesisState([]types.PastCommit{{Height: 2}})
	require.NotNil(t, genesis, "NewGenesisState should not return nil response")
	require.Equal(t, 1, len(genesis.PastCommits), "NewGenesisState should create proper pastcommits")

	pastCommit := genesis.PastCommits[0]
	require.Equal(t, int64(2), pastCommit.Height, "NewGenesisState should create commits with valid data")
}

func TestValidateGenesis(t *testing.T) {
	emptyGenesis := types.GenesisState{}
	require.Nil(t, types.ValidateGenesis(emptyGenesis), "Empty genesis should be valid genesis")

	emptyGenesis = types.NewGenesisState(make([]types.PastCommit, 0))
	require.Nil(t, types.ValidateGenesis(emptyGenesis), "Empty genesis should be valid genesis (using NewGenesisState)")

	genesis := types.NewGenesisState([]types.PastCommit{{Height: 2}})
	err := types.ValidateGenesis(genesis)
	require.Error(t, err, "PastCommit object with height 2 should not be allowed")

	// get random seed from time as source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	genesis = types.NewGenesisState(simulation.RandomPastCommits(r, 10, 0, 0))
	err = types.ValidateGenesis(genesis)
	require.Error(t, err, "PastCommit object without should not be allowed")

	genesis = types.NewGenesisState(simulation.RandomPastCommits(r, 10, 5, 0))
	err = types.ValidateGenesis(genesis)
	require.Equal(t, 10, len(genesis.PastCommits))
	require.Equal(t, 5, len(genesis.PastCommits[0].Txs))
	require.Nil(t, err, "PastCommit object with txs should not throw an error")
}
