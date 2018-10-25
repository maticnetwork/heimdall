package checkpoint

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
	conf "github.com/maticnetwork/heimdall/helper"
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
	var logger = conf.Logger.With("module", "checkpoint")

	store := ctx.KVStore(k.checkpointKey)

	data := createBlock(start, end, root, proposer)
	out, err := json.Marshal(data)
	if err != nil {
		logger.Error("Error Marshalling checkpoint to json ", "Error", err)
	}

	//TODO add block data validation
	logger.Info("Block data inserted ! ", out, "BlockData ")

	store.Set([]byte(strconv.Itoa(int(ctx.BlockHeight()))), []byte(out))

	return ctx.BlockHeight()
}

func (k Keeper) GetCheckpoint(ctx sdk.Context, key int64) []byte {

	store := ctx.KVStore(k.checkpointKey)
	getKey := []byte(strconv.Itoa(int(key)))

	return store.Get(getKey)
}
