package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/params/subspace"
)

const (
	// SlotCost cost for validator
	SlotCost int64 = 1
)

// Default parameter values
const (
	DefaultSprintDuration    uint64 = 64
	DefaultSpanDuration      uint64 = 100 * DefaultSprintDuration
	DefaultFirstSpanDuration uint64 = 256
	DefaultProducerCount     uint64 = 4
)

// Parameter keys
var (
	KeySprintDuration = []byte("SprintDuration")
	KeySpanDuration   = []byte("SpanDuration")
	KeyProducerCount  = []byte("ProducerCount")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	SprintDuration uint64 `json:"sprint_duration" yaml:"sprint_duration"` // sprint duration
	SpanDuration   uint64 `json:"span_duration" yaml:"span_duration"`     // span duration ie number of blocks for which val set is frozen on heimdall
	ProducerCount  uint64 `json:"producer_count" yaml:"producer_count"`   // producer count per span
}

// NewParams creates a new Params object
func NewParams(sprintDuration uint64, spanDuration uint64, producerCount uint64) Params {
	return Params{
		SprintDuration: sprintDuration,
		SpanDuration:   spanDuration,
		ProducerCount:  producerCount,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeySprintDuration, &p.SprintDuration},
		{KeySpanDuration, &p.SpanDuration},
		{KeyProducerCount, &p.ProducerCount},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("SprintDuration: %d\n", p.SprintDuration))
	sb.WriteString(fmt.Sprintf("SpanDuration: %d\n", p.SpanDuration))
	sb.WriteString(fmt.Sprintf("ProducerCount: %d\n", p.ProducerCount))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateSprintDuration(p.SprintDuration); err != nil {
		return err
	}

	if err := validateSpanDuration(p.SprintDuration); err != nil {
		return err
	}

	if err := validateProducerCount(p.SprintDuration); err != nil {
		return err
	}

	return nil
}

//
// Extra functions
//

// ParamKeyTable for auth module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		SprintDuration: DefaultSprintDuration,
		SpanDuration:   DefaultSpanDuration,
		ProducerCount:  DefaultProducerCount,
	}
}

func validateSprintDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid sprint duration: %d", v)
	}

	return nil
}

func validateSpanDuration(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid span duration: %d", v)
	}

	return nil
}

func validateProducerCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid producers count: %d", v)
	}

	return nil
}
