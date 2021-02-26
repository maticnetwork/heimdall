package keeper_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/maticnetwork/heimdall/helper/mocks"

	"github.com/maticnetwork/heimdall/x/staking/test_helper"

	"github.com/maticnetwork/heimdall/helper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/app"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/types/simulation"
	checkPointSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context

	contractCaller *mocks.IContractCaller
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
	suite.contractCaller = &mocks.IContractCaller{}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// tests setter/getters for validatorSignerMaps , validator set/get
func (suite *KeeperTestSuite) TestValidator() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmCommonTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)

		err := initApp.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	// Get random validator ID
	valId := simulation.RandIntBetween(r1, 0, n)

	singerAddress, err := sdk.AccAddressFromHex(validators[valId].Signer)
	if err != nil {
		t.Error("Error getting signer address", err)
	}

	// Get Validator Info from state
	valInfo, err := initApp.StakingKeeper.GetValidatorInfo(ctx, singerAddress)
	if err != nil {
		t.Error("Error while fetching Validator", err)
	}

	// Get Signer Address mapped with ValidatorId
	mappedSignerAddress, isMapped := initApp.StakingKeeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Signer Address not mapped to Validator Id")
	}

	// Check if Validator matches in state
	require.Equal(t, valInfo, *validators[valId], "Validators in state doesnt match")
	require.Equal(t, strings.ToLower(mappedSignerAddress.Hex()), strings.ToLower(validators[0].Signer), "Signer address doesn't match")
}

// tests VotingPower change, validator creation, validator set update when signer changes
func (suite *KeeperTestSuite) TestUpdateSigner() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			1,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmCommonTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
		err := initApp.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	singerAddress, err := sdk.AccAddressFromHex(validators[0].Signer)
	if err != nil {
		t.Error("Error getting signer address", err)
	}

	// Fetch Validator Info from Store
	valInfo, err := initApp.StakingKeeper.GetValidatorInfo(ctx, singerAddress)
	if err != nil {
		t.Error("Error while fetching Validator Info from store", err)
	}
	valInfoSigner, err := sdk.AccAddressFromHex(valInfo.Signer)
	if err != nil {
		t.Error("Error getting validator info", err)
	}

	// Update Signer
	newPrivKey := secp256k1.GenPrivKey()
	newPubKey := hmCommonTypes.NewPubKey(newPrivKey.PubKey().Bytes())
	newSigner := sdk.AccAddress(newPrivKey.PubKey().Address().Bytes())
	newSignerAddress, err := sdk.AccAddressFromHex(newSigner.String())
	if err != nil {
		t.Error("Error getting new signer address", err)
	}

	err = initApp.StakingKeeper.UpdateSigner(ctx, newSignerAddress, newPubKey, valInfoSigner)
	if err != nil {
		t.Error("Error while updating Signer Address -", err)
	}

	// Check Validator Info of Prev Signer
	prevSginerValInfo, err := initApp.StakingKeeper.GetValidatorInfo(ctx, singerAddress)
	if err != nil {
		t.Error("Error while fetching Validator Info for Prev Signer - ", err)
	}
	require.Equal(t, int64(0), prevSginerValInfo.VotingPower, "VotingPower of Prev Signer should be zero")

	// Check Validator Info of Updated Signer
	updatedSignerValInfo, err := initApp.StakingKeeper.GetValidatorInfo(ctx, newSignerAddress)
	if err != nil {
		t.Error("Error while fetching Validator Info for Updater Signer", err)
	}
	require.Equal(t, validators[0].VotingPower, updatedSignerValInfo.VotingPower, "VotingPower of updated signer should match with prev signer VotingPower")

	// Check If ValidatorId is mapped To Updated Signer
	mappedSignerAddress, isMapped := initApp.StakingKeeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Validator Id is not mapped to Signer Address", err)
	}
	require.Equal(t, strings.ToLower(newSignerAddress.String()), strings.ToLower(mappedSignerAddress.Hex()), "Validator ID should be mapped to Updated Signer Address")

	// Check total Validators
	totalValidators := initApp.StakingKeeper.GetAllValidators(ctx)
	require.LessOrEqual(t, 6, len(totalValidators), "Total Validators should be six.")

	// Check current Validators
	currentValidators := initApp.StakingKeeper.GetCurrentValidators(ctx)
	require.LessOrEqual(t, 5, len(currentValidators), "Current Validators should be five.")
}

