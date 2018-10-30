package checkpoint

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
)

type Keeper struct {
	checkpointKey sdk.StoreKey
	cdc           *wire.Codec

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		checkpointKey: key,
		cdc:           cdc,
		codespace:     codespace,
	}
	return keeper
}

type CheckpointBlockHeader struct {
	StartBlock uint64
	EndBlock   uint64
	RootHash   common.Hash
}

func createBlock(start uint64, end uint64, rootHash common.Hash) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
	}
}

func (k Keeper) AddCheckpoint(ctx sdk.Context, start uint64, end uint64, root common.Hash) int64 {
	store := ctx.KVStore(k.checkpointKey)
	data := createBlock(start, end, root)
	out, err := json.Marshal(data)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}
	store.Set([]byte(strconv.Itoa(int(ctx.BlockHeight()))), []byte(out))

	// TODO add block data validation
	CheckpointLogger.Debug("Checkpoint block saved!", "roothash", data.RootHash, "startBlock", data.StartBlock, "endBlock", data.EndBlock)

	// return new block
	return ctx.BlockHeight()
}

func (k Keeper) GetCheckpoint(ctx sdk.Context, key int64) []byte {
	store := ctx.KVStore(k.checkpointKey)
	getKey := []byte(strconv.Itoa(int(key)))

	return store.Get(getKey)
}
