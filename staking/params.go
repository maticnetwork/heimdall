package staking

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (

	// DefaultProposerBonusPercent - Proposer Signer Reward Ratio
	DefaultProposerBonusPercent = int64(10)
)

// ParamStoreKeyProposerBonusPercent - Store's Key for Reward amount
var ParamStoreKeyProposerBonusPercent = []byte("proposerbonuspercent")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyProposerBonusPercent, DefaultProposerBonusPercent,
	)
}
