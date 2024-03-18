package simulation

import (
	"crypto/rand"
	"math/big"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// GenRandomVal generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []types.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := types.NewPubKey(util.AppendPrefix(privKey1.PubKey().Bytes()))

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
