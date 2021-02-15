package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultSprintDuration    uint64 = 64
	DefaultSpanDuration             = 100 * DefaultSprintDuration
	DefaultFirstSpanDuration uint64 = 256
	DefaultProducerCount     uint64 = 4
)

// Parameter keys
var (
	KeySprintDuration = []byte("SprintDuration")
	KeySpanDuration   = []byte("SpanDuration")
	KeyProducerCount  = []byte("ProducerCount")
)

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		SprintDuration: DefaultSprintDuration,
		SpanDuration:   DefaultSpanDuration,
		ProducerCount:  DefaultProducerCount,
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeySprintDuration, &p.SprintDuration, validateSprintDuration),
		paramtypes.NewParamSetPair(KeySpanDuration, &p.SpanDuration, validateSpanDuration),
		paramtypes.NewParamSetPair(KeyProducerCount, &p.ProducerCount, validateProducerCount),
	}
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
