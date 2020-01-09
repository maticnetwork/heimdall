package staking_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
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

// // Tests setters and getters for validator reward
// func TestValidatorRewards(t *testing.T) {
// 	ctx, keeper, _ := cmn.CreateTestInput(t, false)
// 	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
// 	curVal := keeper.GetCurrentValidators(ctx)
// 	// check initial reward
// 	initReward := big.NewInt(100)
// 	keeper.SetValidatorIDToReward(ctx, curVal[0].ID, initReward)
// 	valReward, err := keeper.GetRewardByDividendAccountID(ctx, types.DividendAccountID(curVal[0].ID))
// 	require.Equal(t, initReward, valReward, "Validator Initial Reward should be %v but it is %v", initReward, valReward)
// 	// check updated reward
// 	rewardAdded := big.NewInt(50)
// 	keeper.SetValidatorIDToReward(ctx, curVal[0].ID, rewardAdded)
// 	updatedReward, err := keeper.GetRewardByDividendAccountID(ctx, types.DividendAccountID(curVal[0].ID))
// 	rewardSum := big.NewInt(0).Add(initReward, rewardAdded)
// 	require.Equal(t, rewardSum, updatedReward, "Validator Updated Reward should be %v but it is %v", rewardSum, updatedReward)
// 	// zero reward for Invalid Validator ID
// 	rewardNonValId, err := keeper.GetRewardByDividendAccountID(ctx, types.DividendAccountID(curVal[1].ID))
// 	require.Equal(t, big.NewInt(0), rewardNonValId, "Reward should be zero but it is %v", rewardNonValId)
// 	// check validator reward map
// 	keeper.SetValidatorIDToReward(ctx, curVal[1].ID, big.NewInt(35))
// 	keeper.SetValidatorIDToReward(ctx, curVal[2].ID, big.NewInt(45))
// 	valRewardMap := keeper.GetAllDividendAccounts(ctx)
// 	t.Log("Validator Reward Map - ", valRewardMap)
// 	require.Equal(t, 3, len(valRewardMap), "Validator Reward map size should be %v but it is %v", 3, len(valRewardMap))
// 	require.Equal(t, rewardSum, valRewardMap[curVal[0].ID], "Validator Reward should be %v but it is %v", rewardSum, valRewardMap[curVal[0].ID])
// 	require.Equal(t, big.NewInt(35), valRewardMap[curVal[1].ID], "Validator Reward should be %v but it is %v", big.NewInt(35), valRewardMap[curVal[0].ID])
// 	require.Equal(t, big.NewInt(45), valRewardMap[curVal[2].ID], "Validator Reward should be %v but it is %v", big.NewInt(45), valRewardMap[curVal[0].ID])

// 	// Generate Merkle Root Out of Rewards after sorting by valID
// 	rewardRootHash, err := checkpointTypes.GetRewardRootHash(valRewardMap)
// 	require.Empty(t, err, "Error when generating reward root hash from validator reward state tree")
// 	t.Log("Reward root hash - ", types.BytesToHeimdallHash(rewardRootHash))

// }

