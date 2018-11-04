package helper

import (
	"encoding/hex"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
)

func GetValidators() (validators []abci.ValidatorUpdate) {
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating validatorSetInstance", "error", err)
	}

	powers, _, err := stakeManagerInstance.GetValidatorSet(nil)
	if err != nil {
		Logger.Error("Error getting validator set", "error", err)
	}

	for index := range powers {
		pubkey, err := stakeManagerInstance.GetPubkey(nil, big.NewInt(int64(index)))
		if err != nil {
			Logger.Error("Error getting pubkey for index", "error", err)
		}

		var pubkeyBytes secp256k1.PubKeySecp256k1
		_pubkey, _ := hex.DecodeString(pubkey)
		copy(pubkeyBytes[:], _pubkey)
		// todo use new valiator update here
		validator := abci.ValidatorUpdate{
			Power:  powers[index].Int64(),
			PubKey: tmtypes.TM2PB.PubKey(pubkeyBytes),
		}

		validators = append(validators, validator)
	}

	return validators
}

func GetLastBlock() uint64 {
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating validatorSetInstance", "error", err)
		return 0
	}

	lastBlock, err := stakeManagerInstance.StartBlock(nil)
	if err != nil {
		Logger.Error("Unable to fetch last block from mainchain", "error", err)
		return 0
	}

	return lastBlock.Uint64()
}
