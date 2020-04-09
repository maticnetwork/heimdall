package app

import (
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func TestHeimdalldExport(t *testing.T) {
	db := db.NewMemDB()
	happ := NewHeimdallApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	genesisState := NewDefaultGenesisState()
	err := setGenesis(happ, genesisState)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newHapp := NewHeimdallApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	_, _, err = newHapp.ExportAppStateAndValidators()
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestMakePulp(t *testing.T) {
	pulp := MakePulp()
	require.NotNil(t, pulp, "Pulp should be nil")
}

func setGenesis(app *HeimdallApp, genesisState GenesisState) error {
	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app.Commit()
	return nil
}
