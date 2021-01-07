package clerk_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	abcitypes "github.com/tendermint/tendermint/abci/types"
// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	"github.com/cosmos/cosmos-sdk/x/staking/types"
// 	"github.com/maticnetwork/heimdall/app"
// )

// func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
// 	app := app.Setup(false)
// 	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

// 	app.InitChain(
// 		abcitypes.RequestInitChain{
// 			AppStateBytes: []byte("{}"),
// 			ChainId:       "test-chain-id",
// 		},
// 	)

// 	acc := app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.BondedPoolName))
// 	require.NotNil(t, acc)

// 	acc = app.AccountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.NotBondedPoolName))
// 	require.NotNil(t, acc)
// }
