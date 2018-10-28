package helper

import (
	"encoding/hex"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GetValidators() (validators []abci.Validator) {
	stakeManagerInstance, err := GetStakeManagerInstance()

	ValidatorAddrs, err := stakeManagerInstance.GetCurrentValidatorSet(nil)
	if err != nil {
		Logger.Info("Error getting validator set", "Error", err)
	}

	for index := range ValidatorAddrs {
		if ValidatorAddrs[index].String() != "" {
			validatorStruct, error := stakeManagerInstance.Stakers(nil, ValidatorAddrs[index])
			if error != nil {
				Logger.Error("Error Fetching Staker", "Error", error, "Index", index)
			}
			pubkey := validatorStruct.Pubkey
			var pubkeyBytes secp256k1.PubKeySecp256k1
			_pubkey, _ := hex.DecodeString(pubkey)
			copy(pubkeyBytes[:], _pubkey)

			// todo add a check to check pubkey corresponds to address
			validator := abci.Validator{
				Address: ValidatorAddrs[index].Bytes(),
				Power:   validatorStruct.Amount.Int64(),
				PubKey:  tmtypes.TM2PB.PubKey(pubkeyBytes),
			}

			Logger.Info("New Validator Generated", "Validator", validator.String())

			validators = append(validators, validator)
		} else {
			Logger.Info("Validator Empty", "Index", index)
		}
	}

	return validators
}

// SubmitProof submit header
//func SubmitProof(voteSignBytes []byte, sigs []byte, extradata []byte, start uint64, end uint64, rootHash common.Hash) {
//	Logger.Info("Root Hash Generated ", "Start", start, "End", end, "RootHash", rootHash)
//	// get validator set instance from config
//	validatorSetInstance, err := GetValidatorSetInstance()
//
//	Logger.Info("Inputs to submitProof", " Vote", hex.EncodeToString(voteSignBytes), "Signatures", hex.EncodeToString(sigs), "Tx_Data", hex.EncodeToString(extradata))
//	// submit proof
//	result, proposer, err := validatorSetInstance.Validate(nil, voteSignBytes, sigs, extradata)
//	if err != nil {
//		Logger.Error("Checkpoint Submission Errored", "Error", err)
//	} else {
//		Logger.Info("Submitted Proof Successfully ", "Status", result, "Proposer", proposer)
//	}
//}
