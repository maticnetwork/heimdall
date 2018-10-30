package helper

import (
	"encoding/hex"

	"math/big"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

//func GetValidators() (validators []abci.Validator) {
//	stakeManagerInstance, err := GetStakeManagerInstance()
//	if err!=nil{
//		Logger.Error("Error creating stakeManagerInstance","Error",err)
//	}
//
//	ValidatorAddrs, err := stakeManagerInstance.GetCurrentValidatorSet(nil)
//	if err != nil {
//		Logger.Info("Error getting validator set", "Error", err)
//	}
//
//	for index := range ValidatorAddrs {
//		if ValidatorAddrs[index].String() != "" {
//			validatorStruct, error := stakeManagerInstance.Stakers(nil, ValidatorAddrs[index])
//			if error != nil {
//				Logger.Error("Error Fetching Staker", "Error", error, "Index", index)
//			}
//			pubkey := validatorStruct.Pubkey
//			var pubkeyBytes secp256k1.PubKeySecp256k1
//			_pubkey, _ := hex.DecodeString(pubkey)
//			copy(pubkeyBytes[:], _pubkey)
//
//			// todo add a check to check pubkey corresponds to address
//			validator := abci.Validator{
//				Address: ValidatorAddrs[index].Bytes(),
//				Power:   validatorStruct.Amount.Int64(),
//				PubKey:  tmtypes.TM2PB.PubKey(pubkeyBytes),
//			}
//
//
//			validators = append(validators, validator)
//		} else {
//			Logger.Info("Validator Empty", "Index", index)
//		}
//	}
//
//	return validators
//}

func GetValidatorsFromMock() (validators []abci.Validator) {
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating validatorSetInstance", "error", err)
	}

	powers, ValidatorAddrs, err := stakeManagerInstance.GetValidatorSet(nil)
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

		validator := abci.Validator{
			Address: ValidatorAddrs[index].Bytes(),
			Power:   powers[index].Int64(),
			PubKey:  tmtypes.TM2PB.PubKey(pubkeyBytes),
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
