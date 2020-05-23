package slashing_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/slashing"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"

	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
)

//
// Test suite
//

// KeeperTestSuite integrate test suite context object
type KeeperTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	sideHandler    hmTypes.SideTxHandler
	postHandler    hmTypes.PostTxHandler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = slashing.NewSideTxHandler(suite.app.SlashingKeeper, &suite.contractCaller)
	suite.postHandler = slashing.NewPostTxHandler(suite.app.SlashingKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

//
// Tests
//

func (suite *KeeperTestSuite) TestGetSetValidatorSigningInfo() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.SlashingKeeper

	valID := hmTypes.NewValidatorID(uint64(1))

	info, found := keeper.GetValidatorSigningInfo(ctx, valID)
	require.False(t, found)

	newInfo := hmTypes.NewValidatorSigningInfo(
		valID,
		int64(4),
		int64(3),
		int64(10),
	)
	keeper.SetValidatorSigningInfo(ctx, valID, newInfo)
	info, found = keeper.GetValidatorSigningInfo(ctx, valID)
	require.True(t, found)
	require.Equal(t, info.StartHeight, int64(4))
	require.Equal(t, info.IndexOffset, int64(3))
	require.Equal(t, info.MissedBlocksCounter, int64(10))
}

func (suite *KeeperTestSuite) TestGetSetValidatorMissedBlockBitArray() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	slashingKeeper := app.SlashingKeeper
	valID := hmTypes.NewValidatorID(uint64(1))
	index := int64(0)
	missed := false

	response := slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)
	require.Equal(t, missed, response) // treat empty key as not missed

	slashingKeeper.SetValidatorMissedBlockBitArray(ctx, valID, index, missed)
	response = slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)
	require.Equal(t, missed, response)

	missed = true
	index = int64(1)
	slashingKeeper.SetValidatorMissedBlockBitArray(ctx, valID, index, missed)
	response = slashingKeeper.GetValidatorMissedBlockBitArray(ctx, valID, index)
	require.Equal(t, missed, response) // now should be missed
}

// Test a new validator entering the validator set
// Ensure that SigningInfo.StartHeight is set correctly
// and that they are not immediately jailed
func (suite *KeeperTestSuite) TestHandleNewValidator() {
	// initial setup
	t, app, ctx := suite.T(), suite.app, suite.ctx
	ctx.WithBlockHeight(100)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	validatorId := uint64(1)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	blockNumber := big.NewInt(10)
	nonce := big.NewInt(3)

	slashingKeeper := app.SlashingKeeper
	// params := slashingKeeper.GetParams(ctx)

	stakingPostHandler := staking.NewPostTxHandler(suite.app.StakingKeeper, &suite.contractCaller)

	msgValJoin := stakingTypes.NewMsgValidatorJoin(
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		validatorId,
		uint64(1),
		sdk.NewInt(1000000000000000000),
		pubkey,
		txHash,
		logIndex,
		blockNumber.Uint64(),
		nonce.Uint64(),
	)

	result := stakingPostHandler(ctx, msgValJoin, abci.SideTxResultType_Yes)
	require.True(t, result.IsOK(), "expected validator join to be ok, got %v", result)

	valSigningInfo, found := slashingKeeper.GetValidatorSigningInfo(ctx, hmTypes.NewValidatorID(validatorId))
	require.True(t, found)
	require.Equal(t, ctx.BlockHeight(), valSigningInfo.StartHeight, "start height mismatch")
	require.Equal(t, int64(0), valSigningInfo.IndexOffset, "index offset should be zero")
	require.Equal(t, int64(0), valSigningInfo.MissedBlocksCounter, "missed block counter should be zero")
}

