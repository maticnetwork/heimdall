package bor

import (
	"crypto/rand"
	"math/big"
	"reflect"
	"testing"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestShuffleList(t *testing.T) {
	var list1 []uint64
	seed1 := [32]byte{1, 128, 12}
	seed2 := [32]byte{2, 128, 12}
	for i := 0; i < 10; i++ {
		list1 = append(list1, uint64(i))
	}

	list2 := make([]uint64, len(list1))
	copy(list2, list1)

	list1, err := ShuffleList(list1, seed1)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	list2, err = ShuffleList(list2, seed2)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	if reflect.DeepEqual(list1, list2) {
		t.Errorf("2 shuffled lists shouldn't be equal")
	}
	// if !reflect.DeepEqual(list1, []uint64{0, 7, 8, 6, 3, 9, 4, 5, 2, 1}) {
	// 	t.Errorf("list 1 was incorrectly shuffled got: %v", list1)
	// }
	// if !reflect.DeepEqual(list2, []uint64{0, 5, 2, 1, 6, 8, 7, 3, 4, 9}) {
	// 	t.Errorf("list 2 was incorrectly shuffled got: %v", list2)
	// }
}

func TestValShuffle(t *testing.T) {
	seedHash1 := common.HexToHash("0xc46afc66ad9f4b237414c23a0cf0c469aeb60f52176565990644a9ee36a17667")
	initialVals := GenRandomVal(50, 0, 100, uint64(10), true, 1)

	_, err := XXXSelectNextProducers(seedHash1, initialVals, 40)
	require.NoError(t, err)
}

// Generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []types.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := types.NewPubKey(privKey1.PubKey().Bytes())
		if randomise {
			startBlock = generateRandNumber(10)
			power = int64(generateRandNumber(100))
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

func generateRandNumber(max int64) uint64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 1
	}

	return nBig.Uint64()
}
