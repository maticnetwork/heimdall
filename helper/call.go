package helper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/staking"
	"math/big"
)

//
//func GetValidators() (validators []abci.ValidatorUpdate) {
//	stakeManagerInstance, err := GetStakeManagerInstance()
//	if err != nil {
//		Logger.Error("Error creating validatorSetInstance", "error", err)
//	}
//
//	powers, _, err := stakeManagerInstance.GetValidatorSet(nil)
//	if err != nil {
//		Logger.Error("Error getting validator set", "error", err)
//	}
//
//	for index := range powers {
//		pubkey, err := stakeManagerInstance.GetPubkey(nil, big.NewInt(int64(index)))
//		if err != nil {
//			Logger.Error("Error getting pubkey for index", "error", err)
//		}
//
//		var pubkeyBytes secp256k1.PubKeySecp256k1
//		_pubkey, _ := hex.DecodeString(pubkey)
//		copy(pubkeyBytes[:], _pubkey)
//		// todo use new valiator update here
//		validator := abci.ValidatorUpdate{
//			Power:  powers[index].Int64(),
//			PubKey: tmtypes.TM2PB.PubKey(pubkeyBytes),
//		}
//
//		validators = append(validators, validator)
//	}
//
//	return validators
//}

func GetHeaderInfo(headerId uint64) (root common.Hash, start uint64, end uint64, err error) {
	// get rootchain instance
	rootChainInstance, err := GetRootChainInstance()
	if err != nil {
		Logger.Error("Error creating rootchain instance while fetching headerBlock", "error", err, "headerBlockIndex", headerId)
		return common.HexToHash(""), 0, 0, err
	}

	// get header from rootchain
	headerBlock, err := rootChainInstance.HeaderBlock(nil, big.NewInt(int64(headerId)))
	if err != nil {
		Logger.Error("Unable to fetch header block from rootchain", "headerBlockIndex", headerId)
	}

	return headerBlock.Root, headerBlock.Start.Uint64(), headerBlock.End.Uint64(), nil
}

func GetValidatorInfo(addr common.Address) (validator staking.Validator, err error) {
	// get stakemanager intance
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating stakeManagerInstance while getting validator info", "error", err, "validatorAddrress", addr)
		return
	}

	amount, startEpoch, endEpoch, signer, _, err := stakeManagerInstance.GetStakerDetails(nil, addr)
	if err != nil {
		Logger.Error("Error fetching validator information from stakemanager", "Error", err, "ValidatorAddress", addr)
		return
	}
	validator = staking.Validator{
		ValidatorAddress: addr,
		Power:            amount.Int64(),
		StartEpoch:       startEpoch.Int64(),
		EndEpoch:         endEpoch.Int64(),
		Signer:           signer,
	}

	return validator, nil
}
