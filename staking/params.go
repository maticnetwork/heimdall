package staking

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultCheckpointReward - Total Checkpoint reward
	DefaultCheckpointReward = uint64(10000)

	// DefaultProposerToSignerRewards - Proposer Signer Reward Ratio
	DefaultProposerToSignerRewards = uint64(10)
)

// ParamStoreKeyCheckpointReward - Store's Key for Reward amount
var ParamStoreKeyCheckpointReward = []byte("checkpointreward")

// ParamStoreKeyProposerToSignerRewards - Store's Key for Reward amount
var ParamStoreKeyProposerToSignerRewards = []byte("proposertosignerrewards")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyCheckpointReward, DefaultCheckpointReward,
		ParamStoreKeyProposerToSignerRewards, DefaultProposerToSignerRewards,
	)
}
