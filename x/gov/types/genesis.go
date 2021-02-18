package types

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	// DefaultIndex is the default capability global index
	DefaultIndex uint64 = 1
	// DefaultPeriod default period for deposits & voting
	DefaultPeriod time.Duration = 86400 * 2 * time.Second // 2 days
)

// NewGenesisState creates a new genesis state for the governance module
func NewGenesisState(startingProposalID uint64, dp DepositParams, vp VotingParams, tp TallyParams) GenesisState {
	return GenesisState{
		StartingProposalId: startingProposalID,
		DepositParams:      dp,
		VotingParams:       vp,
		TallyParams:        tp,
	}
}

// DefaultGenesis get raw genesis raw message for testing
func DefaultGenesis() *GenesisState {
	minDepositTokens := sdk.NewIntFromBigInt(new(big.Int).Mul(big.NewInt(10), hmTypes.CoinDecimals))
	return &GenesisState{
		StartingProposalId: 1,
		DepositParams: DepositParams{
			MinDeposit:       sdk.Coins{sdk.NewCoin(hmTypes.FeeToken, minDepositTokens)},
			MaxDepositPeriod: DefaultPeriod,
		},
		VotingParams: VotingParams{
			VotingPeriod: DefaultPeriod,
		},
		TallyParams: TallyParams{
			Quorum:    sdk.NewDecWithPrec(334, 3),
			Threshold: sdk.NewDecWithPrec(5, 1),
			Veto:      sdk.NewDecWithPrec(334, 3),
		},
	}
}

// Equal checks whether 2 GenesisState structs are equivalent.
func (data GenesisState) Equal(data2 GenesisState) bool {
	b1 := ModuleCdc.MustMarshalBinaryBare(&data)
	b2 := ModuleCdc.MustMarshalBinaryBare(&data2)
	return bytes.Equal(b1, b2)
}

// IsEmpty returns if a GenesisState is empty or has data in it
func (data GenesisState) IsEmpty() bool {
	emptyGenState := GenesisState{}
	return data.Equal(emptyGenState)
}

// Validate checks if parameters are within valid ranges
func (gs GenesisState) Validate() error {
	threshold := gs.TallyParams.Threshold
	if threshold.IsNegative() || threshold.GT(sdk.OneDec()) {
		return fmt.Errorf("Governance vote threshold should be positive and less or equal to one, is %s",
			threshold.String())
	}

	veto := gs.TallyParams.Veto
	if veto.IsNegative() || veto.GT(sdk.OneDec()) {
		return fmt.Errorf("Governance vote veto threshold should be positive and less or equal to one, is %s",
			veto.String())
	}

	if !gs.DepositParams.MinDeposit.IsValid() {
		return fmt.Errorf("Governance deposit amount must be a valid sdk.Coins amount, is %s",
			gs.DepositParams.MinDeposit.String())
	}

	return nil
}
