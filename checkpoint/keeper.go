package checkpoint

import (
	"encoding/json"
	"strconv"

	"bytes"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"time"
)

type Keeper struct {
	checkpointKey sdk.StoreKey
	cdc           *codec.Codec
	stakingkeeper staking.Keeper

	// codespace
	codespace sdk.CodespaceType
}

var (
	ACKCountKey         = []byte{0x01}
	BufferCheckpointKey = []byte{0x02}
	HeaderBlockKey      = []byte{0x03}
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
	TimeStamp  time.Time
}

func createBlock(start uint64, end uint64, rootHash common.Hash, proposer common.Address) CheckpointBlockHeader {
	return CheckpointBlockHeader{
		StartBlock: start,
		EndBlock:   end,
		RootHash:   rootHash,
		Proposer:   proposer,
		TimeStamp:  time.Now().UTC(),
	}
}

// Add checkpoint to buffer or final headerBlocks
func (k Keeper) AddCheckpointToKey(ctx sdk.Context, start uint64, end uint64, root common.Hash, proposer common.Address, key []byte) sdk.Error {
	store := ctx.KVStore(k.checkpointKey)

	checkpointBuffer, _ := k.GetCheckpointFromBuffer(ctx)

	// Reject new checkpoint if checkpoint exists in buffer
	if bytes.Equal(key, BufferCheckpointKey) && !bytes.Equal(store.Get(BufferCheckpointKey), []byte("")) && time.Now().UTC().Before(checkpointBuffer.TimeStamp.Add(time.Minute*5)) {
		return ErrNoACK(k.codespace)
	}

	// Flush Checkpoint If 5 minutes have passed since it was added to buffer and NoAck received
	if bytes.Equal(key, BufferCheckpointKey) && !bytes.Equal(store.Get(BufferCheckpointKey), []byte("")) && time.Now().UTC().After(checkpointBuffer.TimeStamp.Add(time.Minute*5)) {
		k.FlushCheckpointBuffer(ctx)
	}

	// create Checkpoint block and marshall
	data := createBlock(start, end, root, proposer)
	out, err := json.Marshal(data)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}

	// store in key provided
	store.Set(key, []byte(out))

	return nil
}

// Flush Checkpoint Buffer
func (k Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.checkpointKey)
	store.Set(BufferCheckpointKey, []byte(""))
}

// Get checkpoint in buffer
func (k Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.checkpointKey)

	// Get checkpoint and unmarshall
	var checkpoint CheckpointBlockHeader
	err := json.Unmarshal(store.Get(BufferCheckpointKey), &checkpoint)

	return checkpoint, err
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

	// TODO maybe this needs to be set to 1
	// set to 0
	key := []byte(strconv.Itoa(int(0)))
	store.Set(ACKCountKey, key)
}

// appends prefix to headerNumber
func GetHeaderKey(headerNumber int) []byte {
	headerNumberBytes := strconv.Itoa(headerNumber)
	return append(HeaderBlockKey, headerNumberBytes...)
}

// gets last checkpoint , headerIndex = TotalACKs * ChildBlockInterval
func (k Keeper) GetLastCheckpoint(ctx sdk.Context) CheckpointBlockHeader {
	store := ctx.KVStore(k.checkpointKey)

	ACKs := k.GetACKCount(ctx)

	// fetch last checkpoint key (NumberOfACKs*ChildBlockInterval)
	lastCheckpointKey := (helper.GetConfig().ChildBlockInterval) * (ACKs)

	// fetch checkpoint and unmarshall
	var checkpoint CheckpointBlockHeader
	err := json.Unmarshal(store.Get(GetHeaderKey(lastCheckpointKey)), &checkpoint)
	if err != nil {
		CheckpointLogger.Error("Unable to fetch last checkpoint from store", "Key", lastCheckpointKey, "ACKCount", ACKs)
	}

	// return checkpoint
	return checkpoint
}
