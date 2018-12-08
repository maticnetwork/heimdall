package test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func MakeTestCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	// custom types
	checkpoint.RegisterWire(cdc)
	staking.RegisterWire(cdc)

	cdc.Seal()
	return cdc
}

func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, common.Keeper) {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	keyCheckpoint := sdk.NewKVStoreKey("checkpoint")
	keyStaker := sdk.NewKVStoreKey("staker")
	keyMaster := sdk.NewKVStoreKey("master")
	ms.MountStoreWithDB(keyCheckpoint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaker, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMaster, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	cdc := MakeTestCodec()
	//pulp := MakeTestPulp()

	masterKeeper := common.NewKeeper(cdc, keyMaster, keyStaker, keyCheckpoint, common.DefaultCodespace)
	// set empty values in cache by default
	masterKeeper.SetCheckpointAckCache(ctx, common.EmptyBufferValue)
	masterKeeper.SetCheckpointCache(ctx, common.EmptyBufferValue)
	masterKeeper.UpdateACKCountWithValue(ctx, 0)

	return ctx, masterKeeper
}

func GenRandCheckpointHeader() {

}