func (suite *KeeperTestSuite) TestCurrentValidator() {
	type TestDataItem struct {
		name        string
		startblock  uint64
		nonce       uint64
		VotingPower int64
		ackcount    uint64
		result      bool
		resultmsg   string
	}

	dataItems := []TestDataItem{
		{"VotingPower zero", uint64(0), uint64(0), int64(0), uint64(1), false, "should not be current validator as VotingPower is zero."},
		{"start epoch greater than ackcount", uint64(3), uint64(0), int64(10), uint64(1), false, "should not be current validator as start epoch greater than ackcount."},
	}
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx

	stakingKeeper, checkpointKeeper := intiApp.StakingKeeper, intiApp.CheckpointKeeper
	for i, item := range dataItems {
		suite.Run(item.name, func() {
			// Create a Validator [startEpoch, endEpoch, VotingPower]
			privKep := secp256k1.GenPrivKey()
			pubkey := hmCommonTypes.NewPubKey(privKep.PubKey().Bytes())
			newVal := hmTypes.Validator{
				ID:               hmTypes.NewValidatorID(1 + uint64(i)),
				StartEpoch:       item.startblock,
				EndEpoch:         item.startblock,
				Nonce:            item.nonce,
				VotingPower:      item.VotingPower,
				Signer:           hmCommonTypes.HexToHeimdallAddress(pubkey.Address().String()).String(),
				PubKey:           pubkey.String(),
				ProposerPriority: 0,
			}
			// check current validator
			err := stakingKeeper.AddValidator(ctx, newVal)
			require.NoError(t, err)
			checkpointKeeper.UpdateACKCountWithValue(ctx, item.ackcount)

			isCurrentVal := stakingKeeper.IsCurrentValidatorByAddress(ctx, GetAccAddressFromString(newVal.Signer))
			require.Equal(t, item.result, isCurrentVal, item.resultmsg)
		})
	}
}

// // GetAccAddressFromString return sdk.AccAddress from return
func GetAccAddressFromString(address string) sdk.AccAddress {
	return sdk.AccAddress([]byte(address))
}

func (suite *KeeperTestSuite) TestRemoveValidatorSetChange() {
	// create sub test to check if validator remove
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper

	// load 4 validators to state
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	initValSet := keeper.GetValidatorSet(ctx)

	currentValSet := initValSet.Copy()
	prevValidatorSet := initValSet.Copy()

	prevValidatorSet.Validators[0].StartEpoch = 20

	validator := *prevValidatorSet.Validators[0]
	err := keeper.AddValidator(ctx, validator)
	require.Empty(t, err, "Unable to update validator set")

	setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
	err = currentValSet.UpdateWithChangeSet(setUpdates)
	require.NoError(t, err)
	updatedValSet := currentValSet

	require.Equal(t, len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators), "Validator set should be reduced by one ")

	for _, val := range updatedValSet.Validators {
		if val.Signer == prevValidatorSet.Validators[0].Signer {
			require.Fail(t, "Validator is not removed from updated validator set")
		}
	}
}

func (suite *KeeperTestSuite) TestAddValidatorSetChange() {
	// create sub test to check if validator remove
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper

	// load 4 validators to state
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	initValSet := keeper.GetValidatorSet(ctx)

	validators := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)
	prevValSet := initValSet.Copy()

	valToBeAdded := validators[0]
	currentValSet := initValSet.Copy()

	err := keeper.AddValidator(ctx, valToBeAdded)
	require.NoError(t, err)
	setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
	err = currentValSet.UpdateWithChangeSet(setUpdates)

	require.NoError(t, err)
	require.Equal(t, len(prevValSet.Validators)+1, len(currentValSet.Validators), "Number of validators should be increased by 1")
	require.Equal(t, true, currentValSet.HasAddress(valToBeAdded.GetSigner().Bytes()), "New Validator should be added")
	require.Equal(t, prevValSet.GetTotalVotingPower()+int64(valToBeAdded.VotingPower), currentValSet.GetTotalVotingPower(), "Total VotingPower should be increased")

}

