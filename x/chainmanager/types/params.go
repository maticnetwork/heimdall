package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/helper"
	// "github.com/maticnetwork/heimdall/params/paramtypes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	borCommon "github.com/maticnetwork/bor/common"
)

// Default parameter values
const (
	DefaultMainchainTxConfirmations  uint64 = 6
	DefaultMaticchainTxConfirmations uint64 = 10
)

var (
	// DefaultStateReceiverAddress is used set Default State Receiver address
	DefaultStateReceiverAddress sdk.AccAddress = sdk.AccAddress(borCommon.FromHex("0x0000000000000000000000000000000000001001"))
	// DefaultValidatorSetAddress is used set Default Validator Set address
	DefaultValidatorSetAddress sdk.AccAddress = sdk.AccAddress(borCommon.FromHex("0x0000000000000000000000000000000000001000"))
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
	// TODO fix ParamSetPair issue, enabling this throwing `assignment to entry in nil map`
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMainchainTxConfirmations, &p.MainchainTxConfirmations, validateMainchainTxConfirmations),
		paramtypes.NewParamSetPair(KeyMaticchainTxConfirmations, &p.MaticchainTxConfirmations, validateMaticchainTxConfirmations),
		paramtypes.NewParamSetPair(KeyChainParams, &p.ChainParams, validateChainParams),
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	// TODO add ProtoCodec instead of AminoCodec
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
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
	addr, err := sdk.AccAddressFromHex(p.ChainParams.MaticTokenAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(MaticTokenAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.StakingManagerAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(StakingManagerAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.SlashManagerAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(SlashManagerAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.RootChainAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(RootChainAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.StakingInfoAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(StakingInfoAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.StateSenderAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(StateSenderAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.StateReceiverAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(StateReceiverAddress, addr); err != nil {
		return err
	}

	addr, err = sdk.AccAddressFromHex(p.ChainParams.ValidatorSetAddress)
	if err != nil {
		return err
	}
	if err := validateAccAddress(ValidatorSetAddress, addr); err != nil {
		return err
	}

	return nil
}

func validateAccAddress(key string, value sdk.AccAddress) error {
	if value.String() == "" {
		return fmt.Errorf("Invalid value %s in chain_params", key)
	}

	// TODO add validation based on Key and Address

	return nil
}

func validateMainchainTxConfirmations(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("Mainchain Tx Confirmations must be positive: %d", v)
	}

	return nil
}

func validateMaticchainTxConfirmations(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("Maticchain Tx Confirmations must be positive: %d", v)
	}

	return nil
}

func validateChainParams(i interface{}) error {
	_, ok := i.(*ChainParams)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
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
			StateReceiverAddress: DefaultStateReceiverAddress.String(),
			ValidatorSetAddress:  DefaultValidatorSetAddress.String(),
		},
	}
}
