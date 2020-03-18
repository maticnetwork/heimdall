package types

import (
	"github.com/maticnetwork/heimdall/params/subspace"
)

const (

	// DefaultProposerBonusPercent - Proposer Signer Reward Ratio
	DefaultProposerBonusPercent = int64(10)
)

// ParamStoreKeyProposerBonusPercent - Store's Key for Reward amount
var ParamStoreKeyProposerBonusPercent = []byte("proposerbonuspercent")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable(
		ParamStoreKeyProposerBonusPercent, DefaultProposerBonusPercent,
	)
}
