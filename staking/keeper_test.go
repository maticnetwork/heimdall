package staking_test

import (
	"encoding/hex"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	cmn "github.com/maticnetwork/heimdall/test"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"math/big"
	"testing"
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

// tests VotingPower change, validator creation, validator set update when signer changes
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
	require.Equal(t, int64(0), prevSginerValInfo.VotingPower, "VotingPower of Prev Signer should be zero")

	// Check Validator Info of Updated Signer
	updatedSignerValInfo, err := keeper.GetValidatorInfo(ctx, newSigner.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info for Updater Signer", err)
	}
	require.Equal(t, validators[0].VotingPower, updatedSignerValInfo.VotingPower, "VotingPower of updated signer should match with prev signer VotingPower")

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
		name        string
		startblock  uint64
		endblock    uint64
		VotingPower int64
		ackcount    uint64
		result      bool
		resultmsg   string
	}

	dataItems := []TestDataItem{
		{"VotingPower zero", uint64(0), uint64(0), int64(0), uint64(1), false, "should not be current validator as VotingPower is zero."},
		{"start epoch greater than ackcount", uint64(3), uint64(0), int64(10), uint64(1), false, "should not be current validator as start epoch greater than ackcount."},
	}

	ctx, stakingKeeper, checkpointKeeper := cmn.CreateTestInput(t, false)
	for i, item := range dataItems {
		t.Run(item.name, func(t *testing.T) {
			// Create a Validator [startEpoch, endEpoch, VotingPower]
			privKep := secp256k1.GenPrivKey()
			pubkey := types.NewPubKey(privKep.PubKey().Bytes())
			newVal := types.Validator{
				ID:               types.NewValidatorID(1 + uint64(i)),
				StartEpoch:       item.startblock,
				EndEpoch:         item.startblock,
				VotingPower:      item.VotingPower,
				Signer:           types.HexToHeimdallAddress(pubkey.Address().String()),
				PubKey:           pubkey,
				ProposerPriority: 0,
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

func TestValidatorSetChange(t *testing.T) {

	// create sub test to check if validator remove
	t.Run("remove", func(t *testing.T) {
		ctx, keeper, _ := cmn.CreateTestInput(t, false)

		// load 4 validators to state
		cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)

		currentValSet := initValSet.Copy()
		prevValidatorSet := initValSet.Copy()

		// remove validator (making IsCurrentValidator return false)
		prevValidatorSet.Validators[0].StartEpoch = 20

		t.Log("Updated Validators in state")
		for _, v := range prevValidatorSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
		}

		err := keeper.AddValidator(ctx, *prevValidatorSet.Validators[0])
		require.Empty(t, err, "Unable to update validator set")

		// apply updates
		// helper.UpdateValidators(
		// 	currentValSet,                // pointer to current validator set -- UpdateValidators will modify it
		// 	keeper.GetAllValidators(ctx), // All validators
		// 	5,                            // ack count
		// )
		setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
		currentValSet.UpdateWithChangeSet(setUpdates)
		updatedValSet := currentValSet
		t.Log("Validators in updated validator set")
		for _, v := range updatedValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
		}
		// check if 1 validator is removed
		require.Equal(t, len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators), "Validator set should be reduced by one ")
		// remove first validator from initial validator set and equate with new
		t.Log("appended set-", updatedValSet, "one", prevValidatorSet.Validators[1:])
		for _, val := range updatedValSet.Validators {
			if val.Signer == prevValidatorSet.Validators[0].Signer {
				require.Fail(t, "Validator is not removed from updatedvalidator set")
			}
		}
	})

	t.Run("add", func(t *testing.T) {
		ctx, keeper, _ := cmn.CreateTestInput(t, false)

		// load 4 validators to state
		cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)

		validators := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
		prevValSet := initValSet.Copy()
		valToBeAdded := validators[0]
		currentValSet := initValSet.Copy()
		//prevValidatorSet := initValSet.Copy()
		keeper.AddValidator(ctx, valToBeAdded)

		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
		}
		t.Log("Val to be Added")
		t.Log("-->", "Address", valToBeAdded.Signer.String(), "StartEpoch", valToBeAdded.StartEpoch, "EndEpoch", valToBeAdded.EndEpoch, "VotingPower", valToBeAdded.VotingPower)

		// helper.UpdateValidators(
		// 	currentValSet,                // pointer to current validator set -- UpdateValidators will modify it
		// 	keeper.GetAllValidators(ctx), // All validators
		// 	5,                            // ack count
		// )
		setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
		currentValSet.UpdateWithChangeSet(setUpdates)

		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
		}

		require.Equal(t, len(prevValSet.Validators)+1, len(currentValSet.Validators), "Number of validators should be increased by 1")
		require.Equal(t, true, currentValSet.HasAddress(valToBeAdded.Signer.Bytes()), "New Validator should be added")
		require.Equal(t, prevValSet.TotalVotingPower()+int64(valToBeAdded.VotingPower), currentValSet.TotalVotingPower(), "Total VotingPower should be increased")
	})

	t.Run("update", func(t *testing.T) {
		ctx, keeper, _ := cmn.CreateTestInput(t, false)

		// load 4 validators to state
		cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
		initValSet := keeper.GetValidatorSet(ctx)
		t.Log("init val set-", initValSet)
		keeper.IncrementAccum(ctx, 2)
		prevValSet := initValSet.Copy()
		currentValSet := keeper.GetValidatorSet(ctx)
		t.Log("current Val set - ", currentValSet)
		valToUpdate := currentValSet.Validators[0]
		newSigner := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
		t.Log("Validators in old validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "ProposerPriority", v.ProposerPriority, "Signer", v.Signer.String(), "Total VotingPower", currentValSet.TotalVotingPower())
		}
		keeper.UpdateSigner(ctx, newSigner[0].Signer, newSigner[0].PubKey, valToUpdate.Signer)
		// helper.UpdateValidators(
		// 	&currentValSet,               // pointer to current validator set -- UpdateValidators will modify it
		// 	keeper.GetAllValidators(ctx), // All validators
		// 	5,                            // ack count
		// )
		setUpdates := helper.GetUpdatedValidators(&currentValSet, keeper.GetAllValidators(ctx), 5)
		currentValSet.UpdateWithChangeSet(setUpdates)
		t.Log("Validators in updated validator set")
		for _, v := range currentValSet.Validators {
			t.Log("-->", "Address", v.Signer.String(), "ProposerPriority", v.ProposerPriority, "Signer", v.Signer.String(), "Total VotingPower", currentValSet.TotalVotingPower())
		}

		require.Equal(t, len(prevValSet.Validators), len(currentValSet.Validators), "Number of validators should remain same")

		index, _ := currentValSet.GetByAddress(valToUpdate.Signer.Bytes())
		require.Equal(t, -1, index, "Prev Validator should not be present in CurrentValSet")
		index, val := currentValSet.GetByAddress(newSigner[0].Signer.Bytes())
		t.Log("currentValSet - ", currentValSet)
		// require.Equal(t, 0, index, "New Signer should be present in Current val set")
		require.Equal(t, newSigner[0].Signer, val.Signer, "Signer address should change")
		require.Equal(t, newSigner[0].PubKey, val.PubKey, "Signer pubkey should change")
		// require.Equal(t, valToUpdate.ProposerPriority, val.ProposerPriority, "Validator ProposerPriority should not change")
		require.Equal(t, prevValSet.TotalVotingPower(), currentValSet.TotalVotingPower(), "Total VotingPower should not change")
		// TODO not sure if proposer check is needed
		//require.Equal(t, &initValSet.Proposer.Address, &currentValSet.Proposer.Address, "Proposer should not change")
	})

	/* Validator Set changes When
		1. When ackCount changes
		2. When new validator joins
		3. When validator updates stake
		4. When signer is updatedctx
		5. When Validator Exits
	**/

}

