package app

import (
	"math/rand"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/simulation"
	simTypes "github.com/maticnetwork/heimdall/types/simulation"
)

func TestHeimdallAppExport(t *testing.T) {
	db := db.NewMemDB()
	happ := NewHeimdallApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	genesisState := NewDefaultGenesisState()

	// Get state bytes
	stateBytes, err := codec.MarshalJSONIndent(happ.Codec(), genesisState)
	require.NoError(t, err)

	// Initialize the chain
	happ.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	// Set commit
	happ.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	newHapp := NewHeimdallApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	_, _, err = newHapp.ExportAppStateAndValidators()
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestHeimdallAppExportWithRand(t *testing.T) {
	config, db, dir, logger, _, err := SetupSimulation("goleveldb-app-sim", "Simulation")
	require.NoError(t, err)
	require.NotEmpty(t, dir)
	defer func() {
		db.Close()
		require.NoError(t, os.RemoveAll(dir))
	}()
	// create seed
	config.Seed = int64(rand.Int())
	seed := rand.New(rand.NewSource(config.Seed))

	// create app
	app := NewHeimdallApp(logger, db)

	params := simulation.RandomParams(seed)
	accs := simTypes.RandomAccounts(seed, params.NumKeys())
	genesisTimestamp := simTypes.RandTimestamp(seed)

	sm := app.SimulationManager()
	appParams := make(simTypes.AppParams)
	genesisState, _ := AppStateRandomizedFn(sm, seed, app.Codec(), accs, genesisTimestamp, appParams)
	require.NotEmpty(t, string(genesisState))

	// Get state bytes
	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	require.NoError(t, err)
	require.NotEmpty(t, string(stateBytes))

	// Initialize the chain
	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	// Set commit
	app.Commit()

	// Making a new app object with the db, so that initchain hasn't been called
	newHapp := NewHeimdallApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db)
	exportedState, _, err := newHapp.ExportAppStateAndValidators()
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
	require.NotEmpty(t, string(exportedState))
}

func TestSetup(t *testing.T) {
	happ := Setup(false)
	ctx := happ.BaseApp.NewContext(false, abci.Header{})

	require.NotNil(t, happ.GetModuleManager())
	require.LessOrEqual(t, 0, len(happ.AccountKeeper.GetAllAccounts(ctx)))
}

func TestSetupWithGenesisAccounts(t *testing.T) {
	r := rand.New(rand.NewSource(42))          // seed = 42
	accounts := simTypes.RandomAccounts(r, 10) // create 10 accounts

	// genesis accounts
	var genesisAccs authTypes.GenesisAccounts
	for _, acc := range accounts {
		bacc := authTypes.NewBaseAccountWithAddress(acc.Address)
		gacc, _ := authTypes.NewGenesisAccountI(&bacc)
		genesisAccs = append(genesisAccs, gacc)
	}

	happ := SetupWithGenesisAccounts(genesisAccs)
	ctx := happ.BaseApp.NewContext(false, abci.Header{})
	require.NotNil(t, happ.GetModuleManager())
	require.LessOrEqual(t, 10, len(happ.AccountKeeper.GetAllAccounts(ctx)))
}

func TestValidateGenesis(t *testing.T) {
	happ := Setup(false)

	// not validate app state
	require.Panics(t, func() {
		happ.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: []byte("{}"),
			},
		)
	})
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}
