package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// GenRandomVal generate random validators
func GenRandomVal(count int, startBlock uint64, power int64, timeAlive uint64, randomise bool, startID uint64) (validators []hmTypes.Validator) {
	for i := 0; i < count; i++ {
		privKey1 := secp256k1.GenPrivKey()
		pubkey := hmCommonTypes.NewPubKey(privKey1.PubKey().Bytes())
		signer := sdk.AccAddress(privKey1.PubKey().Address().Bytes())

		if randomise {
			rand.Seed(time.Now().UnixNano())
			startBlock = uint64(rand.Intn(9))
			startBlock += 1
			power = int64(rand.Intn(99))
			power += 1
		}

		newVal := hmTypes.Validator{
			ID:               hmTypes.NewValidatorID(startID + uint64(i)),
			StartEpoch:       startBlock,
			EndEpoch:         startBlock + timeAlive,
			VotingPower:      power,
			Signer:           signer.String(),
			PubKey:           pubkey.String(),
			ProposerPriority: 0,
		}
		validators = append(validators, newVal)
	}
	return
}
