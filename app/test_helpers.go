package app

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/app/helpers"
	authTypes "github.com/maticnetwork/heimdall/auth/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// Setup initializes a new App. A Nop logger is set in App.
func Setup(isCheckTx bool) *HeimdallApp {
	db := dbm.NewMemDB()
	app := NewHeimdallApp(log.NewNopLogger(), db)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:    []abci.ValidatorUpdate{},
				AppStateBytes: stateBytes,
			},
		)
	}

	return app
}

// SetupWithGenesisAccounts initializes a new Heimdall with the provided genesis
// accounts and possible balances.
func SetupWithGenesisAccounts(genAccs []authTypes.GenesisAccount) *HeimdallApp {
	// setup with isCheckTx
	app := Setup(true)

	// initialize the chain with the passed in genesis accounts
	genesisState := NewDefaultGenesisState()

	authGenesis := authTypes.NewGenesisState(authTypes.DefaultParams(), genAccs)
	genesisState[authTypes.ModuleName] = app.Codec().MustMarshalJSON(authGenesis)

	// bankGenesis := authTypes.NewGenesisState(authTypes.DefaultGenesisState().SendEnabled)
	// genesisState[authTypes.ModuleName] = app.Codec().MustMarshalJSON(bankGenesis)

	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
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
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})

	return app
}

// GenerateAccountStrategy account strategy
type GenerateAccountStrategy func(int) []hmTypes.HeimdallAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []hmTypes.HeimdallAddress {
	testAddrs := make([]hmTypes.HeimdallAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := secp256k1.GenPrivKey().PubKey()
		testAddrs[i] = hmTypes.BytesToHeimdallAddress(pk.Address().Bytes())
	}

	return testAddrs
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []hmTypes.HeimdallAddress {
	var addresses []hmTypes.HeimdallAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		addresses = append(addresses, hmTypes.HexToHeimdallAddress(buffer.String()))
		buffer.Reset()
	}

	return addresses
}

// AddTestAddrsFromPubKeys adds the addresses into the SimApp providing only the public keys.
func AddTestAddrsFromPubKeys(app *HeimdallApp, ctx sdk.Context, pubKeys []crypto.PubKey, accAmt sdk.Int) {
	initCoins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, accAmt))

	setTotalSupply(app, ctx, accAmt, len(pubKeys))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, pubKey := range pubKeys {
		saveAccount(app, ctx, hmTypes.BytesToHeimdallAddress(pubKey.Address().Bytes()), initCoins)
	}
}

// setTotalSupply provides the total supply based on accAmt * totalAccounts.
func setTotalSupply(app *HeimdallApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	// totalSupply := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, accAmt.MulRaw(int64(totalAccounts))))
	// prevSupply := app.SupplyKeeper.GetSupply(ctx)
	// app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *HeimdallApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []hmTypes.HeimdallAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createRandomAccounts)
}

// AddTestAddrsIncremental constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *HeimdallApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []hmTypes.HeimdallAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

func addTestAddrs(app *HeimdallApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []hmTypes.HeimdallAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, accAmt))
	setTotalSupply(app, ctx, accAmt, accNum)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		saveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// saveAccount saves the provided account into the simapp with balance based on initCoins.
func saveAccount(app *HeimdallApp, ctx sdk.Context, addr hmTypes.HeimdallAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	_, err := app.BankKeeper.AddCoins(ctx, addr, initCoins)
	if err != nil {
		panic(err)
	}
}

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *HeimdallApp, addr hmTypes.HeimdallAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, abci.Header{})
	account := app.AccountKeeper.GetAccount(ctxCheck, addr)
	require.True(t, balances.IsEqual(account.GetCoins()))
}

// SignCheckDeliver checks a generated signed transaction and simulates a
// block commitment with the given transaction. A test assertion is made using
// the parameter 'expPass' against the result. A corresponding result is
// returned.
func SignCheckDeliver(
	t *testing.T, cdc *codec.Codec, app *bam.BaseApp, header abci.Header, msgs []sdk.Msg,
	accNums, seq []uint64, expSimPass, expPass bool, priv ...crypto.PrivKey,
) (sdk.Result, error) {
	// generate tx
	tx := helpers.GenTx(
		msgs,
		sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
		helpers.DefaultGenTxGas,
		"",
		accNums,
		seq,
		priv...,
	)

	txBytes, err := cdc.MarshalBinaryBare(tx)
	require.Nil(t, err)

	// Must simulate now as CheckTx doesn't run Msgs anymore
	res := app.Simulate(txBytes, tx)

	if expSimPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	// Simulate a sending a transaction and committing a block
	app.BeginBlock(abci.RequestBeginBlock{Header: header})
	res = app.Deliver(tx)

	if expPass {
		require.NoError(t, err)
		require.NotNil(t, res)
	} else {
		require.Error(t, err)
		require.Nil(t, res)
	}

	app.EndBlock(abci.RequestEndBlock{})
	app.Commit()

	return res, err
}

// GenSequenceOfTxs generates a set of signed transactions of messages, such
// that they differ only by having the sequence numbers incremented between
// every transaction.
func GenSequenceOfTxs(msgs []sdk.Msg, accNums []uint64, initSeqNums []uint64, numToGenerate int, priv ...crypto.PrivKey) []authTypes.StdTx {
	txs := make([]authTypes.StdTx, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		txs[i] = helpers.GenTx(
			msgs,
			sdk.Coins{sdk.NewInt64Coin(sdk.DefaultBondDenom, 0)},
			helpers.DefaultGenTxGas,
			"",
			accNums,
			initSeqNums,
			priv...,
		)
		incrementAllSequenceNumbers(initSeqNums)
	}

	return txs
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}

// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
func CreateTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") // base pubkey string
		buffer.WriteString(numString)                                                       // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}

	return publicKeys
}

// NewPubKeyFromHex returns a PubKey from a hex string.
func NewPubKeyFromHex(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd secp256k1.PubKeySecp256k1
	copy(pkEd[:], pkBytes)
	return pkEd
}
