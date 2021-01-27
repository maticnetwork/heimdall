package gov_test

import (
	"bytes"
	"log"
	"sort"
	"encoding/json"

	// "github.com/tendermint/tendermint/crypto"
	// "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
)

// TODO - Merge this with keeper/common_test.go
//
// Create test app
//

// returns context and app with params set on gov keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	govGenesis := types.NewGenesisState(types.DefaultGenesis().StartingProposalId, types.DefaultGenesis().DepositParams, types.DefaultGenesis().VotingParams, types.DefaultGenesis().TallyParams)

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(app.AppCodec())

	genesisState[types.ModuleName] = app.AppCodec().MustMarshalJSON(&govGenesis)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})

	return app, ctx, cliCtx
}

var (
	valTokens           = sdk.TokensFromConsensusPower(42)
	TestProposal        = types.NewTextProposal("Test", "description")
	// TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
	// TestCommissionRates = stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
)

// SortAddresses - Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	byteAddrs := make([][]byte, len(addrs))

	for i, addr := range addrs {
		byteAddrs[i] = addr.Bytes()
	}

	SortByteArrays(byteAddrs)

	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// SortByteArrays - sorts the provided byte array
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

const contextKeyBadProposal = "contextKeyBadProposal"

// var (
// 	pubkeys = []crypto.PubKey{
// 		ed25519.GenPrivKey().PubKey(),
// 		ed25519.GenPrivKey().PubKey(),
// 		ed25519.GenPrivKey().PubKey(),
// 	}
// )

// func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) {
// 	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")

// 	for i := 0; i < len(addrs); i++ {
// 		valTokens := sdk.TokensFromConsensusPower(powerAmt[i])
// 		valCreateMsg, err := stakingtypes.NewMsgCreateValidator(
// 			addrs[i], pubkeys[i], sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
// 			TestDescription, TestCommissionRates, sdk.OneInt(),
// 		)
// 		require.NoError(t, err)
// 		handleAndCheck(t, stakingHandler, ctx, valCreateMsg)
// 	}
// }

// func handleAndCheck(t *testing.T, h sdk.Handler, ctx sdk.Context, msg sdk.Msg) {
// 	res, err := h(ctx, msg)
// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// }

