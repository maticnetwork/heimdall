package supply

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/tendermint/tendermint/libs/log"

	auth "github.com/maticnetwork/heimdall/auth"
	bank "github.com/maticnetwork/heimdall/bank"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
)

// DefaultCodespace from the supply module
var DefaultCodespace sdk.CodespaceType = supplyTypes.ModuleName

// Keys for supply store
// Items are stored with the following key: values
//
// - 0x00: Supply
var (
	SupplyKey = []byte{0x00}
)

// Keeper of the supply store
type Keeper struct {
	cdc           *codec.Codec
	storeKey      sdk.StoreKey
	ak            auth.AccountKeeper
	bk            bank.Keeper
	paramSubspace subspace.Subspace
	permAddrs     map[string]supplyTypes.PermissionsForAddress
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramstore subspace.Subspace,
	maccPerms map[string][]string,
	ak auth.AccountKeeper,
	bk bank.Keeper,
) Keeper {

	// set the addresses
	permAddrs := make(map[string]supplyTypes.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = supplyTypes.NewPermissionsForAddress(name, perms)
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		ak:            ak,
		bk:            bk,
		paramSubspace: paramstore,
		permAddrs:     permAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", supplyTypes.ModuleName))
}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) (supply supplyTypes.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(SupplyKey)
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Context, supply supplyTypes.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(SupplyKey, b)
}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc supplyTypes.ModuleAccountInterface) error {
	permAddr := k.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}
	return nil
}
