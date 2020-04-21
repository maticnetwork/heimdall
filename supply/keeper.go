package supply

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	auth "github.com/maticnetwork/heimdall/auth"
	bank "github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/params/subspace"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

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
	return ctx.Logger().With("module", supplyTypes.ModuleName)
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

//
// Account related methods
//

// GetModuleAddress returns an address based on the module name
func (k Keeper) GetModuleAddress(moduleName string) (addr hmTypes.HeimdallAddress) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return
	}
	addr = permAddr.GetAddress()
	return
}

// RemoveModuleAddress removes module address
func (k Keeper) RemoveModuleAddress(moduleName string) bool {
	_, ok := k.permAddrs[moduleName]
	if !ok {
		return false
	}
	delete(k.permAddrs, moduleName)
	return true
}

// GetModuleAddressAndPermissions returns an address and permissions based on the module name
func (k Keeper) GetModuleAddressAndPermissions(moduleName string) (addr hmTypes.HeimdallAddress, permissions []string) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return addr, permissions
	}
	return permAddr.GetAddress(), permAddr.GetPermissions()
}

// GetModuleAccountAndPermissions gets the module account from the auth account store and its
// registered permissions
func (k Keeper) GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (supplyTypes.ModuleAccountInterface, []string) {
	addr, perms := k.GetModuleAddressAndPermissions(moduleName)
	if addr.Empty() {
		return nil, []string{}
	}

	acc := k.ak.GetAccount(ctx, addr)
	if acc != nil {
		macc, ok := acc.(supplyTypes.ModuleAccountInterface)
		if !ok {
			panic("account is not a module account")
		}
		return macc, perms
	}

	// create a new module account
	macc := supplyTypes.NewEmptyModuleAccount(moduleName, perms...)
	maccI := (k.ak.NewAccount(ctx, macc)).(supplyTypes.ModuleAccountInterface) // set the account number
	k.SetModuleAccount(ctx, maccI)

	return maccI, perms
}

// GetModuleAccount gets the module account from the auth account store
func (k Keeper) GetModuleAccount(ctx sdk.Context, moduleName string) supplyTypes.ModuleAccountInterface {
	acc, _ := k.GetModuleAccountAndPermissions(ctx, moduleName)
	return acc
}

// SetModuleAccount sets the module account to the auth account store
func (k Keeper) SetModuleAccount(ctx sdk.Context, macc supplyTypes.ModuleAccountInterface) {
	k.ak.SetAccount(ctx, macc)
}

//
// Bank related methods
//

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context,
	senderModule string,
	recipientAddr hmTypes.HeimdallAddress,
	amt sdk.Coins,
) sdk.Error {
	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	// ignore tags
	return k.bk.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k Keeper) SendCoinsFromModuleToModule(
	ctx sdk.Context,
	senderModule string,
	recipientModule string,
	amt sdk.Coins,
) sdk.Error {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return supplyTypes.ErrNoAccountCreated(supplyTypes.DefaultCodespace)
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(
	ctx sdk.Context,
	senderAddr hmTypes.HeimdallAddress,
	recipientModule string,
	amt sdk.Coins,
) sdk.Error {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return supplyTypes.ErrNoAccountCreated(supplyTypes.DefaultCodespace)
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}
