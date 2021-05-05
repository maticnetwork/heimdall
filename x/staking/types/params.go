package types

import (
	paramtypes "github.com/maticnetwork/heimdall/x/params/types"
)

const (

	// DefaultProposerBonusPercent - Proposer Signer Reward Ratio
	DefaultProposerBonusPercent = int64(10)
)

// ParamStoreKeyProposerBonusPercent - Store's Key for Reward amount
var ParamStoreKeyProposerBonusPercent = []byte("proposerbonuspercent")

var KeyBondDenom = []byte("BondDenom")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyProposerBonusPercent, DefaultProposerBonusPercent, validateProposerBonusPercent),
	// ParamStoreKeyProposerBonusPercent, DefaultProposerBonusPercent,
	)
}

func validateProposerBonusPercent(i interface{}) error {
	// v, ok := i.(sdk.Coin)
	// if !ok {
	// 	return fmt.Errorf("invalid parameter type: %T", i)
	// }

	// if !v.IsValid() {
	// 	return fmt.Errorf("invalid constant fee: %s", v)
	// }

	return nil
}
