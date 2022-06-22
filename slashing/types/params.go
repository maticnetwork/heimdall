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
  EnableSlashing:   %t`,
		p.SignedBlocksWindow, p.MinSignedPerWindow,
		p.DowntimeJailDuration, p.SlashFractionDoubleSign, p.MaxEvidenceAge,
		p.SlashFractionDowntime, p.SlashFractionLimit, p.JailFractionLimit, p.EnableSlashing)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{Key: KeySignedBlocksWindow, Value: &p.SignedBlocksWindow},
		{Key: KeyMinSignedPerWindow, Value: &p.MinSignedPerWindow},
		{Key: KeyDowntimeJailDuration, Value: &p.DowntimeJailDuration},
		{Key: KeySlashFractionDoubleSign, Value: &p.SlashFractionDoubleSign},
		{Key: KeySlashFractionDowntime, Value: &p.SlashFractionDowntime},
		{Key: KeySlashFractionLimit, Value: &p.SlashFractionLimit},
		{Key: KeyJailFractionLimit, Value: &p.JailFractionLimit},
		{Key: KeyMaxEvidenceAge, Value: &p.MaxEvidenceAge},
		{Key: KeyEnableSlashing, Value: &p.EnableSlashing},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(
		DefaultSignedBlocksWindow, DefaultMinSignedPerWindow, DefaultDowntimeJailDuration,
		DefaultSlashFractionDoubleSign, DefaultSlashFractionDowntime, DefaultSlashFractionLimit, DefaultJailFractionLimit, DefaultMaxEvidenceAge, DefaultEnableSlashing,
	)
}
