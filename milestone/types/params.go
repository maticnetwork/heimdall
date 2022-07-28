package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/params/subspace"
)

// Default parameter values
const (
	DefaultSprintLength uint64 = 64
)

// Parameter keys
var (
	KeySprintLength = []byte("SprintLength")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	SprintLength uint64 `json:"sprint_length" yaml:"sprint_length"`
}

// NewParams creates a new Params object
func NewParams(
	sprintLength uint64,
) Params {
	return Params{
		SprintLength: sprintLength,
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
		{KeySprintLength, &p.SprintLength},
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
		SprintLength: DefaultSprintLength,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder

	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("SprintLength: %s\n", p.SprintLength))

	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if p.SprintLength == 0 {
		return fmt.Errorf("Sprint Length should be non-zero")
	}

	return nil
}
