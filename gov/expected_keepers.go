package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// SupplyKeeper defines the supply Keeper for module accounts
type SupplyKeeper interface {
	GetModuleAddress(name string) hmTypes.HeimdallAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyTypes.ModuleAccountInterface

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyTypes.ModuleAccountInterface)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr hmTypes.HeimdallAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr hmTypes.HeimdallAddress, recipientModule string, amt sdk.Coins) sdk.Error
}
