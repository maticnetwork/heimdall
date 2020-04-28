package staking_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	cmn "github.com/maticnetwork/heimdall/test"

	"github.com/maticnetwork/heimdall/helper"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// tests setter/getters for validatorSignerMaps , validator set/get
func (suite *KeeperTestSuite) TestValidator() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)
	dividendAccounts := make([]hmTypes.DividendAccount, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
		// create dividend account for validator
		dividendAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.NewDividendAccountID(uint64(validators[i].ID)),
			big.NewInt(0).String(),
			big.NewInt(0).String(),
		)
		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	// Get Validator Info from state
	valInfo, err := app.StakingKeeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator", err)
	}
	// Get Signer Address mapped with ValidatorId
	mappedSignerAddress, isMapped := app.StakingKeeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Signer Address not mapped to Validator Id")
	}
	// Check if Validator matches in state
	require.Equal(t, valInfo, *validators[0], "Validators in state doesnt match")
	require.Equal(t, types.HexToHeimdallAddress(mappedSignerAddress.Hex()), validators[0].Signer, "Signer address doesn't match")
}

// tests VotingPower change, validator creation, validator set update when signer changes
func (suite *KeeperTestSuite) TestUpdateSigner() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	n := 5

	validators := make([]*hmTypes.Validator, n)
	accounts := simulation.RandomAccounts(r1, n)
	dividendAccounts := make([]hmTypes.DividendAccount, n)

	for i := range validators {
		// validator
		validators[i] = hmTypes.NewValidator(
			hmTypes.NewValidatorID(uint64(int64(i))),
			0,
			0,
			int64(simulation.RandIntBetween(r1, 10, 100)), // power
			hmTypes.NewPubKey(accounts[i].PubKey.Bytes()),
			accounts[i].Address,
		)
		// create dividend account for validator
		dividendAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.NewDividendAccountID(uint64(validators[i].ID)),
			big.NewInt(0).String(),
			big.NewInt(0).String(),
		)
		err := app.StakingKeeper.AddValidator(ctx, *validators[i])
		if err != nil {
			t.Error("Error while adding validator to store", err)
		}
	}

	// Fetch Validator Info from Store
	valInfo, err := app.StakingKeeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info from store", err)
	}

	// Update Signer
	newPrivKey := secp256k1.GenPrivKey()
	newPubKey := types.NewPubKey(newPrivKey.PubKey().Bytes())
	newSigner := types.HexToHeimdallAddress(newPubKey.Address().String())
	err = app.StakingKeeper.UpdateSigner(ctx, newSigner, newPubKey, valInfo.Signer)
	if err != nil {
		t.Error("Error while updating Signer Address -", err)
	}

	// Check Validator Info of Prev Signer
	prevSginerValInfo, err := app.StakingKeeper.GetValidatorInfo(ctx, validators[0].Signer.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info for Prev Signer - ", err)
	}
	require.Equal(t, int64(0), prevSginerValInfo.VotingPower, "VotingPower of Prev Signer should be zero")

	// Check Validator Info of Updated Signer
	updatedSignerValInfo, err := app.StakingKeeper.GetValidatorInfo(ctx, newSigner.Bytes())
	if err != nil {
		t.Error("Error while fetching Validator Info for Updater Signer", err)
	}
	require.Equal(t, validators[0].VotingPower, updatedSignerValInfo.VotingPower, "VotingPower of updated signer should match with prev signer VotingPower")

	// Check If ValidatorId is mapped To Updated Signer
	signerAddress, isMapped := app.StakingKeeper.GetSignerFromValidatorID(ctx, validators[0].ID)
	if !isMapped {
		t.Error("Validator Id is not mapped to Signer Address", err)
	}
	require.Equal(t, newSigner, types.HexToHeimdallAddress(signerAddress.Hex()), "Validator ID should be mapped to Updated Signer Address")

	// Check total Validators
	totalValidators := app.StakingKeeper.GetAllValidators(ctx)
	require.LessOrEqual(t, 6, len(totalValidators), "Total Validators should be six.")
}

