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
	tmTypes "github.com/tendermint/tendermint/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
	ackCount uint64,
) error {
	for _, validator := range validators {
		address := validator.Signer.Bytes()
		_, val := currentSet.GetByAddress(address)
		if val != nil && !validator.IsCurrentValidator(ackCount) {
			// remove val
			_, removed := currentSet.Remove(address)
			if !removed {
				return fmt.Errorf("Failed to remove validator %X", address)
			}
		} else if val == nil && validator.IsCurrentValidator(ackCount) {
			// add val
			added := currentSet.Add(validator)
			if !added {
				return fmt.Errorf("Failed to add new validator %v", validator)
			}
		} else if val != nil {
			validator.Accum = val.Accum             // use last accum
			updated := currentSet.Update(validator) // update validator
			validator.Accum = 0                     // reset accum
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

func GetPubObjects(pubkey crypto.PubKey) secp256k1.PubKeySecp256k1 {
	var pubObject secp256k1.PubKeySecp256k1
	cdc.MustUnmarshalBinaryBare(pubkey.Bytes(), &pubObject)
	return pubObject
}

// StringToPubkey converts string to Pubkey
func StringToPubkey(pubkeyStr string) (secp256k1.PubKeySecp256k1, error) {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	_pubkey, err := hex.DecodeString(pubkeyStr)
	if err != nil {
		Logger.Error("Decoding of pubkey(string) to pubkey failed", "Error", err, "PubkeyString", pubkeyStr)
		return pubkeyBytes, err
	}
	// copy
	copy(pubkeyBytes[:], _pubkey)

	return pubkeyBytes, nil
}

// BytesToPubkey converts bytes to Pubkey
func BytesToPubkey(pubKey []byte) secp256k1.PubKeySecp256k1 {
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
func SendTendermintRequest(cliCtx context.CLIContext, txBytes []byte, mode string) (sdk.TxResponse, error) {
	Logger.Info("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes), "txHash", hex.EncodeToString(tmhash.Sum(txBytes[4:])))
	if mode != "" {
		cliCtx.BroadcastMode = mode
	}
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

// Creates message and sends tx
// Used from cli- waits till transaction is included in block
func CreateAndSendTx(msg sdk.Msg, cliCtx context.CLIContext) (err error) {
	txBytes, err := CreateTxBytes(msg)
	if err != nil {
		return err
	}

	resp, err := SendTendermintRequest(cliCtx, txBytes, BroadcastBlock)
	if err != nil {
		return err
	}

	fmt.Printf("Transaction sent %v", resp.TxHash)
	return nil
}
