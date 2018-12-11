package test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"math/rand"
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
	helper.InitHeimdallConfig()
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
	masterKeeper.UpdateACKCountWithValue(ctx, 1)

	return ctx, masterKeeper
}

// TODO check why initHeimdall not working here
// create random header block
func GenRandCheckpointHeader() (headerBlock types.CheckpointBlockHeader, err error) {
	//start := rand.Intn(100) + 10
	//end := start + 256
	//var headerBlock types.CheckpointBlockHeader
	//roothash, err := checkpoint.GetHeaders(uint64(start), uint64(end))
	//if err != nil {
	//	return headerBlock, err
	//}
	proposer := ethcmn.Address{}
	headerBlock = types.CreateBlock(uint64(4733040), uint64(4733050), ethcmn.HexToHash("0x5ba1680c5f5d5da8c7e3c08ba5d168c69da7a7104cf4beab94f7c0c955551f35"), proposer, rand.Uint64())
	return headerBlock, nil
}

// TODO autogenerate validator instead of
func GenRandomVal() types.Validator {
	pubkey := types.NewPubKey([]byte("0x5ba1680c5f5d5da8c7e3c08ba5d168c69da7a7104cf4beab94f7c0c955551f35"))
	return types.Validator{
		Address:    ethcmn.HexToAddress("0x660b992672675153ed263424E5dD48c2cD2DBf4f"),
		StartEpoch: 2,
		EndEpoch:   1,
		Power:      10,
		Signer:     pubkey.Address(),
		PubKey:     pubkey,
		Accum:      0,
	}
}
