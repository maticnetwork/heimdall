package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/maticnetwork/heimdall/params/subspace"
)

// Default parameter values
const (
	DefaultCheckpointBufferTime time.Duration = 1000 * time.Second
)

// Parameter keys
var (
	KeyCheckpointBufferTime = []byte("CheckpointBufferTime")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	CheckpointBufferTime time.Duration `json:"checkpoint_buffer_time" yaml:"checkpoint_buffer_time"`
}

// NewParams creates a new Params object
func NewParams(
	checkpointBufferTime time.Duration,
) Params {
	return Params{
		CheckpointBufferTime: checkpointBufferTime,
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
		{KeyCheckpointBufferTime, &p.CheckpointBufferTime},
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
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("CheckpointBufferTime: %s\n", p.CheckpointBufferTime))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	return nil
}
