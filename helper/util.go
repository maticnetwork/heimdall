package helper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmTypes "github.com/tendermint/tendermint/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ZeroHash represents empty hash
var ZeroHash = common.Hash{}

// ZeroAddress represents empty address
var ZeroAddress = common.Address{}

// ZeroPubKey represents empty pub key
var ZeroPubKey = hmTypes.PubKey{}

// UpdateValidators updates validators in validator set
func UpdateValidators(
	currentSet *hmTypes.ValidatorSet,
	validators []*hmTypes.Validator,
	validatorToSigner map[string]common.Address,
	ackCount uint64,
) error {
	var filteredValidators []*hmTypes.Validator
	for _, v := range validators {
		key := hex.EncodeToString(v.Address.Bytes())
		s, exists := validatorToSigner[key]
		if exists && bytes.Equal(v.Signer.Bytes(), s.Bytes()) {
			filteredValidators = append(filteredValidators, v)
		}
	}

	for _, validator := range filteredValidators {
		address := validator.Address.Bytes()
		_, val := currentSet.GetByAddress(address)
		if !validator.IsCurrentValidator(ackCount) {
			// remove val
			_, removed := currentSet.Remove(address)
			if !removed {
				return fmt.Errorf("Failed to remove validator %X", address)
			}
		} else if val == nil {
			// add val
			added := currentSet.Add(validator)
			if !added {
				return fmt.Errorf("Failed to add new validator %v", validator)
			}
		} else {
			// update val
			updated := currentSet.Update(validator)
			if !updated {
				return fmt.Errorf("Failed to update validator %X to %v", address, validator)
			}
		}
	}

	return nil
}

// GetPkObjects from crypto priv key
func GetPkObjects(privKey crypto.PrivKey) (secp256k1.PrivKeySecp256k1, secp256k1.PubKeySecp256k1) {
	var privObject secp256k1.PrivKeySecp256k1
	var pubObject secp256k1.PubKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)

	return privObject, pubObject
}

// StringToPubkey converts string to Pubkey
func StringToPubkey(pubkeyStr string) (crypto.PubKey, error) {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	_pubkey, err := hex.DecodeString(pubkeyStr)
	if err != nil {
		Logger.Error("Decoding of pubkey(string) to pubkey failed", "Error", err, "PubkeyString", pubkeyStr)
		return nil, err
	}
	// copy
	copy(pubkeyBytes[:], _pubkey)

	return pubkeyBytes, nil
}

// BytesToPubkey converts bytes to Pubkey
func BytesToPubkey(pubKey []byte) crypto.PubKey {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	copy(pubkeyBytes[:], pubKey)
	return pubkeyBytes
}

// CreateTxBytes creates tx bytes from Msg
func CreateTxBytes(msg sdk.Msg) ([]byte, error) {
	// tx := hmTypes.NewBaseTx(msg)
	pulp := hmTypes.GetPulpInstance()
	txBytes, err := pulp.EncodeToBytes(msg)
	if err != nil {
		Logger.Error("Error generating TX Bytes", "error", err)
		return []byte(""), err
	}
	return txBytes, nil
}

// SendTendermintRequest sends request to tendermint
func SendTendermintRequest(cliCtx context.CLIContext, txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	Logger.Info("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes))
	return cliCtx.BroadcastTx(txBytes)
}

// GetSigs returns sigs bytes from vote
func GetSigs(votes []tmTypes.Vote) (sigs []byte) {
	sort.Slice(votes, func(i, j int) bool {
		return bytes.Compare(votes[i].ValidatorAddress.Bytes(), votes[j].ValidatorAddress.Bytes()) < 0
	})

	// loop votes and append to sig to sigs
	for _, vote := range votes {
		sigs = append(sigs, vote.Signature...)
	}
	return
}

// GetVoteBytes returns vote bytes
func GetVoteBytes(votes []tmTypes.Vote, ctx sdk.Context) []byte {
	// sign bytes for vote
	return votes[0].SignBytes(ctx.ChainID())
}

func RandomAddress() {

}