// tests setter/getters for Dividend account
func TestDividendAccount(t *testing.T) {
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	dividendAccount := types.DividendAccount{
		ID:            types.NewDividendAccountID(1),
		FeeAmount:     big.NewInt(0).String(),
		SlashedAmount: big.NewInt(0).String(),
	}
	keeper.AddDividendAccount(ctx, dividendAccount)
	ok := keeper.CheckIfDividendAccountExists(ctx, dividendAccount.ID)
	t.Log(ok)

	dividendAccountInStore, _ := keeper.GetDividendAccountByID(ctx, dividendAccount.ID)

	t.Log(dividendAccountInStore)
}

func TestDividendAccountTree(t *testing.T) {

	divAccounts := cmn.GenRandomDividendAccount(3, 1, true)

	accountRoot, _ := checkpointTypes.GetAccountRootHash(divAccounts)
	accountProof, _ := checkpointTypes.GetAccountProof(divAccounts, types.NewDividendAccountID(1))
	leafHash, _ := divAccounts[0].CalculateHash()
	t.Log("accounts", divAccounts)
	t.Log("account root", types.BytesToHeimdallHash(accountRoot))
	t.Log("leaf hash", leafHash)
	t.Log("leaf hash hex", hex.EncodeToString(leafHash))
	t.Log("leaf hash hex bytes", types.BytesToHexBytes(leafHash))
	t.Log("leaft hash heimdall", types.BytesToHeimdallHash(leafHash))
	t.Log("account proof", hex.EncodeToString(accountProof))

}

// func TestDividendAccountHash(t *testing.T) {

// 	divAccounts := cmn.GenRandomDividendAccount(1, 1, true)
// 	accounthash, _ := divAccounts[0].CalculateHash()
// 	t.Log("account hash", accounthash)

// 	_, _ = checkpointTypes.GetRewardRootHash(divAccounts)

// }
