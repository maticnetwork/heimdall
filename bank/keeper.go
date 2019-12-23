package bank

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/types"
)

var (
	// ValidatorTopupKey represents validator topup key
	ValidatorTopupKey = []byte{0x80} // prefix for each key to a validator
)

// Keeper manages transfers between accounts
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// code space
	codespace sdk.CodespaceType
	// param subspace
	paramSpace params.Subspace
	// account keeper
	ak auth.AccountKeeper
}

// NewKeeper returns a new Keeper
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	ak auth.AccountKeeper,
) Keeper {
	ps := paramSpace.WithKeyTable(bankTypes.ParamKeyTable())
	return Keeper{
		key:        key,
		cdc:        cdc,
		codespace:  codespace,
		paramSpace: ps,
		ak:         ak,
	}
}

// GetTopupKey drafts the validator key for addresses
func GetTopupKey(address []byte) []byte {
	return append(ValidatorTopupKey, address...)
}

// Codespace returns the keeper's codespace.
func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

// Logger returns a module-specific logger
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", bankTypes.ModuleName)
}

// SetCoins sets the coins at the addr.
func (keeper Keeper) SetCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) sdk.Error {

	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	return setCoins(ctx, keeper.ak, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper Keeper) SubtractCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}
	return subtractCoins(ctx, keeper.ak, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper Keeper) AddCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}
	return addCoins(ctx, keeper.ak, addr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper Keeper) InputOutputCoins(
	ctx sdk.Context, inputs []bankTypes.Input, outputs []bankTypes.Output,
) (sdk.Tags, sdk.Error) {

	return inputOutputCoins(ctx, keeper.ak, inputs, outputs)
}

// GetTopup returns validator toptup information
func (keeper Keeper) GetTopup(ctx sdk.Context, addr types.HeimdallAddress) (*bankTypes.ValidatorTopup, error) {
	store := ctx.KVStore(keeper.key)

	// check if topup exists
	key := GetTopupKey(addr.Bytes())
	if !store.Has(key) {
		return nil, nil
	}

	// unmarshall validator and return
	validatorTopup, err := bankTypes.UnmarshallValidatorTopup(keeper.cdc, store.Get(key))
	if err != nil {
		return nil, err
	}

	// return true if validator
	return &validatorTopup, nil
}

// SetTopup sets validator topup object
func (keeper Keeper) SetTopup(ctx sdk.Context, addr types.HeimdallAddress, validatorTopup bankTypes.ValidatorTopup) error {
	store := ctx.KVStore(keeper.key)

	// validator topup
	bz, err := bankTypes.MarshallValidatorTopup(keeper.cdc, validatorTopup)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(GetTopupKey(addr.Bytes()), bz)
	keeper.Logger(ctx).Debug("Validator topup stored", "key", hex.EncodeToString(GetTopupKey(addr.Bytes())), "totalTopups", validatorTopup.Copy().TotalTopups)

	return nil
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt types.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return sendCoins(ctx, keeper.ak, fromAddr, toAddr, amt)
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper Keeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper Keeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCoins(ctx sdk.Context, addr types.HeimdallAddress) types.Coins {
	return getCoins(ctx, keeper.ak, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) bool {
	return hasCoins(ctx, keeper.ak, addr, amt)
}

func getCoins(ctx sdk.Context, am auth.AccountKeeper, addr types.HeimdallAddress) types.Coins {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return types.NewCoins()
	}
	return acc.GetCoins()
}

func setCoins(ctx sdk.Context, am auth.AccountKeeper, addr types.HeimdallAddress, amt types.Coins) sdk.Error {
	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCoins(amt)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// HasCoins returns whether or not an account has at least amt coins.
func hasCoins(ctx sdk.Context, am auth.AccountKeeper, addr types.HeimdallAddress, amt types.Coins) bool {
	return getCoins(ctx, am, addr).IsAllGTE(amt)
}

func getAccount(ctx sdk.Context, ak auth.AccountKeeper, addr types.HeimdallAddress) authTypes.Account {
	return ak.GetAccount(ctx, addr)
}

func setAccount(ctx sdk.Context, ak auth.AccountKeeper, acc authTypes.Account) {
	ak.SetAccount(ctx, acc)
}

// subtractCoins subtracts amt coins from an account with the given address addr.
func subtractCoins(ctx sdk.Context, ak auth.AccountKeeper, addr types.HeimdallAddress, amt types.Coins) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := types.NewCoins(), types.NewCoins()

	acc := getAccount(ctx, ak, addr)
	if acc != nil {
		oldCoins = acc.GetCoins()
		spendableCoins = acc.SpendableCoins(ctx.BlockHeader().Time)
	}

	// So the check here is sufficient instead of subtracting from oldCoins.
	_, hasNeg := spendableCoins.SafeSub(amt)
	if hasNeg {
		return amt, nil, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", spendableCoins, amt),
		)
	}

	newCoins := oldCoins.Sub(amt) // should not panic as spendable coins was already checked
	err := setCoins(ctx, ak, addr, newCoins)
	tags := sdk.NewTags(TagKeySender, addr.String())

	return newCoins, tags, err
}

// AddCoins adds amt to the coins at the addr.
func addCoins(ctx sdk.Context, am auth.AccountKeeper, addr types.HeimdallAddress, amt types.Coins) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins := getCoins(ctx, am, addr)
	newCoins := oldCoins.Add(amt)

	if newCoins.IsAnyNegative() {
		return amt, nil, sdk.ErrInsufficientCoins(
			fmt.Sprintf("insufficient account funds; %s < %s", oldCoins, amt),
		)
	}

	err := setCoins(ctx, am, addr, newCoins)
	tags := sdk.NewTags(TagKeyRecipient, addr.String())

	return newCoins, tags, err
}

// SendCoins moves coins from one account to another
// Returns ErrInvalidCoins if amt is invalid.
func sendCoins(ctx sdk.Context, am auth.AccountKeeper, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt types.Coins) (sdk.Tags, sdk.Error) {
	// Safety check ensuring that when sending coins the keeper must maintain the
	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	_, subTags, err := subtractCoins(ctx, am, fromAddr, amt)
	if err != nil {
		return nil, err
	}

	_, addTags, err := addCoins(ctx, am, toAddr, amt)
	if err != nil {
		return nil, err
	}

	return subTags.AppendTags(addTags), nil
}

// InputOutputCoins handles a list of inputs and outputs
// NOTE: Make sure to revert state changes from tx on error
func inputOutputCoins(ctx sdk.Context, am auth.AccountKeeper, inputs []bankTypes.Input, outputs []bankTypes.Output) (sdk.Tags, sdk.Error) {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := bankTypes.ValidateInputsOutputs(inputs, outputs); err != nil {
		return nil, err
	}

	allTags := sdk.EmptyTags()

	for _, in := range inputs {
		_, tags, err := subtractCoins(ctx, am, in.Address, in.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	for _, out := range outputs {
		_, tags, err := addCoins(ctx, am, out.Address, out.Coins)
		if err != nil {
			return nil, err
		}
		allTags = allTags.AppendTags(tags)
	}

	return allTags, nil
}
