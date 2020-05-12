package test

import (
	"bytes"
	"math/big"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/maticnetwork/bor/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/types"
)

func MakeTestCodec() *codec.Codec {
	cdc := codec.New()

	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)

	bankTypes.RegisterCodec(cdc)

	// custom types
	borTypes.RegisterCodec(cdc)
	checkpointTypes.RegisterCodec(cdc)
	stakingTypes.RegisterCodec(cdc)

	cdc.Seal()
	return cdc
}

// NOTE commented this out as it causes a build error
// some of the functions are passed incorrect data. Should not effect any
// functionality
//
// // create random header block
// func GenRandCheckpointHeader(start int, headerSize int) (headerBlock types.CheckpointBlockHeader, err error) {
// 	start = start
// 	end := start + headerSize
// 	maxCheckpointLenght := uint64(1024)
// 	roothash, err := checkpointTypes.GetHeaders(uint64(start), uint64(end), maxCheckpointLenght)
// 	if err != nil {
// 		return headerBlock, err
// 	}
// 	proposer := ethcmn.Address{}
// 	headerBlock = types.CreateBlock(uint64(start), uint64(end), types.HexToHeimdallHash(hex.EncodeToString(roothash)), types.HexToHeimdallHash(hex.EncodeToString(roothash)), types.HexToHeimdallAddress(proposer.String()), uint64(time.Now().UTC().Unix()))
//
// 	return headerBlock, nil
// }

// Generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []types.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := types.NewPubKey(privKey1.PubKey().Bytes())
		if randomise {
			startBlock := uint64(rand.Intn(10))
			// todo find a way to genrate non zero random number
			if startBlock == 0 {
				startBlock = 1
			}
			power := uint64(rand.Intn(100))
			if power == 0 {
				power = 1
			}
		}

		newVal := types.Validator{
			ID:               types.NewValidatorID(startID + uint64(i)),
			StartEpoch:       startBlock,
			EndEpoch:         startBlock + timeAlive,
			VotingPower:      power,
			Signer:           types.HexToHeimdallAddress(pubkey.Address().String()),
			PubKey:           pubkey,
			ProposerPriority: 0,
		}
		validators = append(validators, newVal)
	}
	return
}

// GenRandomDividendAccount random dividend accounts
func GenRandomDividendAccount(count int, startID uint64, randomise bool) (dividendAccounts []types.DividendAccount) {
	for i := 0; i < count; i++ {

		if randomise {
			fee := big.NewInt(int64(rand.Intn(100))).String()
			slashedAmount := big.NewInt(int64(rand.Intn(100))).String()

			newAcc := types.DividendAccount{
				ID:            types.NewDividendAccountID(startID + uint64(i)),
				FeeAmount:     fee,
				SlashedAmount: slashedAmount,
			}
			dividendAccounts = append(dividendAccounts, newAcc)
		}
	}
	return
}

// Load Validator Set
func LoadValidatorSet(count int, t *testing.T, keeper staking.Keeper, ctx sdk.Context, randomise bool, timeAlive int) types.ValidatorSet {
	// create 4 validators
	validators := GenRandomVal(4, 0, 10, uint64(timeAlive), randomise, 1)
	var valSet types.ValidatorSet
	// add validators to new Validator set and state
	for _, validator := range validators {
		err := keeper.AddValidator(ctx, validator)
		require.Empty(t, err, "Unable to set validator, Error: %v", err)
		// add validator to validator set
		// valSet.Add(&validator)
		valSet.UpdateWithChangeSet([]*types.Validator{&validator})
	}

	err := keeper.UpdateValidatorSetInStore(ctx, valSet)
	require.Empty(t, err, "Unable to update validator set")
	return valSet
}

// finds address in give validator slice
func FindSigner(signer ethcmn.Address, vals []types.Validator) bool {
	for _, val := range vals {
		if bytes.Compare(signer.Bytes(), val.Signer.Bytes()) == 0 {
			return true
		}
	}
	return false
}

// print validators
func PrintVals(t *testing.T, valset *types.ValidatorSet) {
	t.Log("Printing validators")
	for i, val := range valset.Validators {
		t.Log("Validator ===> ", "Index", i, "ValidatorInfo", val.String())
	}
}
