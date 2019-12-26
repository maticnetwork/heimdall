package app

import (
	"encoding/json"

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
	// appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	appState, err = json.Marshal(result)
	return appState, validators, err
}
