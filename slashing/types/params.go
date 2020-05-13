package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/params/subspace"
)

// Default parameter namespace
const (
	DefaultParamspace           = ModuleName
	DefaultSignedBlocksWindow   = int64(100)
	DefaultDowntimeJailDuration = 60 * 10 * time.Second
)

var (
	DefaultMinSignedPerWindow      = sdk.NewDecWithPrec(5, 1)
	DefaultSlashFractionDoubleSign = sdk.NewDec(1).Quo(sdk.NewDec(20))
	DefaultSlashFractionDowntime   = sdk.NewDec(1).Quo(sdk.NewDec(100))
	DefaultSlashFractionLimit      = sdk.NewDec(1).Quo(sdk.NewDec(3))
	DefaultJailFractionLimit       = sdk.NewDec(1).Quo(sdk.NewDec(3))
	DefaultMaxEvidenceAge          = 60 * 2 * time.Second
	DefaultEnableSlashing          = false
)

// Parameter store keys
var (
	KeySignedBlocksWindow      = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow      = []byte("MinSignedPerWindow")
	KeyDowntimeJailDuration    = []byte("DowntimeJailDuration")
	KeySlashFractionDoubleSign = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime   = []byte("SlashFractionDowntime")
	KeySlashFractionLimit      = []byte("SlashFractionLimit")
	KeyJailFractionLimit       = []byte("JailFractionLimit")
	KeyMaxEvidenceAge          = []byte("MaxEvidenceAge")
	KeyEnableSlashing          = []byte("EnableSlashing")
)

var _ subspace.ParamSet = &Params{}

// Params - used for initializing default parameter for slashing at genesis
type Params struct {
	SignedBlocksWindow      int64         `json:"signed_blocks_window" yaml:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window" yaml:"min_signed_per_window"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration" yaml:"downtime_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign" yaml:"slash_fraction_double_sign"` // fraction amount to slash on double sign
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime" yaml:"slash_fraction_downtime"`       // fraction amount to slash on downtime
	SlashFractionLimit      sdk.Dec       `json:"slash_fraction_limit" yaml:"slash_fraction_limit"`             // if totalSlashedAmount crossed SlashFraction of totalValidatorPower, emit Slash-limit event
	JailFractionLimit       sdk.Dec       `json:"jail_fraction_limit" yaml:"jail_fraction_limit"`               // if slashedAmount crossed JailFraction of validatorPower, Jail him
	MaxEvidenceAge          time.Duration `json:"max_evidence_age" yaml:"max_evidence_age"`
	EnableSlashing          bool          `json:"enable_slashing" yaml:"enable_slashing"`
}

// NewParams creates a new Params object
func NewParams(
	signedBlocksWindow int64, minSignedPerWindow sdk.Dec, downtimeJailDuration time.Duration,
	slashFractionDoubleSign, slashFractionDowntime sdk.Dec, slashFractionLimit sdk.Dec, jailFractionLimit sdk.Dec, maxEvidenceAge time.Duration, enableSlashing bool,
) Params {

	return Params{
		SignedBlocksWindow:      signedBlocksWindow,
		MinSignedPerWindow:      minSignedPerWindow,
		DowntimeJailDuration:    downtimeJailDuration,
		SlashFractionDoubleSign: slashFractionDoubleSign,
		SlashFractionDowntime:   slashFractionDowntime,
		MaxEvidenceAge:          maxEvidenceAge,
		SlashFractionLimit:      slashFractionLimit,
		JailFractionLimit:       jailFractionLimit,
		EnableSlashing:          enableSlashing,
	}
}

// ParamKeyTable for slashing module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Slashing Params:
  SignedBlocksWindow:      %d
  MinSignedPerWindow:      %s
  DowntimeJailDuration:    %s
  SlashFractionDoubleSign: %s
  MaxEvidenceAge: %s
  SlashFractionDowntime:   %s
  SlashFractionLimit:   %s
  JailFractionDowntime:   %s
  EnableSlashing:   %s`,
		p.SignedBlocksWindow, p.MinSignedPerWindow,
		p.DowntimeJailDuration, p.SlashFractionDoubleSign, p.MaxEvidenceAge,
		p.SlashFractionDowntime, p.SlashFractionLimit, p.JailFractionLimit, p.EnableSlashing)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{KeySignedBlocksWindow, &p.SignedBlocksWindow},
		{KeyMinSignedPerWindow, &p.MinSignedPerWindow},
		{KeyDowntimeJailDuration, &p.DowntimeJailDuration},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
		{KeySlashFractionLimit, &p.SlashFractionLimit},
		{KeyJailFractionLimit, &p.JailFractionLimit},
		{KeyMaxEvidenceAge, &p.MaxEvidenceAge},
		{KeyEnableSlashing, &p.EnableSlashing},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultSignedBlocksWindow, DefaultMinSignedPerWindow, DefaultDowntimeJailDuration,
		DefaultSlashFractionDoubleSign, DefaultSlashFractionDowntime, DefaultSlashFractionLimit, DefaultJailFractionLimit, DefaultMaxEvidenceAge, DefaultEnableSlashing,
	)
}

func validateSignedBlocksWindow(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("signed blocks window must be positive: %d", v)
	}

	return nil
}

func validateMinSignedPerWindow(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min signed per window cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min signed per window too large: %s", v)
	}

	return nil
}

func validateDowntimeJailDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("downtime jail duration must be positive: %s", v)
	}

	return nil
}

func validateSlashFractionDoubleSign(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("double sign slash fraction cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("double sign slash fraction too large: %s", v)
	}

	return nil
}

func validateSlashFractionDowntime(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("downtime slash fraction cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("downtime slash fraction too large: %s", v)
	}

	return nil
}
