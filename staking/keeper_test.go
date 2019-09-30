package staking_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/staking"
	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// tests setter/getters for validatorSignerMaps , validator set/get
func TestValidator(t *testing.T) {
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	validators := cmn.GenRandomVal(1, 0, 10, uint64(10), true, 1)
	// Add Validator to state
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}
	// Get Validator Info from state
	valInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator", err)
	}
	// Get Signer Address mapped with ValidatorId
	mappedSignerAddress, isMapped := keeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Signer Address not mapped to Validator Id")
	}
	// Check if Validator matches in state
	require.Equal(t, valInfo, validators[0], "Validators in state doesnt match")
	require.Equal(t, types.HexToHeimdallAddress(mappedSignerAddress.Hex()), validators[0].Signer, "Signer address doesn't match")
}

// tests power change, validator creation, validator set update when signer changes
func TestUpdateSigner(t *testing.T) {
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	validators := cmn.GenRandomVal(1, 0, 10, uint64(10), true, 1)

	// Add validator to State
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	// Fetch Validator Info from Store
	valInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info from store", err)
	}

	// Update Signer
	newPrivKey := secp256k1.GenPrivKey()
	newPubKey := types.NewPubKey(newPrivKey.PubKey().Bytes())
	newSigner := types.HexToHeimdallAddress(newPubKey.Address().String())
	err = keeper.UpdateSigner(ctx, newSigner, newPubKey, valInfo.Signer)
	if err != nil {
		t.Error("Error while updating Signer Address -", err)
	}

	// Check Validator Info of Prev Signer
	prevSginerValInfo, err := keeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info for Prev Signer - ", err)
	}
	require.Equal(t, uint64(0), prevSginerValInfo.Power, "Power of Prev Signer should be zero")

	// Check Validator Info of Updated Signer
	updatedSignerValInfo, err := keeper.GetValidatorInfo(ctx, newSigner.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info for Updater Signer", err)
	}
	require.Equal(t, validators[0].Power, updatedSignerValInfo.Power, "power of updated signer should match with prev signer power")

	// Check If ValidatorId is mapped To Updated Signer
	signerAddress, isMapped := keeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Validator Id is not mapped to Signer Address", err)
	}
	require.Equal(t, newSigner, types.HexToHeimdallAddress(signerAddress.Hex()), "Validator ID should be mapped to Updated Signer Address")

	// Check total Validators
	totalValidators := keeper.GetAllValidators(ctx)
	require.Equal(t, 2, len(totalValidators), "Total Validators should be two.")
}

func TestCurrentValidator(t *testing.T) {
	type TestDataItem struct {
		name       string
		startblock uint64
		endblock   uint64
		power      uint64
		ackcount   uint64
		result     bool
		resultmsg  string
	}

	dataItems := []TestDataItem{
		{"power zero", uint64(0), uint64(0), uint64(0), uint64(1), false, "should not be current validator as power is zero."},
		{"start epoch greater than ackcount", uint64(3), uint64(0), uint64(10), uint64(1), false, "should not be current validator as start epoch greater than ackcount."},
	}

	ctx, stakingKeeper, checkpointKeeper := cmn.CreateTestInput(t, false)
	for i, item := range dataItems {
		t.Run(item.name, func(t *testing.T) {
			// Create a Validator [startEpoch, endEpoch, Power]
			privKep := secp256k1.GenPrivKey()
			pubkey := types.NewPubKey(privKep.PubKey().Bytes())
			newVal := types.Validator{
				ID:         types.NewValidatorID(1 + uint64(i)),
				StartEpoch: item.startblock,
				EndEpoch:   item.startblock,
				Power:      item.power,
				Signer:     types.HexToHeimdallAddress(pubkey.Address().String()),
				PubKey:     pubkey,
				Accum:      0,
			}
			// check current validator
			stakingKeeper.AddValidator(ctx, newVal)
			checkpointKeeper.UpdateACKCountWithValue(ctx, item.ackcount)
			t.Log("Ack count - ", checkpointKeeper.GetACKCount(ctx))
			isCurrentVal := stakingKeeper.IsCurrentValidatorByAddress(ctx, newVal.Signer.Bytes())
			require.Equal(t, item.result, isCurrentVal, item.resultmsg)
		})
	}
}

func TestValidatorSet(t *testing.T) {
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	valSet := LoadValidatorSet(4, t, keeper, ctx, true, 10)

	storedValSet := keeper.GetValidatorSet(ctx)
	require.Equal(t, valSet, storedValSet, "Validator Set in state doesnt match ")
}

func LoadValidatorSet(count int, t *testing.T, keeper staking.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	// create 4 validators
	validators := cmn.GenRandomVal(4, 0, 10, uint64(timeAlive), randomise, 1)
	var valSet types.ValidatorSet

	// add validators to new Validator set and state
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.Empty(t, err, "Unable to set validator, Error: %v", err)
		// add validator to validator set
		valSet.Add(&validator)
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")
	vals := keeper.GetAllValidators(ctx)
	t.Log("Vals inserted", vals)
	return valSet
}
