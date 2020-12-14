package types

import (
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/helper"
	// "github.com/maticnetwork/heimdall/params/paramtypes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultMainchainTxConfirmations  uint64 = 6
	DefaultMaticchainTxConfirmations uint64 = 10
)

var (
	DefaultStateReceiverAddress sdk.AccAddress = sdk.AccAddress("0x0000000000000000000000000000000000001001")
	DefaultValidatorSetAddress  sdk.AccAddress = sdk.AccAddress("0x0000000000000000000000000000000000001000")
)

// Parameter keys
var (
	KeyMainchainTxConfirmations  = []byte("MainchainTxConfirmations")
	KeyMaticchainTxConfirmations = []byte("MaticchainTxConfirmations")
	KeyChainParams               = []byte("ChainParams")
)

var _ paramtypes.ParamSet = &Params{}

func (cp ChainParams) String() string {
	return fmt.Sprintf(`
	BorChainID: 									%s
  MaticTokenAddress:            %s
	StakingManagerAddress:        %s
	SlashManagerAddress:        %s
	RootChainAddress:             %s
  StakingInfoAddress:           %s
	StateSenderAddress:           %s
	StateReceiverAddress: 				%s
	ValidatorSetAddress:					%s`,
		cp.BorChainID, cp.MaticTokenAddress, cp.StakingManagerAddress, cp.SlashManagerAddress, cp.RootChainAddress, cp.StakingInfoAddress, cp.StateSenderAddress, cp.StateReceiverAddress, cp.ValidatorSetAddress)
}

// NewParams creates a new Params object
func NewParams(
	mainchainTxConfirmations uint64,
	maticchainTxConfirmations uint64,
	chainParams *ChainParams,
) Params {
	return Params{
		MainchainTxConfirmations:  mainchainTxConfirmations,
		MaticchainTxConfirmations: maticchainTxConfirmations,
		ChainParams:               chainParams,
	}
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		// {KeyMainchainTxConfirmations, &p.MainchainTxConfirmations},
		// {KeyMaticchainTxConfirmations, &p.MaticchainTxConfirmations},
		// {KeyChainParams, &p.ChainParams},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	// bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	// bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	// return bytes.Equal(bz1, bz2)
	return true
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("MainchainTxConfirmations: %d\n", p.MainchainTxConfirmations))
	sb.WriteString(fmt.Sprintf("MaticchainTxConfirmations: %d\n", p.MaticchainTxConfirmations))
	sb.WriteString(fmt.Sprintf("ChainParams: %s\n", p.ChainParams.String()))
	return sb.String()
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateAddress("matic_token_address", p.ChainParams.MaticTokenAddress); err != nil {
		return err
	}

	if err := validateAddress("staking_manager_address", p.ChainParams.StakingManagerAddress); err != nil {
		return err
	}

	if err := validateAddress("slash_manager_address", p.ChainParams.SlashManagerAddress); err != nil {
		return err
	}

	if err := validateAddress("root_chain_address", p.ChainParams.RootChainAddress); err != nil {
		return err
	}

	if err := validateAddress("staking_info_address", p.ChainParams.StakingInfoAddress); err != nil {
		return err
	}

	if err := validateAddress("state_sender_address", p.ChainParams.StateSenderAddress); err != nil {
		return err
	}

	if err := validateAddress("state_receiver_address", p.ChainParams.StateReceiverAddress); err != nil {
		return err
	}

	if err := validateAddress("validator_set_address", p.ChainParams.ValidatorSetAddress); err != nil {
		return err
	}

	return nil
}

func validateAddress(key string, value sdk.AccAddress) error {
	if value.String() == "" {
		return fmt.Errorf("Invalid value %s in chain_params", key)
	}

	return nil
}

//
// Extra functions
//

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() *Params {
	return &Params{
		MainchainTxConfirmations:  DefaultMainchainTxConfirmations,
		MaticchainTxConfirmations: DefaultMaticchainTxConfirmations,
		ChainParams: &ChainParams{
			BorChainID:           helper.DefaultBorChainID,
			StateReceiverAddress: DefaultStateReceiverAddress,
			ValidatorSetAddress:  DefaultValidatorSetAddress,
		},
	}
}
