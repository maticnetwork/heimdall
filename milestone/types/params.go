package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/params/subspace"
)

// Default parameter values
const (
	DefaultMilestoneLength uint64 = 64
)

// Parameter keys
var (
	KeyMilestoneLength = []byte("MilestoneLength")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	MilestoneLength uint64 `json:"milestone_length" yaml:"milestone_length"`
}

type Count struct {
	Count uint64 `json:"count" yaml:"count"`
}

// NewParams creates a new Params object
func NewParams(
	milestoneLength uint64,
) Params {
	return Params{
		MilestoneLength: milestoneLength,
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyMilestoneLength, &p.MilestoneLength},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)

	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MilestoneLength: DefaultMilestoneLength,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder

	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("MilestoneLength: %s\n", p.MilestoneLength))

	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if p.MilestoneLength == 0 {
		return fmt.Errorf("Milestone Length should be non-zero")
	}

	return nil
}
