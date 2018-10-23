package helper

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/contracts/stakemanager"
)

func GetValidators() (validators []abci.Validator) {
	validatorSetInstance := GetValidatorSetInstance(KovanClient)
	powers, ValidatorAddrs, err := validatorSetInstance.GetValidatorSet(nil)
	if err != nil {
		fmt.Printf(" The error is %v", err)
	}

	for index := range powers {
		pubkey, error := validatorSetInstance.GetPubkey(nil, big.NewInt(int64(index)))
		if error != nil {
			fmt.Errorf(" Error getting pubkey for index %v", error)
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
		fmt.Printf("New Validator is %v", validator)
		validators = append(validators, validator)
		//validatorPubKeys[index] = pubkey
	}

	return validators
}

func GetProposer() common.Address {
	validatorSetInstance := GetValidatorSetInstance(KovanClient)
	fmt.Printf("kovan client %v", KovanClient)
	currentProposer, err := validatorSetInstance.Proposer(nil)
	if err != nil {
		fmt.Printf("error getting proposer")
	}
	return currentProposer
}

// SubmitProof submit header
func SubmitProof(voteSignBytes []byte, sigs []byte, extradata []byte, start uint64, end uint64, rootHash common.Hash) {
	fmt.Printf("Root hash obtained for blocks from %v to %v is %v", start, end, rootHash)
	//auth := GenerateAuthObj(clientKovan)
	//auth.Value = big.NewInt(0)
	//todo change this to tx , right now its a call

	validatorSetInstance := GetValidatorSetInstance(KovanClient)
	fmt.Printf("validatorset instance %v", validatorSetInstance)
	fmt.Printf("inputs , vote: %v , sigs: %v , extradata %v ", hex.EncodeToString(voteSignBytes), hex.EncodeToString(sigs), hex.EncodeToString(extradata))
	res, proposer, error := validatorSetInstance.Validate(nil, voteSignBytes, sigs, extradata)
	if error != nil {
		fmt.Errorf("Er	ror hua ")
	}
	fmt.Printf("Submitted Proof Successfully %v %v %v ", res, proposer.String(), error)
}

func getValidatorByIndex(_index int64) abci.Validator {
	stakeManagerInstance, err := stakemanager.NewContracts(common.HexToAddress(GetConfig().StakeManagerAddress), KovanClient)
	if err != nil {
		log.Fatal(err)
	}
	validator, _ := stakeManagerInstance.Validators(nil, big.NewInt(_index))
	var _pubkey secp256k1.PubKeySecp256k1
	_pub, _ := hex.DecodeString(validator.Pubkey)
	copy(_pubkey[:], _pub[:])
	_address, _ := hex.DecodeString(_pubkey.Address().String())

	abciValidator := abci.Validator{
		Address: _address,
		Power:   validator.Power.Int64(),
		PubKey:  tmtypes.TM2PB.PubKey(_pubkey),
	}
	return abciValidator
}
