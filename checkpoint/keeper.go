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

var (
	ACKCountKey         = []byte{0x01}
	LastCheckpointKey   = []byte{0x02}
	BufferCheckpointKey = []byte{0x03}
	HeaderBlockKey      = []byte{0x04}
)

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
func (k Keeper) AddCheckpointToKey(ctx sdk.Context, start uint64, end uint64, root common.Hash, proposer common.Address, key []byte) {
	store := ctx.KVStore(k.checkpointKey)

	// create Checkpoint block and marshall
	data := createBlock(start, end, root, proposer)
	out, err := json.Marshal(data)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}

	// store in key provided
	store.Set(key, []byte(out))

}

func (k Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)
	store.Set(LastCheckpointKey, []byte(""))
}

func (k Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.checkpointKey)

	// Get checkpoint and unmarshall
	var checkpoint CheckpointBlockHeader
	err := json.Unmarshal(store.Get(BufferCheckpointKey), &checkpoint)

	return checkpoint, err
}

// sets last header block
// we can eliminate this by having CHILD_BLOCK_INTERVAL*(TotalACKS+1)
func (k Keeper) SetLastCheckpointKey(ctx sdk.Context, _key int64) {
	store := ctx.KVStore(k.checkpointKey)

	// set last checkpoint key
	key := []byte(strconv.Itoa(int(_key)))
	store.Set(LastCheckpointKey, key)
}

func (k Keeper) GetLastCheckpointKey(ctx sdk.Context) (key int64) {
	store := ctx.KVStore(k.checkpointKey)

	// get last checkpoint
	keyInt, err := strconv.Atoi(string(store.Get(LastCheckpointKey)))
	if err != nil {
		CheckpointLogger.Error("Unable to convert key to int")
	}
	return int64(keyInt)
}

// update ACK count by 1
func (k Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.Itoa(ACKCount + 1))

	// update
	store.Set(ACKCountKey, ACKs)
}

// Get current ACK count
func (k Keeper) GetACKCount(ctx sdk.Context) int {
	store := ctx.KVStore(k.checkpointKey)

	// get current ACK count
	ACKs, err := strconv.Atoi(string(store.Get(ACKCountKey)))
	if err != nil {
		CheckpointLogger.Error("Unable to convert key to int")
	}

	return ACKs
}

// Set ACK Count to 0
func (k Keeper) InitACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)

	// set to 0
	key := []byte(strconv.Itoa(int(0)))
	store.Set(ACKCountKey, key)
}

func GetHeaderKey(headerNumber int) []byte {
	headerNumberBytes := strconv.Itoa(headerNumber)
	return append(HeaderBlockKey, headerNumberBytes...)
}