// Test a validator through uptime, downtime, revocation,
// unrevocation, starting height reset, and revocation again
func (suite *KeeperTestSuite) TestHandleAbsentValidator() {

	// initial setup
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	txHash := hmTypes.HexToHeimdallHash("123")
	index := simulation.RandIntBetween(r1, 0, 100)
	logIndex := uint64(index)
	validatorId := uint64(1)

	privKey1 := secp256k1.GenPrivKey()
	pubkey := hmTypes.NewPubKey(privKey1.PubKey().Bytes())
	address := pubkey.Address()

	blockNumber := big.NewInt(10)
	nonce := big.NewInt(3)

	slashingKeeper := app.SlashingKeeper
	params := slashingKeeper.GetParams(ctx)
	params.SignedBlocksWindow = 100
	stakingPostHandler := staking.NewPostTxHandler(suite.app.StakingKeeper, &suite.contractCaller)
	power := int64(1000)
	bigPower, _ := helper.GetAmountFromPower(power)
	msgValJoin := stakingTypes.NewMsgValidatorJoin(
		hmTypes.BytesToHeimdallAddress(address.Bytes()),
		validatorId,
		uint64(1),
		sdk.NewIntFromBigInt(bigPower),
		pubkey,
		txHash,
		logIndex,
		blockNumber.Uint64(),
		nonce.Uint64(),
	)

	result := stakingPostHandler(ctx, msgValJoin, abci.SideTxResultType_Yes)
	require.True(t, result.IsOK(), "expected validator join to be ok, got %v", result)

	val, _ := app.StakingKeeper.GetValidatorInfo(ctx, address.Bytes())
	valSigningInfo, found := slashingKeeper.GetValidatorSigningInfo(ctx, hmTypes.NewValidatorID(validatorId))
	require.True(t, found)
	require.Equal(t, ctx.BlockHeight(), valSigningInfo.StartHeight, "start height mismatch")
	require.Equal(t, int64(0), valSigningInfo.IndexOffset, "index offset should be zero")
	require.Equal(t, int64(0), valSigningInfo.MissedBlocksCounter, "missed block counter should be zero")

	height := int64(0)

	// 1000 first blocks OK
	for ; height < params.SignedBlocksWindow; height++ {
		ctx = ctx.WithBlockHeight(height)
		slashingKeeper.HandleValidatorSignature(ctx, val.Signer.Bytes(), val.VotingPower, true)
	}

	// 500 blocks missed
	for ; height < params.SignedBlocksWindow+(params.SignedBlocksWindow-slashingKeeper.MinSignedPerWindow(ctx)); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashingKeeper.HandleValidatorSignature(ctx, val.Signer.Bytes(), val.VotingPower, false)
	}
	info, found := slashingKeeper.GetValidatorSigningInfo(ctx, val.ID)
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	require.Equal(t, params.SignedBlocksWindow-slashingKeeper.MinSignedPerWindow(ctx), info.MissedBlocksCounter)

	// 501st block missed
	ctx = ctx.WithBlockHeight(height)
	slashingKeeper.HandleValidatorSignature(ctx, val.Signer.Bytes(), val.VotingPower, false)
	info, found = slashingKeeper.GetValidatorSigningInfo(ctx, val.ID)
	require.True(t, found)
	require.Equal(t, int64(0), info.StartHeight)
	// counter now reset to zero
	require.Equal(t, int64(0), info.MissedBlocksCounter)

	bufferSlashInfo, found := slashingKeeper.GetBufferValSlashingInfo(ctx, val.ID)
	require.True(t, found)
	expectedSlashAmount := sdk.NewDec(power).Mul(params.SlashFractionDowntime).TruncateInt().Int64()
	require.Equal(t, bufferSlashInfo.SlashedAmount, uint64(expectedSlashAmount), "slash amount mismatches")

	for ; height < params.SignedBlocksWindow+(params.SignedBlocksWindow-slashingKeeper.MinSignedPerWindow(ctx)); height++ {
		ctx = ctx.WithBlockHeight(height)
		slashingKeeper.HandleValidatorSignature(ctx, val.Signer.Bytes(), val.VotingPower, false)
	}

	/* 	// validator should have been jailed
	   	validator, _ = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	   	require.Equal(t, sdk.Unbonding, validator.GetStatus())

	   	slashAmt := amt.ToDec().Mul(keeper.SlashFractionDowntime(ctx)).RoundInt64()

	   	// validator should have been slashed
	   	require.Equal(t, amt.Int64()-slashAmt, validator.GetTokens().Int64())

	   	// 502nd block *also* missed (since the LastCommit would have still included the just-unbonded validator)
	   	height++
	   	ctx = ctx.WithBlockHeight(height)
	   	keeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	   	info, found = keeper.getValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	   	require.True(t, found)
	   	require.Equal(t, int64(0), info.StartHeight)
	   	require.Equal(t, int64(1), info.MissedBlocksCounter)

	   	// end block
	   	staking.EndBlocker(ctx, sk)

	   	// validator should not have been slashed any more, since it was already jailed
	   	validator, _ = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	   	require.Equal(t, amt.Int64()-slashAmt, validator.GetTokens().Int64())

	   	// unrevocation should fail prior to jail expiration
	   	got = slh(ctx, NewMsgUnjail(addr))
	   	require.False(t, got.IsOK())

	   	// unrevocation should succeed after jail expiration
	   	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Unix(1, 0).Add(keeper.DowntimeJailDuration(ctx))})
	   	got = slh(ctx, NewMsgUnjail(addr))
	   	require.True(t, got.IsOK())

	   	// end block
	   	staking.EndBlocker(ctx, sk)

	   	// validator should be rebonded now
	   	validator, _ = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	   	require.Equal(t, sdk.Bonded, validator.GetStatus())

	   	// validator should have been slashed
	   	bondPool = sk.GetBondedPool(ctx)
	   	require.Equal(t, amt.Int64()-slashAmt, bondPool.GetCoins().AmountOf(sk.BondDenom(ctx)).Int64())

	   	// Validator start height should not have been changed
	   	info, found = keeper.getValidatorSigningInfo(ctx, sdk.ConsAddress(val.Address()))
	   	require.True(t, found)
	   	require.Equal(t, int64(0), info.StartHeight)
	   	// we've missed 2 blocks more than the maximum, so the counter was reset to 0 at 1 block more and is now 1
	   	require.Equal(t, int64(1), info.MissedBlocksCounter)

	   	// validator should not be immediately jailed again
	   	height++
	   	ctx = ctx.WithBlockHeight(height)
	   	keeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	   	validator, _ = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	   	require.Equal(t, sdk.Bonded, validator.GetStatus())

	   	// 500 signed blocks
	   	nextHeight := height + keeper.MinSignedPerWindow(ctx) + 1
	   	for ; height < nextHeight; height++ {
	   		ctx = ctx.WithBlockHeight(height)
	   		keeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	   	}

	   	// end block
	   	staking.EndBlocker(ctx, sk)

	   	// validator should be jailed again after 500 unsigned blocks
	   	nextHeight = height + keeper.MinSignedPerWindow(ctx) + 1
	   	for ; height <= nextHeight; height++ {
	   		ctx = ctx.WithBlockHeight(height)
	   		keeper.HandleValidatorSignature(ctx, val.Address(), power, false)
	   	}

	   	// end block
	   	staking.EndBlocker(ctx, sk)

	   	validator, _ = sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(val))
	   	require.Equal(t, sdk.Unbonding, validator.GetStatus()) */
}
