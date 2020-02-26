package bank

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/bank/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	// DefaultValue default value
	DefaultValue = []byte{0x01}
	// ValidatorTopupKey represents validator topup key
	ValidatorTopupKey = []byte{0x80} // prefix for each key to a validator
	// TopupSequencePrefixKey represents topup sequence prefix key
	TopupSequencePrefixKey = []byte{0x81}
)

// ValidatorManager manager to access validator info
type ValidatorManager interface {
	// AddFeeToDividendAccount add fee to dividend account
	AddFeeToDividendAccount(ctx sdk.Context, valID hmTypes.ValidatorID, fee *big.Int) sdk.Error
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
	paramSpace params.Subspace
	// account keeper
	ak auth.AccountKeeper
	// validator manager
	vm ValidatorManager
}

// NewKeeper returns a new Keeper
func NewKeeper(
	cdc *codec.Codec,
	key sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	ak auth.AccountKeeper,
	vm ValidatorManager,
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
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt hmTypes.Coins,
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
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt hmTypes.Coins,
) (hmTypes.Coins, sdk.Error) {

	if !amt.IsValid() {
		return nil, sdk.ErrInvalidCoins(amt.String())
	}

	oldCoins, spendableCoins := hmTypes.NewCoins(), hmTypes.NewCoins()

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
	ctx sdk.Context, addr hmTypes.HeimdallAddress, amt hmTypes.Coins,
) (hmTypes.Coins, sdk.Error) {

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

// InputOutputCoins handles a list of inputs and outputs
func (keeper Keeper) InputOutputCoins(
	ctx sdk.Context, inputs []types.Input, outputs []types.Output,
) sdk.Error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of Coins.
	if err := types.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	for _, in := range inputs {
		_, err := keeper.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(types.AttributeKeySender, in.Address.String()),
			),
		)
	}

	for _, out := range outputs {
		_, err := keeper.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeTransfer,
				sdk.NewAttribute(types.AttributeKeyRecipient, out.Address.String()),
			),
		)
	}

	return nil
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr hmTypes.HeimdallAddress, toAddr hmTypes.HeimdallAddress, amt hmTypes.Coins,
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
func (keeper Keeper) GetCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress) hmTypes.Coins {
	acc := keeper.ak.GetAccount(ctx, addr)
	if acc == nil {
		return hmTypes.NewCoins()
	}
	return acc.GetCoins()
}

// HasCoins returns whether or not an account has at least amt coins.
func (keeper Keeper) HasCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress, amt hmTypes.Coins) bool {
	return keeper.GetCoins(ctx, addr).IsAllGTE(amt)
}

//
// Topup methods
//

// GetTopupKey drafts the topup key for address
func GetTopupKey(address []byte) []byte {
	return append(ValidatorTopupKey, address...)
}

// GetTopupSequenceKey drafts topup sequence for address
func GetTopupSequenceKey(sequence uint64) []byte {
	return append(TopupSequencePrefixKey, []byte(strconv.FormatUint(sequence, 10))...)
}

// GetValidatorTopup returns validator toptup information
func (keeper Keeper) GetValidatorTopup(ctx sdk.Context, addr hmTypes.HeimdallAddress) (*types.ValidatorTopup, error) {
	store := ctx.KVStore(keeper.key)

	// check if topup exists
	key := GetTopupKey(addr.Bytes())
	if !store.Has(key) {
		return nil, nil
	}

	// unmarshall validator and return
	validatorTopup, err := types.UnmarshallValidatorTopup(keeper.cdc, store.Get(key))
	if err != nil {
		return nil, err
	}

	// return true if validator
	return &validatorTopup, nil
}

// SetValidatorTopup sets validator topup object
func (keeper Keeper) SetValidatorTopup(ctx sdk.Context, addr hmTypes.HeimdallAddress, validatorTopup types.ValidatorTopup) error {
	store := ctx.KVStore(keeper.key)

	// validator topup
	bz, err := types.MarshallValidatorTopup(keeper.cdc, validatorTopup)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(GetTopupKey(addr.Bytes()), bz)
	keeper.Logger(ctx).Debug("Validator topup stored", "key", hex.EncodeToString(GetTopupKey(addr.Bytes())), "totalTopups", validatorTopup.Copy().TotalTopups)

	return nil
}

// SetTopupSequence sets mapping for sequence id to bool
func (keeper Keeper) SetTopupSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(keeper.key)
	store.Set(GetTopupSequenceKey(sequence), DefaultValue)
}

// HasTopupSequence checks if topup already exists
func (keeper Keeper) HasTopupSequence(ctx sdk.Context, sequence uint64) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(GetTopupSequenceKey(sequence))
}