func (suite *KeeperTestSuite) TestCurrentValidator() {
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
	t, app, ctx := suite.T(), suite.app, suite.ctx

	stakingKeeper, checkpointKeeper := app.StakingKeeper, app.CheckpointKeeper
	for i, item := range dataItems {
		suite.Run(item.name, func() {
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

func (suite *KeeperTestSuite) TestRemoveValidatorSetChange() {
	// create sub test to check if validator remove
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper

	// load 4 validators to state
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	initValSet := keeper.GetValidatorSet(ctx)

	currentValSet := initValSet.Copy()
	prevValidatorSet := initValSet.Copy()
	fmt.Println("currentValSet", len(currentValSet.Validators))
	fmt.Println("prevValidatorSet", len(prevValidatorSet.Validators))

	// remove validator (making IsCurrentValidator return false)
	fmt.Println("Updating", prevValidatorSet.Validators[0].ID.String(), prevValidatorSet.Validators[0].StartEpoch)

	prevValidatorSet.Validators[0].StartEpoch = 20
	fmt.Println("remove validator  currentValSet", len(currentValSet.Validators))
	fmt.Println("remove validator  prevValidatorSet", len(prevValidatorSet.Validators))
	fmt.Println("--->1 keeper.GetAllValidators(ctx)", len(keeper.GetAllValidators(ctx)))
	fmt.Println("Updated Validators in state", prevValidatorSet.Validators[0].ID.String(), prevValidatorSet.Validators[0].StartEpoch)

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
	fmt.Println("------------------------------------------------")
	fmt.Println("urrentValSet", currentValSet.Validators)
	fmt.Println("------------------------------------------------")

	setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
	currentValSet.UpdateWithChangeSet(setUpdates)
	fmt.Println("------------------------------------------------")
	fmt.Println("round 2 urrentValSet", len(currentValSet.Validators), len(setUpdates), currentValSet.Validators)
	fmt.Println("------------------------------------------------")
	updatedValSet := currentValSet
	t.Log("Validators in updated validator set")
	for _, v := range updatedValSet.Validators {
		t.Log("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
	}
	// check if 1 validator is removed
	fmt.Println("ohh fuck", len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators))
	require.Equal(t, len(prevValidatorSet.Validators)-1, len(updatedValSet.Validators), "Validator set should be reduced by one ")
	// remove first validator from initial validator set and equate with new
	t.Log("appended set-", updatedValSet, "one", prevValidatorSet.Validators[1:])
	for _, val := range updatedValSet.Validators {
		if val.Signer == prevValidatorSet.Validators[0].Signer {
			require.Fail(t, "Validator is not removed from updatedvalidator set")
		}
	}

}

func (suite *KeeperTestSuite) TestAddValidatorSetChange() {
	// create sub test to check if validator remove
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper

	// load 4 validators to state
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	initValSet := keeper.GetValidatorSet(ctx)

	validators := cmn.GenRandomVal(1, 0, 10, 10, false, 1)
	prevValSet := initValSet.Copy()
	fmt.Println("in add keeper.GetAllValidators(ctx)", len(initValSet.Validators))

	valToBeAdded := validators[0]
	currentValSet := initValSet.Copy()
	//prevValidatorSet := initValSet.Copy()
	keeper.AddValidator(ctx, valToBeAdded)
	fmt.Println("currentValSet new added keeper.GetAllValidators(ctx)", len(keeper.GetValidatorSet(ctx).Validators))
	fmt.Println("Val to be Added")
	fmt.Println("-->", "Address", valToBeAdded.Signer.String(), "StartEpoch", valToBeAdded.StartEpoch, "EndEpoch", valToBeAdded.EndEpoch, "VotingPower", valToBeAdded.VotingPower)

	// helper.UpdateValidators(
	// 	currentValSet,                // pointer to current validator set -- UpdateValidators will modify it
	// 	keeper.GetAllValidators(ctx), // All validators
	// 	5,                            // ack count
	// )
	setUpdates := helper.GetUpdatedValidators(currentValSet, keeper.GetAllValidators(ctx), 5)
	currentValSet.UpdateWithChangeSet(setUpdates)

	fmt.Println("Validators in updated validator set", len(currentValSet.Validators))
	// for _, v := range currentValSet.Validators {
	// 	fmt.Println("-->", "Address", v.Signer.String(), "StartEpoch", v.StartEpoch, "EndEpoch", v.EndEpoch, "VotingPower", v.VotingPower)
	// }
	fmt.Println("venky currentValSet GetUpdatedValidators", len(prevValSet.Validators)+1, len(currentValSet.Validators))

	require.Equal(t, len(prevValSet.Validators)+1, len(currentValSet.Validators), "Number of validators should be increased by 1")
	require.Equal(t, true, currentValSet.HasAddress(valToBeAdded.Signer.Bytes()), "New Validator should be added")
	require.Equal(t, prevValSet.TotalVotingPower()+int64(valToBeAdded.VotingPower), currentValSet.TotalVotingPower(), "Total VotingPower should be increased")

}

func (suite *KeeperTestSuite) TestUpdateValidatorSetChange() {
	// create sub test to check if validator remove
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper

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

	/* Validator Set changes When
		1. When ackCount changes
		2. When new validator joins
		3. When validator updates stake
		4. When signer is updatedctx
		5. When Validator Exits
	**/

}

// tests setter/getters for Dividend account
func (suite *KeeperTestSuite) TestDividendAccount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	dividendAccount := types.DividendAccount{
		ID:            types.NewDividendAccountID(1),
		FeeAmount:     big.NewInt(0).String(),
		SlashedAmount: big.NewInt(0).String(),
	}
	app.StakingKeeper.AddDividendAccount(ctx, dividendAccount)
	ok := app.StakingKeeper.CheckIfDividendAccountExists(ctx, dividendAccount.ID)
	t.Log(ok)

	dividendAccountInStore, _ := app.StakingKeeper.GetDividendAccountByID(ctx, dividendAccount.ID)

	t.Log(dividendAccountInStore)
}

func (suite *KeeperTestSuite) TestDividendAccountTree() {
	t := suite.T()

	divAccounts := make([]hmTypes.DividendAccount, 5)
	for i := 0; i < len(divAccounts); i++ {
		divAccounts[i] = hmTypes.NewDividendAccount(
			hmTypes.NewDividendAccountID(uint64(1)),
			big.NewInt(0).String(),
			big.NewInt(0).String(),
		)
	}

	accountRoot, _ := checkpointTypes.GetAccountRootHash(divAccounts)
	accountProof, _, _ := checkpointTypes.GetAccountProof(divAccounts, types.NewDividendAccountID(1))
	leafHash, _ := divAccounts[0].CalculateHash()
	t.Log("accounts", divAccounts)
	t.Log("account root", types.BytesToHeimdallHash(accountRoot))
	t.Log("leaf hash", leafHash)
	t.Log("leaf hash hex", hex.EncodeToString(leafHash))
	t.Log("leaf hash hex bytes", types.BytesToHexBytes(leafHash))
	t.Log("leaft hash heimdall", types.BytesToHeimdallHash(leafHash))
	t.Log("account proof", hex.EncodeToString(accountProof))

}

func (suite *KeeperTestSuite) TestGetCurrentValidators() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)
	activeValidatorInfo, err := keeper.GetActiveValidatorInfo(ctx, validators[0].Signer.Bytes())
	require.NoError(t, err)
	require.Equal(t, validators[0], activeValidatorInfo)
}

func (suite *KeeperTestSuite) TestGetCurrentProposer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	currentValSet := keeper.GetValidatorSet(ctx)
	currentProposer := keeper.GetCurrentProposer(ctx)
	require.Equal(t, currentValSet.GetProposer(), currentProposer)
}

