package staking_test


//TODO uncomment and fix test case
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
