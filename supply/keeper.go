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
	"github.com/maticnetwork/heimdall/types"
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

//
// Account related methods
//

// GetModuleAddress returns an address based on the module name
func (k Keeper) GetModuleAddress(moduleName string) (addr types.HeimdallAddress) {
	permAddr, ok := k.permAddrs[moduleName]
	if !ok {
		return
	}
	addr = permAddr.GetAddress()
	return
}

// GetModuleAddressAndPermissions returns an address and permissions based on the module name
func (k Keeper) GetModuleAddressAndPermissions(moduleName string) (addr types.HeimdallAddress, permissions []string) {
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
	recipientAddr types.HeimdallAddress,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {
	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	// ignore tags
	return k.bk.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k Keeper) SendCoinsFromModuleToModule(
	ctx sdk.Context,
	senderModule string,
	recipientModule string,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {

	senderAddr := k.GetModuleAddress(senderModule)
	if senderAddr.Empty() {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return nil, supplyTypes.ErrNoAccountCreated(supplyTypes.DefaultCodespace)
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(
	ctx sdk.Context,
	senderAddr types.HeimdallAddress,
	recipientModule string,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return nil, supplyTypes.ErrNoAccountCreated(supplyTypes.DefaultCodespace)
	}

	return k.bk.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// DelegateCoinsFromAccountToModule delegates coins and transfers
// them from a delegator account to a module account
func (k Keeper) DelegateCoinsFromAccountToModule(
	ctx sdk.Context,
	senderAddr types.HeimdallAddress,
	recipientModule string,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {

	// create the account if it doesn't yet exist
	recipientAcc := k.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		return nil, supplyTypes.ErrNoAccountCreated(supplyTypes.DefaultCodespace)
	}

	if !recipientAcc.HasPermission(supplyTypes.Staking) {
		return nil, supplyTypes.ErrNoPermission(supplyTypes.DefaultCodespace)
	}

	return k.bk.DelegateCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

// UndelegateCoinsFromModuleToAccount undelegates the unbonding coins and transfers
// them from a module account to the delegator account
func (k Keeper) UndelegateCoinsFromModuleToAccount(
	ctx sdk.Context,
	senderModule string,
	recipientAddr types.HeimdallAddress,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {

	acc := k.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", senderModule))
	}

	if !acc.HasPermission(supplyTypes.Staking) {
		return nil, supplyTypes.ErrNoPermission(supplyTypes.DefaultCodespace)
	}

	return k.bk.UndelegateCoins(ctx, acc.GetAddress(), recipientAddr, amt)
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintCoins(
	ctx sdk.Context,
	moduleName string,
	amt types.Coins,
) (sdk.Tags, sdk.Error) {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(supplyTypes.Minter) {
		return nil, supplyTypes.ErrNoPermission(supplyTypes.DefaultCodespace)
	}

	_, tags, err := k.bk.AddCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		return tags, err
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply.Inflate(amt)
	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted %s from %s module account", amt.String(), moduleName))

	return nil, nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amt types.Coins) (sdk.Tags, sdk.Error) {

	// create the account if it doesn't yet exist
	acc := k.GetModuleAccount(ctx, moduleName)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("module account %s does not exist", moduleName))
	}

	if !acc.HasPermission(supplyTypes.Burner) {
		return nil, supplyTypes.ErrNoPermission(supplyTypes.DefaultCodespace)
	}

	_, _, err := k.bk.SubtractCoins(ctx, acc.GetAddress(), amt)
	if err != nil {
		return nil, err
	}

	// update total supply
	supply := k.GetSupply(ctx)
	supply.Deflate(amt)
	k.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned %s from %s module account", amt.String(), moduleName))

	return nil, nil
}
