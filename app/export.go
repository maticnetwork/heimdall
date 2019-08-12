package app

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	auth "github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	bank "github.com/maticnetwork/heimdall/bank"
	bor "github.com/maticnetwork/heimdall/bor"
	checkpoint "github.com/maticnetwork/heimdall/checkpoint"
	staking "github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/supply"
)

// ExportAppStateAndValidators exports the state of heimdall for a genesis file
func (app *HeimdallApp) ExportAppStateAndValidators() (
	appState json.RawMessage,
	validators []tmTypes.GenesisValidator,
	err error,
) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	// iterate to get the accounts
	accounts := []GenesisAccount{}
	appendAccount := func(acc authTypes.Account) (stop bool) {
		account := NewGenesisAccount(acc)
		accounts = append(accounts, account)
		return false
	}
	app.accountKeeper.IterateAccounts(ctx, appendAccount)

	// create new genesis state
	genState := NewGenesisState(
		accounts,

		auth.ExportGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper),
		bank.ExportGenesis(ctx, app.bankKeeper),
		supply.ExportGenesis(ctx, app.supplyKeeper),

		bor.ExportGenesis(ctx, app.borKeeper),
		checkpoint.ExportGenesis(ctx, app.checkpointKeeper),
		staking.ExportGenesis(ctx, app.stakingKeeper),
	)

	// create app state
	// appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	appState, err = json.Marshal(genState)
	return appState, validators, err
}
