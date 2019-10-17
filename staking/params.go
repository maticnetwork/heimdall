package staking

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultCheckpointReward - Total Checkpoint reward
	// DefaultCheckpointReward = big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil).Bytes()

	// DefaultProposerBonusPercent - Proposer Signer Reward Ratio
	DefaultProposerBonusPercent = int64(10)
)

// ParamStoreKeyCheckpointReward - Store's Key for Reward amount
var ParamStoreKeyCheckpointReward = []byte("checkpointreward")

// ParamStoreKeyProposerBonusPercent - Store's Key for Reward amount
var ParamStoreKeyProposerBonusPercent = []byte("proposerbonuspercent")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	DefaultCheckpointReward := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil).Bytes()
	return params.NewKeyTable(
		ParamStoreKeyCheckpointReward, DefaultCheckpointReward,
		ParamStoreKeyProposerBonusPercent, DefaultProposerBonusPercent,
	)
}
