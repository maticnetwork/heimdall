package staking

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultRewardAmount - Reward Amount Given to the Checkpoint Signer
	DefaultRewardAmount = uint64(10)
)

// ParamStoreKeyRewardAmount - Store's Key for Reward amount
var ParamStoreKeyRewardAmount = []byte("rewardamount")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeyRewardAmount, DefaultRewardAmount,
	)
}
