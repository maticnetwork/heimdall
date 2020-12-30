package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultCheckpointBufferTime time.Duration = 1000 * time.Second // Time checkpoint is allowed to stay in buffer (1000 seconds ~ 17 mins)
	DefaultAvgCheckpointLength  uint64        = 256
	DefaultMaxCheckpointLength  uint64        = 1024
	DefaultChildBlockInterval   uint64        = 10000
)

// Parameter keys
var (
	KeyCheckpointBufferTime = []byte("CheckpointBufferTime")
	KeyAvgCheckpointLength  = []byte("AvgCheckpointLength")
	KeyMaxCheckpointLength  = []byte("MaxCheckpointLength")
	KeyChildBlockInterval   = []byte("ChildBlockInterval")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// NewParams creates a new Params object
func NewParams(
	checkpointBufferTime time.Duration,
	checkpointLength uint64,
	maxCheckpointLength uint64,
	childBlockInterval uint64,
) Params {
	return Params{
		CheckpointBufferTime: checkpointBufferTime,
		AvgCheckpointLength:  checkpointLength,
		MaxCheckpointLength:  maxCheckpointLength,
		ChildBlockInterval:   childBlockInterval,
	}
}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// TODO uncomment these lines
		// paramtypes.NewParamSetPair(KeyCheckpointBufferTime, &p.CheckpointBufferTime, validateCheckpointBufferTime),
		// paramtypes.NewParamSetPair(KeyAvgCheckpointLength, &p.AvgCheckpointLength, validateAvgCheckpointLength),
		// paramtypes.NewParamSetPair(KeyMaxCheckpointLength, &p.MaxCheckpointLength, validateMaxCheckpointLength),
		// paramtypes.NewParamSetPair(KeyChildBlockInterval, &p.ChildBlockInterval, validateChildBlockInterval),
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
		CheckpointBufferTime: DefaultCheckpointBufferTime,
		AvgCheckpointLength:  DefaultAvgCheckpointLength,
		MaxCheckpointLength:  DefaultMaxCheckpointLength,
		ChildBlockInterval:   DefaultChildBlockInterval,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CheckpointBufferTime: %s\n", p.CheckpointBufferTime))
	sb.WriteString(fmt.Sprintf("AvgCheckpointLength: %d\n", p.AvgCheckpointLength))
	sb.WriteString(fmt.Sprintf("MaxCheckpointLength: %d\n", p.MaxCheckpointLength))
	sb.WriteString(fmt.Sprintf("ChildBlockInterval: %d\n", p.ChildBlockInterval))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if p.MaxCheckpointLength == 0 || p.AvgCheckpointLength == 0 {
		return fmt.Errorf("MaxCheckpointLength, AvgCheckpointLength should be non-zero")
	}

	if p.MaxCheckpointLength < p.AvgCheckpointLength {
		return fmt.Errorf("AvgCheckpointLength should not be greater than MaxCheckpointLength")
	}

	if p.ChildBlockInterval == 0 {
		return fmt.Errorf("ChildBlockInterval should be greater than zero")
	}

	return nil
}

// TODO add validations
func validateCheckpointBufferTime(i interface{}) error {
	_, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// TODO add validations
func validateAvgCheckpointLength(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// TODO add validations
func validateMaxCheckpointLength(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

// TODO add validations
func validateChildBlockInterval(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
