package checkpoint


import (
sdk "github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/wire"
	"strconv"
	"encoding/json"
	"fmt"
)

type Keeper struct {
	checkpointKey     sdk.StoreKey
	cdc          *wire.Codec
	//validatorSet sdk.ValidatorSet

	// codespace
	codespace sdk.CodespaceType
}
func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		checkpointKey:   key,
		cdc:        cdc,
		codespace:  codespace,
	}
	return keeper
}

func (k Keeper) addCheckpoint(ctx sdk.Context, data []BlockHeader)  {
	store := ctx.KVStore(k.checkpointKey)
	out, err := json.Marshal(data)
	if err != nil {
		panic (err)
	}
	//TODO add block data validation

	fmt.Printf("Block data to be inserted with key %v",[]byte(strconv.Itoa(int(ctx.BlockHeight()))))
	fmt.Printf("Block data to be inserted is %v",out)
	store.Set([]byte(strconv.Itoa(int(ctx.BlockHeight()))),[]byte(out))
}
func(k Keeper) getCheckpoint(ctx sdk.Context,key int64) []byte {
	store := ctx.KVStore(k.checkpointKey)
	getKey := []byte(strconv.Itoa(int(key)))
	return store.Get(getKey)

}