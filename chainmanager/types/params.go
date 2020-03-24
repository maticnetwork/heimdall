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
	BorChainID            string                  `json:"bor_chain_id" yaml:"bor_chain_id"`
	MaticTokenAddress     hmTypes.HeimdallAddress `json:"matic_token_address" yaml:"matic_token_address"`
	StakingManagerAddress hmTypes.HeimdallAddress `json:"staking_manager_address" yaml:"staking_manager_address"`
	RootChainAddress      hmTypes.HeimdallAddress `json:"root_chain_address" yaml:"root_chain_address"`
	StakingInfoAddress    hmTypes.HeimdallAddress `json:"staking_info_address" yaml:"staking_info_address"`
	StateSenderAddress    hmTypes.HeimdallAddress `json:"state_sender_address" yaml:"state_sender_address"`
}

func (cp ChainParams) String() string {
	return fmt.Sprintf(`
	BorChainID: 									%s
  MaticTokenAddress:            %s
	StakingManagerAddress:        %s
	RootChainAddress:             %s
  StakingInfoAddress:           %s
  StateSenderAddress:           %s`,
		cp.BorChainID, cp.MaticTokenAddress, cp.StakingManagerAddress, cp.RootChainAddress, cp.StakingInfoAddress, cp.StateSenderAddress)
}

// Params defines the parameters for the auth module.
type Params struct {
	TxConfirmationTime time.Duration `json:"tx_confirmation_time" yaml:"tx_confirmation_time"` // tx confirmation duration
	ChainParams        ChainParams   `json:"chain_params" yaml:"chain_params"`
}

// NewParams creates a new Params object
func NewParams(txConfirmationTime time.Duration, chainParams ChainParams) Params {
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
	sb.WriteString(fmt.Sprintf("ChainParams: %s\n", p.ChainParams.String()))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateHeimdallAddress(p.ChainParams.MaticTokenAddress, "matic_token_address"); err != nil {
		return err
	}

	if err := validateHeimdallAddress(p.ChainParams.StakingManagerAddress, "staking_manager_address"); err != nil {
		return err
	}

	if err := validateHeimdallAddress(p.ChainParams.RootChainAddress, "root_chain_address"); err != nil {
		return err
	}

	if err := validateHeimdallAddress(p.ChainParams.StakingInfoAddress, "staking_info_address"); err != nil {
		return err
	}

	if err := validateHeimdallAddress(p.ChainParams.StateSenderAddress, "state_sender_address"); err != nil {
		return err
	}

	return nil
}

func validateHeimdallAddress(value hmTypes.HeimdallAddress, key string) error {
	if value.String() == "" {
		return fmt.Errorf("Invalid value %s in chain_params", key)
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
		TxConfirmationTime: DefaultTxConfirmationTime,
	}
}
