package params

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/params/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the global paramstore
type Keeper struct {
	cdc       *codec.Codec
	key       sdk.StoreKey
	tkey      sdk.StoreKey
	codespace sdk.CodespaceType
	spaces    map[string]*subspace.Subspace
}

// NewKeeper constructs a params keeper
func NewKeeper(cdc *codec.Codec, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey, codespace sdk.CodespaceType) (k Keeper) {
	k = Keeper{
		cdc:       cdc,
		key:       key,
		tkey:      tkey,
		codespace: codespace,
		spaces:    make(map[string]*subspace.Subspace),
	}

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Subspace allocate subspace used for keepers
func (k Keeper) Subspace(s string) subspace.Subspace {
	_, ok := k.spaces[s]
	if ok {
		panic("subspace already occupied")
	}

	if s == "" {
		panic("cannot use empty string for subspace")
	}

	space := subspace.NewSubspace(k.cdc, k.key, k.tkey, s)
	k.spaces[s] = &space

	return space
}

// GetSubspace existing substore from keeper
func (k Keeper) GetSubspace(s string) (subspace.Subspace, bool) {
	space, ok := k.spaces[s]
	if !ok {
		return subspace.Subspace{}, false
	}
	return *space, ok
}
