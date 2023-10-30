package types

import (
	"sync/atomic"
)

var milestoneId atomic.Value

func GetMilestoneID() string {
	if milestoneId.Load() == nil {
		return ""
	}

	return milestoneId.Load().(string)
}

func SetMilestoneID(id string) {
	milestoneId.Store(id)
}
