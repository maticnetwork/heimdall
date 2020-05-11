package bank

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	// DefaultValue default value
	DefaultValue = []byte{0x01}
)

// TODO: Remove this later
// ModuleCommunicator manager to access validator info
type ModuleCommunicator interface {
	// GetValidatorFromValID get validator from validator id
	GetValidatorFromValID(ctx sdk.Context, valID hmTypes.ValidatorID) (validator hmTypes.Validator, ok bool)
}

// Keeper manages transfers between accounts
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// code space
	codespace sdk.CodespaceType
	// param subspace
	paramSpace subspace.Subspace
	// account keeper
	ak auth.AccountKeeper
	// module manager
	vm ModuleCommunicator
}

// NewKeeper returns a new Keeper
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	ak auth.AccountKeeper,
	vm ModuleCommunicator,
) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())
	return Keeper{
		key:        key,
		cdc:        cdc,
		codespace:  codespace,
		paramSpace: ps,
		ak:         ak,
		vm:         vm,
	}
}

// Codespace returns the keeper's codespace.
func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

// Logger returns a module-specific logger
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// SetCoins sets the coins at the addr.
func (keeper Keeper) SetCoins(
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt sdk.Coins,
) sdk.Error {

	if !amt.IsValid() && !amt.IsZero() {
		return sdk.ErrInvalidCoins(amt.String())
	}

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		acc = keeper.ak.NewAccountWithAddress(ctx, addr)
	}

	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	keeper.ak.SetAccount(ctx, acc)
	return nil
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper Keeper) SubtractCoins(
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := sdk.NewCoins(), sdk.NewCoins()

	acc := keeper.ak.GetAccount(ctx, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := keeper.SetCoins(ctx, addr, newCoins)

	return newCoins, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper Keeper) AddCoins(
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt sdk.Coins,
) (sdk.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins := keeper.GetCoins(ctx, addr)
	newCoins := oldCoins.Add(amt)

	if newCoins.IsAnyNegative() {
		return amt, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := keeper.SetCoins(ctx, addr, newCoins)
	return newCoins, err
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr hmTypes.HeimdallAddress, toAddr hmTypes.HeimdallAddress, amt sdk.Coins,
) sdk.Error {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
		),
	})

	_, err := keeper.SubtractCoins(ctx, fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = keeper.AddCoins(ctx, toAddr, amt)
	if err != nil {
		return err
	}

	return nil
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper Keeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, types.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper Keeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, types.ParamStoreKeySendEnabled, &enabled)
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress) sdk.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress, amt sdk.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}
