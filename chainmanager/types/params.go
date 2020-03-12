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
	DefaultTxConfirmationTime time.Duration = 6 * 14 * time.Second
)

// Parameter keys
var (
	KeyTxConfirmationTime = []byte("TxConfirmationTime")
)

var _ subspace.ParamSet = &Params{}

// Params defines the parameters for the auth module.
type Params struct {
	TxConfirmationTime time.Duration `json:"tx_confirmation_time" yaml:"tx_confirmation_time"` // tx confirmation duration
}

// NewParams creates a new Params object
func NewParams(txConfirmationTime time.Duration) Params {
	return Params{
		TxConfirmationTime: txConfirmationTime,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyTxConfirmationTime, &p.TxConfirmationTime},
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
	sb.WriteString(fmt.Sprintf("TxConfirmationTime: %d\n", p.TxConfirmationTime))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
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
		TxConfirmationTime: DefaultTxConfirmationTime,
	}
}
