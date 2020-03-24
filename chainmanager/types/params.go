package types

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/maticnetwork/heimdall/params/subspace"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// Default parameter values
const (
	DefaultTxConfirmationTime time.Duration = 6 * 14 * time.Second
)

// Parameter keys
var (
	KeyTxConfirmationTime = []byte("TxConfirmationTime")
	KeyChainParams        = []byte("ChainParams")
)

var _ subspace.ParamSet = &Params{}

// ChainParams chain related params
type ChainParams struct {
	MaticTokenAddress     hmTypes.HeimdallAddress `json:"matic_token_address" yaml:"matic_token_address"`
	StakingManagerAddress hmTypes.HeimdallAddress `json:"staking_manager_address" yaml:"staking_manager_address"`
	RootChainAddress      hmTypes.HeimdallAddress `json:"root_chain_address" yaml:"root_chain_address"`
	StakingInfoAddress    hmTypes.HeimdallAddress `json:"staking_info_address" yaml:"staking_info_address"`
	StateSenderAddress    hmTypes.HeimdallAddress `json:"state_sender_address" yaml:"state_sender_address"`
}

func (cp ChainParams) String() string {
	return fmt.Sprintf(`
  MaticTokenAddress:            %s
	StakingManagerAddress:        %s
	RootChainAddress:             %s
  StakingInfoAddress:           %s
  StateSenderAddress:           %s`,
		cp.MaticTokenAddress, cp.StakingManagerAddress, cp.RootChainAddress, cp.StakingInfoAddress, cp.StateSenderAddress)
}

// Params defines the parameters for the auth module.
type Params struct {
	TxConfirmationTime time.Duration          `json:"tx_confirmation_time" yaml:"tx_confirmation_time"` // tx confirmation duration
	ChainParams        map[string]ChainParams `json:"chain_params" yaml:"chain_params"`
}

// NewParams creates a new Params object
func NewParams(txConfirmationTime time.Duration, chainParams map[string]ChainParams) Params {
	return Params{
		TxConfirmationTime: txConfirmationTime,
		ChainParams:        chainParams,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeyTxConfirmationTime, &p.TxConfirmationTime},
		{KeyChainParams, &p.ChainParams},
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
	sb.WriteString("ChainParams: \n")
	for key, val := range p.ChainParams {
		sb.WriteString(fmt.Sprintf(" %s:\n", key))
		sb.WriteString(fmt.Sprintf("     %s:\n", val.String()))
	}
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