func (suite *KeeperTestSuite) TestGetNextProposer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)

	nextProposer := keeper.GetNextProposer(ctx)
	require.NotNil(t, nextProposer)
}

func (suite *KeeperTestSuite) TestGetValidatorFromValID() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)

	valInfo, ok := keeper.GetValidatorFromValID(ctx, validators[0].ID)
	require.Equal(t, ok, true)
	require.Equal(t, validators[0], valInfo)
}

func (suite *KeeperTestSuite) TestGetLastUpdated() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)

	lastUpdated, ok := keeper.GetLastUpdated(ctx, validators[0].ID)
	require.Equal(t, ok, true)
	require.Equal(t, validators[0].LastUpdated, lastUpdated)
}

func (suite *KeeperTestSuite) TestAddFeeToDividendAccount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetCurrentValidators(ctx)

	amount, _ := big.NewInt(0).SetString("0", 10)
	keeper.AddFeeToDividendAccount(ctx, validators[0].ID, amount)
	dividentAccountId := hmTypes.DividendAccountID(validators[0].ID)
	dividentAccount, _ := keeper.GetDividendAccountByID(ctx, dividentAccountId)
	actualResult, ok := big.NewInt(0).SetString(dividentAccount.FeeAmount, 10)
	require.Equal(t, ok, true)
	require.Equal(t, amount, actualResult)
}

func (suite *KeeperTestSuite) TestGetSpanEligibleValidators() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.StakingKeeper
	cmn.LoadValidatorSet(4, t, keeper, ctx, false, 10)
	validators := keeper.GetSpanEligibleValidators(ctx)
	fmt.Println("validators", len(validators))
	// TODO: Change this later
	require.LessOrEqual(t, 4, len(validators))
}
