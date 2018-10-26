package helper

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GetValidators() (validators []abci.Validator) {
	Logger.With("module", "helper/call")

	stakeManagerInstance := GetStakeManagerInstance()

	ValidatorAddrs, err := stakeManagerInstance.GetCurrentValidatorSet(nil)
	if err != nil {
		Logger.Info("Error getting Validator Set ", "Error", err)
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

			Logger.Info("New Validator Generated ", "Validator", validator)

			validators = append(validators, validator)
		} else {
			Logger.Info("Validator Empty", "Index", index)
		}

	}

	return validators
}

func GetProposer() common.Address {

	validatorSetInstance := GetValidatorSetInstance(MainChainClient)

	currentProposer, err := validatorSetInstance.Proposer(nil)
	if err != nil {
		Logger.Error("Unable to get proposer ", "Error", err)
	}

	return currentProposer
}

// SubmitProof submit header
func SubmitProof(voteSignBytes []byte, sigs []byte, extradata []byte, start uint64, end uint64, rootHash common.Hash) {

	Logger.Info("Root Hash Generated ", "Start", start, "End", end, "RootHash", rootHash)
	// get validator set instance from config
	validatorSetInstance := GetValidatorSetInstance(MainChainClient)

	Logger.Info("Inputs to submitProof", " Vote", hex.EncodeToString(voteSignBytes), "Signatures", hex.EncodeToString(sigs), "Tx_Data", hex.EncodeToString(extradata))
	// submit proof
	result, proposer, error := validatorSetInstance.Validate(nil, voteSignBytes, sigs, extradata)
	if error != nil {
		Logger.Error("Checkpoint Submission Errored", "Error", error)
	} else {
		Logger.Info("Submitted Proof Successfully ", "Status", result, "Proposer", proposer)
	}

}

// To be used later
//
//func getValidatorByIndex(_index int64) abci.Validator {
//
//	stakeManagerInstance, err := stakemanager.NewContracts(common.HexToAddress(GetConfig().StakeManagerAddress), KovanClient)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	validator, _ := stakeManagerInstance.Validators(nil, big.NewInt(_index))
//	var _pubkey secp256k1.PubKeySecp256k1
//	_pub, _ := hex.DecodeString(validator.Pubkey)
//	copy(_pubkey[:], _pub[:])
//	_address, _ := hex.DecodeString(_pubkey.Address().String())
//
//	abciValidator := abci.Validator{
//		Address: _address,
//		Power:   validator.Power.Int64(),
//		PubKey:  tmtypes.TM2PB.PubKey(_pubkey),
//	}
//	return abciValidator
//}
