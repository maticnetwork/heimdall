package checkpoint

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

type Keeper struct {
	checkpointKey sdk.StoreKey
	cdc           *wire.Codec
	//validatorSet sdk.ValidatorSet

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
	rootHash   common.Hash
}

func createBlock(start uint64, end uint64, rootHash common.Hash) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		rootHash:   rootHash,
	}
}
func (k Keeper) addCheckpoint(ctx sdk.Context, start uint64, end uint64, root common.Hash) int64 {
	store := ctx.KVStore(k.checkpointKey)
	data := createBlock(start, end, root)
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	//TODO add block data validation
	fmt.Printf("Block data to be inserted with key %v", []byte(strconv.Itoa(int(ctx.BlockHeight()))))
	fmt.Printf("Block data to be inserted is %v", out)
	store.Set([]byte(strconv.Itoa(int(ctx.BlockHeight()))), []byte(out))
	return ctx.BlockHeight()
}
func (k Keeper) getCheckpoint(ctx sdk.Context, key int64) []byte {
	store := ctx.KVStore(k.checkpointKey)
	getKey := []byte(strconv.Itoa(int(key)))
	return store.Get(getKey)

}
