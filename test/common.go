package test

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"encoding/hex"
	"os"
	"time"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
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

// init for test cases
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, common.Keeper) {
	//t.Parallel()
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	keyCheckpoint := sdk.NewKVStoreKey("checkpoint")
	keyStaker := sdk.NewKVStoreKey("staking")
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
	masterKeeper.UpdateACKCountWithValue(ctx, 1)

	return ctx, masterKeeper
}

// create random header block
func GenRandCheckpointHeader(headerSize int) (headerBlock types.CheckpointBlockHeader, err error) {
	start := rand.Intn(100) + 1
	end := start + headerSize
	roothash, err := checkpoint.GetHeaders(uint64(start), uint64(end))
	if err != nil {
		return headerBlock, err
	}
	proposer := ethcmn.Address{}
	headerBlock = types.CreateBlock(uint64(start), uint64(end), ethcmn.HexToHash(hex.EncodeToString(roothash)), proposer, uint64(time.Now().Unix()))

	return headerBlock, nil
}

// Generate random validators
func GenRandomVal(count int, startBlock uint64, power uint64, timeAlive uint64, randomise bool, startID uint64) (validators []types.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := types.NewPubKey(privKey1.PubKey().Bytes())
		if randomise {
			startBlock := uint64(rand.Intn(10))
			// todo find a way to genrate non zero random number
			if startBlock == 0 {
				startBlock = 1
			}
			power := uint64(rand.Intn(100))
			if power == 0 {
				power = 1
			}
		}

		newVal := types.Validator{
			ID:         types.NewValidatorID(startID + uint64(i)),
			StartEpoch: startBlock,
			EndEpoch:   startBlock + timeAlive,
			Power:      power,
			Signer:     pubkey.Address(),
			PubKey:     pubkey,
			Accum:      0,
		}
		validators = append(validators, newVal)
	}
	return
}

// finds address in give validator slice
func FindSigner(signer ethcmn.Address, vals []types.Validator) bool {
	for _, val := range vals {
		if bytes.Compare(signer.Bytes(), val.Signer.Bytes()) == 0 {
			return true
		}
	}
	return false
}

// print validators
func PrintVals(t *testing.T, valset *types.ValidatorSet) {
	t.Log("Printing validators")
	for i, val := range valset.Validators {
		t.Log("Validator ===> ", "Index", i, "ValidatorInfo", val.String())
	}
}
