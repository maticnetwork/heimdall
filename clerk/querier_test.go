package clerk_test

import (
	"testing"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk"
	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

// TestStateSyncerQueryHandler will test state syncer selection logic
func TestStateSyncerQueryHandler(t *testing.T) {

	type TestDataItem struct {
		name               string
		validatorCount     int
		stateEventCount    uint64
		expectedSyncerSize int
		expectedValIDs     []types.ValidatorID
	}

	dataItems := []TestDataItem{
		{"validators less than 3", int(2), uint64(0), int(2), []types.ValidatorID{types.NewValidatorID(1), types.NewValidatorID(2)}},
		{"validators less than 3", int(2), uint64(1), int(2), []types.ValidatorID{types.NewValidatorID(2), types.NewValidatorID(1)}},
		{"validators  equal 3", int(3), uint64(0), int(3), []types.ValidatorID{types.NewValidatorID(1), types.NewValidatorID(2), types.NewValidatorID(3)}},
		{"validators  greater than 3", int(4), uint64(0), int(3), []types.ValidatorID{types.NewValidatorID(1), types.NewValidatorID(2), types.NewValidatorID(3)}},
		{"validators=5 stateEventCount=0", int(5), uint64(0), int(3), []types.ValidatorID{types.NewValidatorID(1), types.NewValidatorID(2), types.NewValidatorID(3)}},
		{"validators=5 stateEventCount=1", int(5), uint64(1), int(3), []types.ValidatorID{types.NewValidatorID(2), types.NewValidatorID(3), types.NewValidatorID(4)}},
		{"validators=5 stateEventCount=2", int(5), uint64(2), int(3), []types.ValidatorID{types.NewValidatorID(3), types.NewValidatorID(4), types.NewValidatorID(5)}},
		{"validators=5 stateEventCount=3", int(5), uint64(3), int(3), []types.ValidatorID{types.NewValidatorID(4), types.NewValidatorID(5), types.NewValidatorID(1)}},
		{"validators=5 stateEventCount=4", int(5), uint64(4), int(3), []types.ValidatorID{types.NewValidatorID(5), types.NewValidatorID(1), types.NewValidatorID(2)}},
		{"validators=5 stateEventCount=5", int(5), uint64(5), int(3), []types.ValidatorID{types.NewValidatorID(1), types.NewValidatorID(2), types.NewValidatorID(3)}},
		{"validators=5 stateEventCount=6", int(5), uint64(6), int(3), []types.ValidatorID{types.NewValidatorID(2), types.NewValidatorID(3), types.NewValidatorID(4)}},
	}

	for _, item := range dataItems {
		t.Run(item.name, func(t *testing.T) {
			ctx, stakingKeeper, _, clerkKeeper := cmn.CreateTestInput(t, false)
			cmn.LoadValidatorSet(item.validatorCount, t, stakingKeeper, ctx, false, 10)
			clerkKeeper.SetStateSyncEventCount(ctx, item.stateEventCount)

			res, err := clerk.HandlerQueryStateSyncer(ctx, abci.RequestQuery{}, clerkKeeper, stakingKeeper)
			require.Empty(t, err, "Error should be empty")

			cdc := app.MakeCodec()
			var stateSyncerList []types.Validator
			abcierr := cdc.UnmarshalJSON(res, &stateSyncerList)
			require.Empty(t, abcierr, "Error while unmarshalling validator")

			require.Equal(t, item.expectedSyncerSize, len(stateSyncerList), "stateSyncerList size should be %v but it is %v", item.expectedSyncerSize, len(stateSyncerList))
			for i, val := range stateSyncerList {
				require.Equal(t, item.expectedValIDs[i], val.ID, "validator ID mismatch")
			}
		})
	}

}