func (suite *KeeperTestSuite) TestUpdateValidatorSetChange() {
	// create sub test to check if validator remove
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper

	// load 4 validators to state
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	initValSet := keeper.GetValidatorSet(ctx)

	keeper.IncrementAccum(ctx, 2)
	prevValSet := initValSet.Copy()
	currentValSet := keeper.GetValidatorSet(ctx)

	valToUpdate := currentValSet.Validators[0]
	newSigner := stakingSim.GenRandomVal(1, 0, 10, 10, false, 1)

	previousSigner, err := sdk.AccAddressFromHex(valToUpdate.Signer)
	require.NoError(t, err)
	newSignerAddr, _ := sdk.AccAddressFromHex(newSigner[0].Signer)
	err = keeper.UpdateSigner(ctx, newSignerAddr, hmCommonTypes.NewPubKeyFromHex(newSigner[0].PubKey), previousSigner)
	require.NoError(t, err)
	setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
	err = currentValSet.UpdateWithChangeSet(setUpdates)
	require.NoError(t, err)
	require.Equal(t, len(prevValSet.Validators), len(currentValSet.Validators), "Number of validators should remain same")

	index, _ := currentValSet.GetByAddress(previousSigner)
	require.Equal(t, int32(-1), index, "Prev Validator should not be present in CurrentValSet")
	index, val := currentValSet.GetByAddress(newSignerAddr)
	require.NotNil(t, index)
	require.Equal(t, newSigner[0].GetSigner(), val.GetSigner(), "Signer address should change")
	require.Equal(t, newSigner[0].PubKey, val.PubKey, "Signer pubkey should change")

	require.Equal(t, prevValSet.GetTotalVotingPower(), currentValSet.GetTotalVotingPower(), "Total VotingPower should not change")

	/* Validator Set changes When
		1. When ackCount changes
		2. When new validator joins
		3. When validator updates stake
		4. When signer is updatedctx
		5. When Validator Exits
	**/

}

//
func (suite *KeeperTestSuite) TestGetCurrentValidators() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.StakingKeeper
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)
	activeValidatorInfo, err := keeper.GetActiveValidatorInfo(ctx, validators[0].GetSigner())
	require.NoError(t, err)
	require.Equal(t, validators[0], activeValidatorInfo)
}

func (suite *KeeperTestSuite) TestGetCurrentProposer() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	currentValSet := keeper.GetValidatorSet(ctx)
	currentProposer := keeper.GetCurrentProposer(ctx)
	require.Equal(t, currentValSet.GetProposer(), currentProposer)
}

func (suite *KeeperTestSuite) TestGetNextProposer() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx

	checkPointSim.LoadValidatorSet(4, t, intiApp.StakingKeeper, ctx, false, 10)

	nextProposer := intiApp.StakingKeeper.GetNextProposer(ctx)
	require.NotNil(t, nextProposer)
}

func (suite *KeeperTestSuite) TestGetValidatorFromValID() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)

	valInfo, ok := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.Equal(t, ok, true)
	require.Equal(t, validators[0], valInfo)
}

func (suite *KeeperTestSuite) TestGetLastUpdated() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)

	lastUpdated, ok := keeper.GetLastUpdated(ctx, validators[0].ID)
	require.Equal(t, ok, true)
	require.Equal(t, validators[0].LastUpdated, lastUpdated)
}

func (suite *KeeperTestSuite) TestGetSpanEligibleValidators() {
	t, intiApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := intiApp.StakingKeeper
	checkPointSim.LoadValidatorSet(4, t, keeper, ctx, false, 0)

	// Test ActCount = 0
	intiApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 0)

	valActCount0 := keeper.GetSpanEligibleValidators(ctx)
	require.LessOrEqual(t, len(valActCount0), 4)

	intiApp.CheckpointKeeper.UpdateACKCountWithValue(ctx, 20)

	validators := keeper.GetSpanEligibleValidators(ctx)
	require.LessOrEqual(t, len(validators), 4)
}
