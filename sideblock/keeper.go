package sideBlock


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"encoding/json"
	"strings"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *wire.Codec
	//validatorSet sdk.ValidatorSet

	// codespace
	codespace sdk.CodespaceType
}
type Block struct {
	BlockHash string
	TxRoot string
	ReceiptRoot string
}

func createBlock(blockHash string,txroot string,rRoot string) Block  {
	return Block{
		BlockHash:blockHash,
		TxRoot:txroot,
		ReceiptRoot:rRoot,
	}
}
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		codespace:  codespace,
	}
	return keeper
}

// will get you the block details provided with the block hash , cool !
func (k Keeper) getBlock(ctx sdk.Context, blockhash string) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get([]byte(blockhash))
}
// will store the block hash under the block hash for now , later will store the block struct with key as block hash
func (k Keeper) addBlock(ctx sdk.Context,blockHash string,txroot string,rRoot string)  {
	//logger := ctx.Logger().With("module", "x/sideBlock")
	store := ctx.KVStore(k.storeKey)
	// TODO replace the second param with block struct and first will remain block hash
	// we are using block hash as the key here !
	block := createBlock(blockHash,txroot,rRoot)
	addBlockHash(blockHash,ctx,k)
	//converts to json to add to store
	out, err := json.Marshal(block)
	if err != nil {
		panic (err)
	}

	store.Set([]byte(blockHash),[]byte(out))
	//logger.Info("oh okay so the logs so work, ctx is %s", ctx)
}

func addBlockHash(blockhash string,ctx sdk.Context,k Keeper){
	store := ctx.KVStore(k.storeKey)
	blocks := store.Get([]byte("latestBlockHashes"))
	//fmt.Printf("blockhashes are %v",string(blocks))
	blocksString := string(blocks)
	newBlocksString := blocksString+" "+blockhash
	//fmt.Printf("new block string is %v",newBlocksString)
	// TODO check if the blockhash already exists , dont insert then
	store.Set([]byte("latestBlockHashes"),[]byte(newBlocksString))
	// grab the bytes and append the new space separated blockhash after converting to byte
	//store.Set([]byte("latestBlockHashes"),)

}
func InitGenesis(ctx sdk.Context, k Keeper)  {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte("latestBlockHashes"),[]byte("blockhash0"))
}
func FlushBlockHashesKey(ctx sdk.Context,k Keeper)  {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte("latestBlockHashes"),[]byte("blockhash0"))

}
func GetBlocksAfterCheckpoint(ctx sdk.Context,k Keeper) string {
	store := ctx.KVStore(k.storeKey)
	blocks := store.Get([]byte("latestBlockHashes"))
	blocksString := string(blocks)
	blockHashList := strings.Fields(blocksString)
	checkpointBlockDetails :=string("")
	for i:=1 ; i< len(blockHashList)-1; i++ {
		blockHashString := blockHashList[i]
		checkpointBlockDetails= checkpointBlockDetails + string(k.getBlock(ctx,blockHashString))

	}
	//fmt.Printf("Checkpoint details %v",checkpointBlockDetails)
	return checkpointBlockDetails

}