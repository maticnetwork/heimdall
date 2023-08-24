package types

import (
	"sync/atomic"
)

var milestoneId atomic.Value

func GetMilestoneID() string {
	return milestoneId.Load().(string)
}

func SetMilestoneID(id string) {
	milestoneId.Store(id)
}
