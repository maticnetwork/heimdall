package checkpoint

import (
	"encoding/json"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type Keeper struct {
	checkpointKey sdk.StoreKey
	cdc           *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		checkpointKey: key,
		cdc:           cdc,
		codespace:     codespace,
	}
	return keeper
}

type CheckpointBlockHeader struct {
	Proposer   common.Address
	StartBlock uint64
	EndBlock   uint64
	RootHash   common.Hash
}

func createBlock(start uint64, end uint64, rootHash common.Hash, proposer common.Address) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
	}
}

func (k Keeper) AddCheckpoint(ctx sdk.Context, start uint64, end uint64, root common.Hash, proposer common.Address) int64 {
	store := ctx.KVStore(k.checkpointKey)
	data := createBlock(start, end, root, proposer)
	out, err := json.Marshal(data)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}

	Key := []byte(strconv.Itoa(int(ctx.BlockHeight())))
	store.Set([]byte("LastCheckpointKey"), Key)
	store.Set(Key, []byte(out))

	// TODO add block data validation
	CheckpointLogger.Debug("Checkpoint block saved!", "roothash", data.RootHash, "startBlock", data.StartBlock, "endBlock", data.EndBlock)

	// return new block
	return ctx.BlockHeight()
}

func (k Keeper) GetCheckpoint(ctx sdk.Context, key int64) (CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.checkpointKey)
	Key := []byte(strconv.Itoa(int(key)))

	var checkpoint CheckpointBlockHeader
	err := json.Unmarshal(store.Get(Key), &checkpoint)
	return checkpoint, err
}

func (k Keeper) SetLastCheckpointKey(ctx sdk.Context, _key int64) {
	store := ctx.KVStore(k.checkpointKey)

	key := []byte(strconv.Itoa(int(_key)))
	store.Set([]byte("LastCheckpointKey"), key)
}

func (k Keeper) GetLastCheckpointKey(ctx sdk.Context) (key int64) {
	store := ctx.KVStore(k.checkpointKey)
	keyInt, err := strconv.Atoi(string(store.Get([]byte("LastCheckpointKey"))))
	if err != nil {
		CheckpointLogger.Error("Unable to convert key to int")
	}
	return int64(keyInt)
}

// count ACKS
func (k Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)
	ACKCount := k.GetACKCount(ctx)
	ACKs := []byte(strconv.Itoa(ACKCount + 1))
	store.Set([]byte("TotalACK"), ACKs)

}
func (k Keeper) GetACKCount(ctx sdk.Context) int {
	store := ctx.KVStore(k.checkpointKey)
	ACKs, err := strconv.Atoi(string(store.Get([]byte("TotalACK"))))
	if err != nil {
		CheckpointLogger.Error("Unable to convert key to int")
	}
	return ACKs
}

func (k Keeper) InitACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)
	key := []byte(strconv.Itoa(int(0)))
	store.Set([]byte("TotalACK"), key)

}
