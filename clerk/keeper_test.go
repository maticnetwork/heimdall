package clerk_test

import (
	"testing"

	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/stretchr/testify/require"
)

// TestStateSyncEventCount tests setter, getter, and increment of statesynceventcount
func TestStateSyncEventCount(t *testing.T) {
	ctx, _, _, clerkKeeper := cmn.CreateTestInput(t, false)

	// Test genesis State sync event count
	initialCount := clerkKeeper.GetStateSyncEventCount(ctx)
	require.Equal(t, uint64(0), initialCount, "Initial Count of Events should be %v but it is %v", 0, initialCount)

	// Test setter
	clerkKeeper.SetStateSyncEventCount(ctx, uint64(2))
	updatedCount := clerkKeeper.GetStateSyncEventCount(ctx)
	require.Equal(t, uint64(2), updatedCount, "updated Count of Events should be %v but it is %v", 2, updatedCount)

	// Increment
	clerkKeeper.IncrementStateSyncEventCount(ctx)
	countAfterIncrement := clerkKeeper.GetStateSyncEventCount(ctx)
	t.Log(countAfterIncrement)
	require.Equal(t, uint64(3), countAfterIncrement, "Incremented Count of Events should be %v but it is %v", 3, countAfterIncrement)
}
