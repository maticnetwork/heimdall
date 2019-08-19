package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/bits"
	"sort"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
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
	if mode != "" {
		cliCtx.BroadcastMode = mode
	}
	Logger.Info("Broadcasting tx bytes to tendermint", "txBytes", hex.EncodeToString(txBytes), "mode", cliCtx.BroadcastMode, "txHash", hex.EncodeToString(tmhash.Sum(txBytes[4:])))
	return cliCtx.BroadcastTx(txBytes)
}

// GetSigs returns sigs bytes from vote
func GetSigs(votes []*tmTypes.CommitSig) (sigs []byte) {
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
func GetVoteBytes(votes []*tmTypes.CommitSig, chainID string) []byte {
	vote := votes[0]
	v := tmTypes.Vote(*vote)
	// sign bytes for vote
	return v.SignBytes(chainID)
}

// CreateAndSendTx creates message and sends tx
// Used from cli- waits till transaction is included in block
func CreateAndSendTx(msg sdk.Msg, cliCtx context.CLIContext) (resp sdk.TxResponse, err error) {
	txBytes, err := CreateTxBytes(msg)
	if err != nil {
		return resp, err
	}
	resp, err = SendTendermintRequest(cliCtx, txBytes, BroadcastBlock)
	if err != nil {
		return resp, err
	}

	fmt.Printf("Transaction sent %v", resp.TxHash)
	return resp, nil
}

// TendermintTxDecode decodes transaction string and return base tx object
func TendermintTxDecode(txString string) ([]byte, error) {
	decodedTx, err := base64.StdEncoding.DecodeString(txString)
	if err != nil {
		return nil, err
	}

	return []byte(decodedTx), nil
}

// GetMerkleProofList return proof array
// each proof has one byte for direction: 0x0 for left and 0x1 for right
func GetMerkleProofList(proof *merkle.SimpleProof) [][]byte {
	result := [][]byte{}
	computeHashFromAunts(proof.Index, proof.Total, proof.LeafHash, proof.Aunts, &result)
	return result
}

// Use the leafHash and innerHashes to get the root merkle hash.
// If the length of the innerHashes slice isn't exactly correct, the result is nil.
// Recursive impl.
func computeHashFromAunts(index int, total int, leafHash []byte, innerHashes [][]byte, newInnerHashes *[][]byte) []byte {
	if index >= total || index < 0 || total <= 0 {
		return nil
	}
	switch total {
	case 0:
		panic("Cannot call computeHashFromAunts() with 0 total")
	case 1:
		if len(innerHashes) != 0 {
			return nil
		}
		return leafHash
	default:
		if len(innerHashes) == 0 {
			return nil
		}
		numLeft := getSplitPoint(total)
		if index < numLeft {
			leftHash := computeHashFromAunts(index, numLeft, leafHash, innerHashes[:len(innerHashes)-1], newInnerHashes)
			if leftHash == nil {
				return nil
			}
			*newInnerHashes = append(*newInnerHashes, append(rightPrefix, innerHashes[len(innerHashes)-1]...))
			return innerHash(leftHash, innerHashes[len(innerHashes)-1])
		}
		rightHash := computeHashFromAunts(index-numLeft, total-numLeft, leafHash, innerHashes[:len(innerHashes)-1], newInnerHashes)
		if rightHash == nil {
			return nil
		}
		*newInnerHashes = append(*newInnerHashes, append(leftPrefix, innerHashes[len(innerHashes)-1]...))
		return innerHash(innerHashes[len(innerHashes)-1], rightHash)
	}
}

//
// Inner funcitons
//

// getSplitPoint returns the largest power of 2 less than length
func getSplitPoint(length int) int {
	if length < 1 {
		panic("Trying to split a tree with size < 1")
	}
	uLength := uint(length)
	bitlen := bits.Len(uLength)
	k := 1 << uint(bitlen-1)
	if k == length {
		k >>= 1
	}
	return k
}

// TODO: make these have a large predefined capacity
var (
	leafPrefix  = []byte{0}
	innerPrefix = []byte{1}

	leftPrefix  = []byte{0}
	rightPrefix = []byte{1}
)

// returns tmhash(0x00 || leaf)
func leafHash(leaf []byte) []byte {
	return tmhash.Sum(append(leafPrefix, leaf...))
}

// returns tmhash(0x01 || left || right)
func innerHash(left []byte, right []byte) []byte {
	return tmhash.Sum(append(innerPrefix, append(left, right...)...))
}
