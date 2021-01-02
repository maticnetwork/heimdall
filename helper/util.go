package helper

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/common"
	ethcrypto "github.com/maticnetwork/bor/crypto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
)

// ZeroHash represents empty hash
var ZeroHash = common.Hash{}

// ZeroAddress represents empty address
var ZeroAddress = common.Address{}

// ZeroPubKey represents empty pub key
var ZeroPubKey = hmCommonTypes.PubKey{}

// GetPowerFromAmount returns power from amount -- note that this will pollute amount object
func GetPowerFromAmount(amount *big.Int) (*big.Int, error) {
	decimals18 := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)
	if amount.Cmp(decimals18) == -1 {
		return nil, errors.New("amount must be more than 1 token")
	}

	return amount.Div(amount, decimals18), nil
}

// GetPubObjects returns PubKeySecp256k1 public key
func GetPubObjects(pubkey crypto.PubKey) secp256k1.PubKey {
	var pubObject secp256k1.PubKey
	cdc.MustUnmarshalBinaryBare(pubkey.Bytes(), &pubObject)
	return pubObject
}

// UnpackSigAndVotes Unpacks Sig and Votes from Tx Payload
func UnpackSigAndVotes(payload []byte, abi abi.ABI) (votes []byte, sigs []byte, checkpointData []byte, err error) {
	// recover Method from signature and ABI
	method := abi.Methods["submitHeaderBlock"]
	decodedPayload := payload[4:]
	inputDataMap := make(map[string]interface{})
	// unpack method inputs
	err = method.Inputs.UnpackIntoMap(inputDataMap, decodedPayload)
	if err != nil {
		return
	}
	sigs = inputDataMap["sigs"].([]byte)
	checkpointData = inputDataMap["txData"].([]byte)
	votes = inputDataMap["vote"].([]byte)
	return
}

// GetFromAddress get from address
func GetFromAddress(cliCtx client.Context) sdk.AccAddress {
	fromAddress := cliCtx.GetFromAddress()
	if !fromAddress.Empty() {
		return fromAddress
	}

	return GetAddress()
}

// EventByID looks up a event by the topic id
func EventByID(abiObject *abi.ABI, sigdata []byte) *abi.Event {
	for _, event := range abiObject.Events {
		if bytes.Equal(event.ID.Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}

// AppendPubkeyPrefix returns publickey in uncompressed format
func AppendPubkeyPrefix(signerPubKey []byte) []byte {
	// append prefix - "0x04" as heimdall uses publickey in uncompressed format. Refer below link
	// https://superuser.com/questions/1465455/what-is-the-size-of-public-key-for-ecdsa-spec256r1
	prefix := make([]byte, 1)
	prefix[0] = byte(0x04)
	signerPubKey = append(prefix[:], signerPubKey[:]...)
	return signerPubKey
}

// DecompressPubKey decompress pub key
func DecompressPubKey(compressed []byte) ([]byte, error) {
	ecdsaPubkey, err := ethcrypto.DecompressPubkey(compressed)
	if err != nil {
		return nil, err
	}
	return ethcrypto.FromECDSAPub(ecdsaPubkey), nil
}

// CompressPubKey decompress pub key
func CompressPubKey(uncompressedBytes []byte) ([]byte, error) {
	if len(uncompressedBytes) == 64 {
		uncompressedBytes = AppendPubkeyPrefix(uncompressedBytes)
	}
	uncompressed, err := ethcrypto.UnmarshalPubkey(uncompressedBytes)
	if err != nil {
		return nil, err
	}
	return ethcrypto.CompressPubkey(uncompressed), nil
}
