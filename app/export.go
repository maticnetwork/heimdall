package app

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// ExportAppStateAndValidators exports the state of heimdall for a genesis file
func (app *HeimdallApp) ExportAppStateAndValidators() (
	appState json.RawMessage,
	validators []tmTypes.GenesisValidator,
	err error,
) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
	result := app.mm.ExportGenesis(ctx)

	// create app state
	appState, err = jsoniter.ConfigFastest.Marshal(result)

	return appState, validators, err
}