func TestCalculateSignerRewards(t *testing.T) {
	ctx, keeper, _ := cmn.CreateTestInput(t, false)
	checkpointReward := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(22), nil)
	t.Log("checkpoint reward - ", checkpointReward)
	keeper.SetProposerBonusPercent(ctx, stakingTypes.DefaultProposerBonusPercent)
	var valSet = types.ValidatorSet{}
	var newVal = types.Validator{}
	signerRewardshouldbe := make(map[types.ValidatorID]*big.Int)
	signerRewardshouldbe[1], _ = big.NewInt(0).SetString("1500000000000000000000", 10)
	signerRewardshouldbe[2], _ = big.NewInt(0).SetString("3000000000000000000000", 10)
	signerRewardshouldbe[3], _ = big.NewInt(0).SetString("5500000000000000393216", 10)
	// These are pubkeys and signer address for below submitheaderblock trasaction payload
	pubKeys := []string{"045b608112c8d9ca26f50ede110495e6be48cf9bb6d220d0354e3771701d3c9b1c8805d039e194f8938b820fee6d0aff4e2120b385f1e58b62d8649a796e7433a2", "041e8bc59b9c58358c9f2847d9dc62b927bc3fc7ac83e9b1a38b402ffb6d6d2d7be9329f51e6a9a4cfb75bc426def46be8f847a4c2fd335be55f382a08d4f3325a", "047ad78e23df40cecc5c6adf661df02d103aff74a95e2c4de99b1d0855b67d2881c659d7831daeae2c7626b60575a9a4aae62bd1ea225c1f71cb2c63c63a7de4a0"}
	signerAddresses := []string{"a03d8f5af7413e4fd5a37fde9286e390ef8f3c07", "b1bf4473c6b1918a6e37408e1c14df81281411a8", "ba754e3893adb3cabc0afe7932b4b5a3cee3f3ab"}
	// Add these validators to store
	for i := 0; i < len(signerAddresses); i++ {
		newVal = types.Validator{
			ID:               types.NewValidatorID(uint64(i + 1)),
			StartEpoch:       0,
			EndEpoch:         100,
			VotingPower:      int64(i) + 1,
			Signer:           types.HexToHeimdallAddress(signerAddresses[i]),
			PubKey:           types.NewPubKey([]byte(pubKeys[i])),
			ProposerPriority: 0,
		}
		keeper.AddValidator(ctx, newVal)
		valSet.UpdateWithChangeSet([]*types.Validator{&newVal})
	}

	// Add extra validator not part of signer. This is to make sure totalStakePower != signerPower
	nonSignerVals := cmn.GenRandomVal(1, 0, 4, uint64(10), true, 1)
	nonSignerVals[0].ID = types.NewValidatorID(uint64(4))
	keeper.AddValidator(ctx, nonSignerVals[0])
	valSet.UpdateWithChangeSet([]*types.Validator{&nonSignerVals[0]})

	// Set one of the signer as the proposer
	valSet.Proposer = &newVal
	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")

	// updatedValSet := keeper.GetValidatorSet(ctx)
	// t.Log(updatedValSet)

	// Unpack Signers from paylaod
	data := string(rootchain.RootchainABI)
	abi, err := abi.JSON(strings.NewReader(data))
	require.Empty(t, err, "Error while getting RootChainABI")
	payload := "ec83d3ba000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000030ef8f6865696d64616c6c2d39337251774b84766f7465820cb0800294907eb68cd3480777e3fde8897fb1373de6e982cc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c3c65be37c47302d110fbe6453ef84cc69ecc744d9dbd6fa98d68a621541e2e8291758b962f41ed31a1f8c25db2631da2f49fe20ff1e2d683dd0d802fcef6928d5000622073cfbc99994cd06d7a7a8b01e453b57495010d6eb312a68b00ca6f581d0729f494a385bd1fbe2fd8df3da706fa85a54694ab0d9a4177555048aaa7b3371005c6dd42d128482e603c5adc7cfca0f0c730b49d6bd8ba750307d497c21097c922479f51151c53b35f20a1a1cb4790afa95e470b8f70b0318726d6175b9055b340000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045f84394b1bf4473c6b1918a6e37408e1c14df81281411a883543ff8835441f7a07d2842c3044740cfe1e1a5f782bfd3b91de0c634e9933524b5e3daacc854f49b845d94796f000000000000000000000000000000000000000000000000000000"
	decodedPayload, err := hex.DecodeString(payload)
	require.Empty(t, err, "Error while decoding payload")
	voteSignBytes, inputSigs, txData, err := helper.UnpackSigAndVotes(decodedPayload, abi)
	require.Empty(t, err, "Error while unpacking payload")
	t.Log("voteSignBytes", hex.EncodeToString(voteSignBytes))
	t.Log("inputSigs", hex.EncodeToString(inputSigs))
	t.Log("txData", hex.EncodeToString(txData))

	// Calculate Rewards for Signers
	signerRewardMap, err := keeper.CalculateSignerRewards(ctx, voteSignBytes, inputSigs, checkpointReward)
	t.Log("Signer Reward Map - ", signerRewardMap)
	require.Empty(t, err, "Error while calculating rewards for signers", err)
	require.Equal(t, len(pubKeys), len(signerRewardMap), "No of signers should be %v but it is %v", len(pubKeys), len(signerRewardMap))

	// Verify Rewards for validator Signatures
	for i := 0; i < len(signerRewardMap); i++ {
		val, err := keeper.GetValidatorInfo(ctx, types.HexToHeimdallAddress(signerAddresses[i]).Bytes())
		require.Empty(t, err, "Error while getting val info for signer -", types.HexToHeimdallAddress(signerAddresses[i]).Bytes())
		require.Equal(t, signerRewardshouldbe[val.ID], signerRewardMap[val.ID], "Reward for valId %v should be %v but it %v", val.ID, signerRewardshouldbe[val.ID], signerRewardMap[val.ID])
	}
}
