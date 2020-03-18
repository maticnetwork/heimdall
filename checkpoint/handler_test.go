package checkpoint

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	// "testing"
	// "bytes"
	// "encoding/hex"
	// "math/rand"
	// "os"
	// "time"
	// "github.com/maticnetwork/heimdall/helper/mocks"
	// "github.com/cosmos/cosmos-sdk/codec"
	// "github.com/cosmos/cosmos-sdk/store"
	// "github.com/cosmos/cosmos-sdk/x/params"
	ethcmn "github.com/maticnetwork/bor/common"
	// "github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	// "github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
)

func MakeTestCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	bankTypes.RegisterCodec(cdc)

	// custom types
	borTypes.RegisterCodec(cdc)
	RegisterCodec(cdc)
	staking.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}

// init for test cases
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, staking.Keeper, Keeper) {
	//t.Parallel()
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	// TODO create more keys like borKey etc
	keyCheckpoint := sdk.NewKVStoreKey("checkpoint")
	keyStaking := sdk.NewKVStoreKey("staking")
	keyMaster := sdk.NewKVStoreKey("master")
	keyParams := sdk.NewKVStoreKey(subspace.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(subspace.TStoreKey)

	// mount all
	ms.MountStoreWithDB(keyCheckpoint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMaster, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)

	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	cdc := MakeTestCodec()
	//pulp := MakeTestPulp()
	paramsKeeper := params.NewKeeper(cdc, keyParams, tKeyParams)

	dummyStakingKeeper := staking.Keeper{}

	checkpointKeeper := NewKeeper(
		cdc,
		dummyStakingKeeper,
		keyCheckpoint,
		paramsKeeper.Subspace(checkpointTypes.DefaultParamspace),
		common.DefaultCodespace,
	)

	stakingKeeper := staking.NewKeeper(
		cdc,
		keyStaking,
		paramsKeeper.Subspace(stakingTypes.DefaultParamspace),
		common.DefaultCodespace,
		checkpointKeeper,
	)
	checkpointKeeper.sk = stakingKeeper
	return ctx, stakingKeeper, checkpointKeeper
}

// create random header block
func GenRandCheckpointHeader(start int, headerSize int) (headerBlock types.CheckpointBlockHeader, err error) {
	start = start
	end := start + headerSize
	roothash, err := GetHeaders(uint64(start), uint64(end))
	if err != nil {
		return headerBlock, err
	}
	proposer := ethcmn.Address{}
	headerBlock = types.CreateBlock(uint64(start), uint64(end), types.HexToHeimdallHash(hex.EncodeToString(roothash)), types.HexToHeimdallHash(hex.EncodeToString(roothash)), types.HexToHeimdallAddress(proposer.String()), uint64(time.Now().Unix()))

	return headerBlock, nil
}

// Generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []types.Validator) {
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
			ID:               types.NewValidatorID(startID + uint64(i)),
			StartEpoch:       startBlock,
			EndEpoch:         startBlock + timeAlive,
			VotingPower:      power,
			Signer:           types.HexToHeimdallAddress(pubkey.Address().String()),
			PubKey:           pubkey,
			ProposerPriority: 0,
		}
		validators = append(validators, newVal)
	}
	return
}

