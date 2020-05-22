package types

import (
	"fmt"
	"time"

	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Evidence defines the contract which concrete evidence types of misbehavior
// must implement.
type Evidence interface {
	Route() string
	Type() string
	String() string
	Hash() []byte
	ValidateBasic() error

	// The consensus address of the malicious validator at time of infraction
	GetConsensusAddress() sdk.ConsAddress

	// Height at which the infraction occurred
	GetHeight() int64

	// The total power of the malicious validator at time of infraction
	GetValidatorPower() int64

	// The total validator set power at time of infraction
	GetTotalPower() int64
}

// MsgSubmitEvidence defines the specific interface a concrete message must
// implement in order to process submitted evidence. The concrete MsgSubmitEvidence
// must be defined at the application-level.
type MsgSubmitEvidence interface {
	sdk.Msg

	GetEvidence() Evidence
	GetSubmitter() sdk.AccAddress
}

// Equivocation implements the Evidence interface and defines evidence of double
// signing misbehavior.
type Equivocation struct {
	Height           int64     `json:"height,omitempty"`
	Time             time.Time `json:"time"`
	Power            int64     `json:"power,omitempty"`
	ConsensusAddress []byte    `json:"address,omitempty"`
}

// Evidence type constants
const (
	RouteEquivocation = "equivocation"
	TypeEquivocation  = "equivocation"
)

// var _ exported.Evidence = (*Equivocation)(nil)

// Route returns the Evidence Handler route for an Equivocation type.
func (e Equivocation) Route() string { return RouteEquivocation }

// Type returns the Evidence Handler type for an Equivocation type.
func (e Equivocation) Type() string { return TypeEquivocation }

func (e Equivocation) String() string {
	bz, _ := yaml.Marshal(e)
	return string(bz)
}

// Hash returns the hash of an Equivocation object.
func (e Equivocation) Hash() []byte {
	return tmhash.Sum(ModuleCdc.MustMarshalBinaryBare(&e))
}

// ValidateBasic performs basic stateless validation checks on an Equivocation object.
func (e Equivocation) ValidateBasic() error {
	if e.Time.IsZero() {
		return fmt.Errorf("invalid equivocation time: %s", e.Time)
	}
	if e.Height < 1 {
		return fmt.Errorf("invalid equivocation height: %d", e.Height)
	}
	if e.Power < 1 {
		return fmt.Errorf("invalid equivocation validator power: %d", e.Power)
	}
	if e.ConsensusAddress == nil {
		return fmt.Errorf("invalid equivocation validator consensus address: %s", e.ConsensusAddress)
	}

	return nil
}

// GetConsensusAddress returns the validator's consensus address at time of the
// Equivocation infraction.
func (e Equivocation) GetConsensusAddress() sdk.ConsAddress {
	return e.ConsensusAddress
}

// GetHeight returns the height at time of the Equivocation infraction.
func (e Equivocation) GetHeight() int64 {
	return e.Height
}

// GetTime returns the time at time of the Equivocation infraction.
func (e Equivocation) GetTime() time.Time {
	return e.Time
}

// GetValidatorPower returns the validator's power at time of the Equivocation
// infraction.
func (e Equivocation) GetValidatorPower() int64 {
	return e.Power
}

// GetTotalPower is a no-op for the Equivocation type.
func (e Equivocation) GetTotalPower() int64 { return 0 }

// ConvertDuplicateVoteEvidence converts a Tendermint concrete Evidence type to
// SDK Evidence using Equivocation as the concrete type.
func ConvertDuplicateVoteEvidence(dupVote abci.Evidence) Evidence {
	return Equivocation{
		Height:           dupVote.Height,
		Power:            dupVote.Validator.Power,
		ConsensusAddress: sdk.ConsAddress(dupVote.Validator.Address),
		Time:             dupVote.Time,
	}
}
