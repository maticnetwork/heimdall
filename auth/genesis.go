package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	staking "github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"

	"errors"
)

// GenesisState - all auth state that must be provided at genesis
type GenesisState struct {
	CollectedFees types.Coins      `json:"collected_fees"`
	Params        authTypes.Params `json:"params"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params authTypes.Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(authTypes.DefaultParams())
}

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, data GenesisState) {
	ak.SetParams(ctx, data.Params)
}

// InitAccounts
func InitAcccounts(ctx sdk.Context, ak AccountKeeper, data staking.GenesisState) {
	var newValSet types.ValidatorSet
	for _, validator := range data.Validators {
		if ok := newValSet.Add(&validator); !ok {
			panic(errors.New("Error while addings new validator"))
		} else {
			// Add individual validator to state
			ak.NewAccountWithAddress(ctx, types.BytesToHeimdallAddress(validator.Signer.Bytes()))
			fmt.Printf("added new accounts %v", ak.GetAllAccounts(ctx))
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper) GenesisState {
	params := ak.GetParams(ctx)

	return NewGenesisState(params)
}

// ValidateGenesis performs basic validation of auth genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.Params.TxSigLimit == 0 {
		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
	}
	if data.Params.SigVerifyCostED25519 == 0 {
		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
	}
	if data.Params.SigVerifyCostSecp256k1 == 0 {
		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
	}
	if data.Params.MaxMemoCharacters == 0 {
		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
	}
	if data.Params.TxSizeCostPerByte == 0 {
		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
	}
	return nil
}
