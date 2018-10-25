package helper

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GetValidators() (validators []abci.Validator) {
	Logger.With("module", "helper/call")
	validatorSetInstance := GetValidatorSetInstance(MainChainClient)
	powers, ValidatorAddrs, err := validatorSetInstance.GetValidatorSet(nil)
	if err != nil {
		Logger.Info("Error getting Validator Set ", "Error", err)
	}

	for index := range powers {
		pubkey, error := validatorSetInstance.GetPubkey(nil, big.NewInt(int64(index)))
		if error != nil {
			Logger.Error("Error getting pubkey ", "Error", error, "Index", index)
		}

		var pubkeyBytes secp256k1.PubKeySecp256k1
		_pubkey, _ := hex.DecodeString(pubkey)
		copy(pubkeyBytes[:], _pubkey)

		// todo add a check to check pubkey corresponds to address
		validator := abci.Validator{
			Address: ValidatorAddrs[index].Bytes(),
			Power:   powers[index].Int64(),
			PubKey:  tmtypes.TM2PB.PubKey(pubkeyBytes),
		}

		Logger.Info("New Validator Generated ", "Validator", validator)

		validators = append(validators, validator)
		//validatorPubKeys[index] = pubkey
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

	Logger.Info("Root Hash Generated ", "Start ", start, "End ", end, "RootHash ", rootHash)
	// get validator set instance from config
	validatorSetInstance := GetValidatorSetInstance(MainChainClient)

	Logger.Info("Inputs to submitProof", " Vote", hex.EncodeToString(voteSignBytes), "Signatures", hex.EncodeToString(sigs), "Tx Data ", hex.EncodeToString(extradata))
	// submit proof
	result, proposer, error := validatorSetInstance.Validate(nil, voteSignBytes, sigs, extradata)
	if error != nil {
		Logger.Error("Checkpoint Submission Errored : %v", error)
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
