package bank

import (
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

var _ Keeper = (*BaseKeeper)(nil)

// Keeper defines a module interface that facilitates the transfer of coins
// between accounts.
type Keeper interface {
	SendKeeper

	SetCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) sdk.Error
	SubtractCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) (types.Coins, sdk.Tags, sdk.Error)
	AddCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) (types.Coins, sdk.Tags, sdk.Error)
	InputOutputCoins(ctx sdk.Context, inputs []bankTypes.Input, outputs []bankTypes.Output) (sdk.Tags, sdk.Error)
}

// BaseKeeper manages transfers between accounts. It implements the Keeper interface.
type BaseKeeper struct {
	BaseSendKeeper

	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// param subspace
	paramSpace params.Subspace
	// account keeper
	ak auth.AccountKeeper
}

// NewBaseKeeper returns a new BaseKeeper
func NewBaseKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	ak auth.AccountKeeper,
) BaseKeeper {
	ps := paramSpace.WithKeyTable(bankTypes.ParamKeyTable())
	return BaseKeeper{
		key:            key,
		cdc:            cdc,
		paramSpace:     ps,
		ak:             ak,
		BaseSendKeeper: NewBaseSendKeeper(ak, ps, codespace),
	}
}

// Logger returns a module-specific logger
func (keeper BaseKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", bankTypes.ModuleName)
}

// SetCoins sets the coins at the addr.
func (keeper BaseKeeper) SetCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) sdk.Error {

	if !amt.IsValid() {
		return sdk.ErrInvalidCoins(amt.String())
	}
	return setCoins(ctx, keeper.ak, addr, amt)
}

// SubtractCoins subtracts amt from the coins at the addr.
func (keeper BaseKeeper) SubtractCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}
	return subtractCoins(ctx, keeper.ak, addr, amt)
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseKeeper) AddCoins(
	ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins,
) (types.Coins, sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, nil, sdk.ErrInvalidCoins(amt.String())
	}
	return addCoins(ctx, keeper.ak, addr, amt)
}

// InputOutputCoins handles a list of inputs and outputs
func (keeper BaseKeeper) InputOutputCoins(
	ctx sdk.Context, inputs []bankTypes.Input, outputs []bankTypes.Output,
) (sdk.Tags, sdk.Error) {

	return inputOutputCoins(ctx, keeper.ak, inputs, outputs)
}

// between accounts without the possibility of creating coins.
type SendKeeper interface {
	ViewKeeper

	SendCoins(ctx sdk.Context, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt types.Coins) (sdk.Tags, sdk.Error)

	GetSendEnabled(ctx sdk.Context) bool
	SetSendEnabled(ctx sdk.Context, enabled bool)
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	ak         auth.AccountKeeper
	paramSpace params.Subspace
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(ak auth.AccountKeeper,
	paramSpace params.Subspace, codespace sdk.CodespaceType) BaseSendKeeper {

	return BaseSendKeeper{
		BaseViewKeeper: NewBaseViewKeeper(ak, codespace),
		ak:             ak,
		paramSpace:     paramSpace,
	}
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(
	ctx sdk.Context, fromAddr types.HeimdallAddress, toAddr types.HeimdallAddress, amt types.Coins,
) (sdk.Tags, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}
	return sendCoins(ctx, keeper.ak, fromAddr, toAddr, amt)
}

// GetSendEnabled returns the current SendEnabled
// nolint: errcheck
func (keeper BaseSendKeeper) GetSendEnabled(ctx sdk.Context) bool {
	var enabled bool
	keeper.paramSpace.Get(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
	return enabled
}

// SetSendEnabled sets the send enabled
func (keeper BaseSendKeeper) SetSendEnabled(ctx sdk.Context, enabled bool) {
	keeper.paramSpace.Set(ctx, bankTypes.ParamStoreKeySendEnabled, &enabled)
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetCoins(ctx sdk.Context, addr types.HeimdallAddress) types.Coins
	HasCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) bool

	Codespace() sdk.CodespaceType
}

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	ak        auth.AccountKeeper
	codespace sdk.CodespaceType
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(ak auth.AccountKeeper, codespace sdk.CodespaceType) BaseViewKeeper {
	return BaseViewKeeper{ak: ak, codespace: codespace}
}

// GetCoins returns the coins at the addr.
func (keeper BaseViewKeeper) GetCoins(ctx sdk.Context, addr types.HeimdallAddress) types.Coins {
	return getCoins(ctx, keeper.ak, addr)
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper BaseViewKeeper) HasCoins(ctx sdk.Context, addr types.HeimdallAddress, amt types.Coins) bool {
	return hasCoins(ctx, keeper.ak, addr, amt)
}

// Codespace returns the keeper's codespace.
func (keeper BaseViewKeeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
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