// Load Validator Set
func LoadValidatorSet(count int, t *testing.T, keeper staking.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	// create 4 validators
	validators := GenRandomVal(4, 0, 10, uint64(timeAlive), randomise, 1)
	var valSet types.ValidatorSet
	// add validators to new Validator set and state
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.Empty(t, err, "Unable to set validator, Error: %v", err)
		// add validator to validator set
		// valSet.Add(&validator)
		valSet.UpdateWithChangeSet([]*types.Validator{&validator})
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")

	vals := keeper.GetAllValidators(ctx)
	t.Log("Vals added to store", len(vals))
	return valSet
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

// test handler for message
func TestHandleMsgCheckpoint(t *testing.T) {
	contractCallerObj := mocks.IContractCaller{}

	// check valid checkpoint
	t.Run("validCheckpoint", func(t *testing.T) {
		ctx, sk, ck := CreateTestInput(t, false)
		// generate proposer for validator set
		LoadValidatorSet(4, t, sk, ctx, false, 10)
		sk.IncrementAccum(ctx, 1)
		header, err := GenRandCheckpointHeader(0, 10)
		require.Empty(t, err, "Unable to create random header block, Error:%v", err)
		// make sure proposer has min ether
		contractCallerObj.On("GetBalance", sk.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)
		SentValidCheckpoint(header, ck, sk, ctx, contractCallerObj, t)
	})

	// check invalid proposer
	t.Run("invalidProposer", func(t *testing.T) {
		ctx, sk, ck := CreateTestInput(t, false)
		// generate proposer for validator set
		LoadValidatorSet(4, t, sk, ctx, false, 10)
		sk.IncrementAccum(ctx, 1)
		header, err := GenRandCheckpointHeader(0, 10)
		require.Empty(t, err, "Unable to create random header block, Error:%v", err)

		// add wrong proposer to header
		header.Proposer = sk.GetValidatorSet(ctx).Validators[2].Signer
		// make sure proposer has min ether
		contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
		// create checkpoint msg
		genesisValidatorRewards := sk.GetAllValidatorRewards(ctx)
		genesisrewardRootHash, _ := GetRewardRootHash(genesisValidatorRewards)
		header.RewardRootHash = types.BytesToHeimdallHash(genesisrewardRootHash)
		// create checkpoint msg
		msgCheckpoint := NewMsgCheckpointBlock(header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RewardRootHash,
			header.TimeStamp) // send checkpoint to handler
		got := handleMsgCheckpoint(ctx, msgCheckpoint, ck, &contractCallerObj)
		require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got.IsOK())
	})

	t.Run("multipleCheckpoint", func(t *testing.T) {
		t.Run("afterTimeout", func(t *testing.T) {
			ctx, sk, ck := CreateTestInput(t, false)
			// generate proposer for validator set
			LoadValidatorSet(4, t, sk, ctx, false, 10)
			sk.IncrementAccum(ctx, 1)
			header, err := GenRandCheckpointHeader(0, 10)
			require.Empty(t, err, "Unable to create random header block, Error:%v", err)

			// add current proposer to header
			header.Proposer = sk.GetValidatorSet(ctx).Proposer.Signer
			// make sure proposer has min ether
			contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
			// create checkpoint 257 seconds prev to current time
			header.TimeStamp = uint64(time.Now().Add(-(helper.CheckpointBufferTime + time.Second)).Unix())
			t.Log("Sending checkpoint with timestamp", "Timestamp", header.TimeStamp, "Current", time.Now().UTC().Unix())
			// send old checkpoint
			SentValidCheckpoint(header, ck, sk, ctx, contractCallerObj, t)

			header, err = GenRandCheckpointHeader(0, 10)
			header.Proposer = sk.GetValidatorSet(ctx).Proposer.Signer
			// create new checkpoint with current time
			header.TimeStamp = uint64(time.Now().Unix())
			genesisValidatorRewards := sk.GetAllValidatorRewards(ctx)
			genesisrewardRootHash, _ := GetRewardRootHash(genesisValidatorRewards)
			header.RewardRootHash = types.BytesToHeimdallHash(genesisrewardRootHash)
			msgCheckpoint := NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.RewardRootHash, header.TimeStamp)
			// send new checkpoint which should replace old one
			got := handleMsgCheckpoint(ctx, msgCheckpoint, ck, &contractCallerObj)
			require.True(t, got.IsOK(), "expected send-checkpoint to be  ok, got %v", got)
		})

		t.Run("beforeTimeout", func(t *testing.T) {
			ctx, sk, ck := CreateTestInput(t, false)
			// generate proposer for validator set
			LoadValidatorSet(4, t, sk, ctx, false, 10)
			sk.IncrementAccum(ctx, 1)
			header, err := GenRandCheckpointHeader(0, 10)
			require.Empty(t, err, "Unable to create random header block, Error:%v", err)

			// add current proposer to header
			header.Proposer = sk.GetValidatorSet(ctx).Proposer.Signer
			// make sure proposer has min ether
			contractCallerObj.On("GetBalance", header.Proposer).Return(helper.MinBalance, nil)
			// add current proposer to header
			header.Proposer = sk.GetValidatorSet(ctx).Proposer.Signer
			// send old checkpoint
			SentValidCheckpoint(header, ck, sk, ctx, contractCallerObj, t)
			genesisValidatorRewards := sk.GetAllValidatorRewards(ctx)
			genesisrewardRootHash, _ := GetRewardRootHash(genesisValidatorRewards)
			header.RewardRootHash = types.BytesToHeimdallHash(genesisrewardRootHash)
			// create checkpoint msg
			msgCheckpoint := NewMsgCheckpointBlock(header.Proposer, header.StartBlock, header.EndBlock, header.RootHash, header.RewardRootHash, uint64(time.Now().Unix()))

			// send checkpoint to handler
			got := handleMsgCheckpoint(ctx, msgCheckpoint, ck, &contractCallerObj)
			require.True(t, !got.IsOK(), "expected send-checkpoint to be not ok, got %v", got)
		})
	})

}

func SentValidCheckpoint(header types.CheckpointBlockHeader, ck Keeper, sk staking.Keeper, ctx sdk.Context, contractCallerObj mocks.IContractCaller, t *testing.T) {
	// add current proposer to header
	header.Proposer = sk.GetValidatorSet(ctx).Proposer.Signer
	genesisValidatorRewards := sk.GetAllValidatorRewards(ctx)
	genesisrewardRootHash, _ := GetRewardRootHash(genesisValidatorRewards)
	header.RewardRootHash = types.BytesToHeimdallHash(genesisrewardRootHash)
	// create checkpoint msg
	msgCheckpoint := NewMsgCheckpointBlock(header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		header.RewardRootHash,
		header.TimeStamp)
	// send checkpoint to handler
	got := handleMsgCheckpoint(ctx, msgCheckpoint, ck, &contractCallerObj)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)
	storedHeader, err := ck.GetCheckpointFromBuffer(ctx)
	t.Log("Header added to buffer", storedHeader.String())
	require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)
}
