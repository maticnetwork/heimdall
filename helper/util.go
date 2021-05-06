package helper

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"

	"github.com/cosmos/cosmos-sdk/client/input"

	"github.com/spf13/pflag"

	ethTypes "github.com/maticnetwork/bor/core/types"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	authsign "github.com/cosmos/cosmos-sdk/x/auth/signing"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/common"
	ethcrypto "github.com/maticnetwork/bor/crypto"
	secp256k1Crypto "github.com/maticnetwork/bor/crypto/secp256k1"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/tendermint/tendermint/crypto"
)

// ZeroHash represents empty hash
var ZeroHash = common.Hash{}

// ZeroAddress represents empty address
var ZeroAddress = common.Address{}

// ZeroPubKey represents empty pub key
var ZeroPubKey = hmCommonTypes.PubKey{}

const (
	COMPRESSED_PUBKEY_SIZE               = 32
	COMPRESSED_PUBKEY_SIZE_WITH_PREFIX   = 33
	UNCOMPRESSED_PUBKEY_SIZE             = 64
	UNCOMPRESSED_PUBKEY_SIZE_WITH_PREFIX = 65
)

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

// GetHeimdallServerEndpoint returns heimdall server endpoint
func GetHeimdallServerEndpoint(endpoint string) string {
	u, _ := url.Parse(GetConfig().HeimdallServerURL)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

// FetchFromAPI fetches data from any URL
func FetchFromAPI(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// response
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, err
	}

	Logger.Debug("Error while fetching data from URL", "status", resp.StatusCode, "URL", URL)
	return nil, fmt.Errorf("error while fetching data from url: %v, status: %v", URL, resp.StatusCode)
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
	if len(uncompressedBytes) == UNCOMPRESSED_PUBKEY_SIZE {
		uncompressedBytes = AppendPubkeyPrefix(uncompressedBytes)
	}
	uncompressed, err := ethcrypto.UnmarshalPubkey(uncompressedBytes)
	if err != nil {
		return nil, err
	}
	return ethcrypto.CompressPubkey(uncompressed), nil
}

// GetUpdatedValidators updates validators in validator set
func GetUpdatedValidators(
	currentSet *hmTypes.ValidatorSet,
	validators []*hmTypes.Validator,
	ackCount uint64,
) []*hmTypes.Validator {
	updates := make([]*hmTypes.Validator, 0)
	for _, v := range validators {
		// create copy of validator
		validator := v.Copy()

		address := validator.GetSigner()
		_, val := currentSet.GetByAddress(address)
		if val != nil && !validator.IsCurrentValidator(ackCount) {
			// remove validator
			validator.VotingPower = 0
			updates = append(updates, validator)
		} else if val == nil && validator.IsCurrentValidator(ackCount) {
			// add validator
			updates = append(updates, validator)
		} else if val != nil && validator.VotingPower != val.VotingPower {
			updates = append(updates, validator)
		}
	}

	return updates
}

// ToBytes32 is a convenience method for converting a byte slice to a fix
// sized 32 byte array. This method will truncate the input if it is larger
// than 32 bytes.
func ToBytes32(x []byte) [32]byte {
	var y [32]byte
	copy(y[:], x)
	return y
}

// GetTxEncoder returns tx encoder
//func GetTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
//	return legacytx.DefaultTxEncoder(cdc)
//}
//
//// GetTxDecoder returns tx decoder
//func GetTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
//	return legacytx.DefaultTxDecoder(cdc)
//}

func GenerateOrBroadcastTxCli(clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) error {
	txf := tx.NewFactoryCLI(clientCtx, flagSet)
	if clientCtx.GenerateOnly {
		return tx.GenerateTx(clientCtx, txf, msgs...)
	}

	return BuildAndBroadcastMsgsWithCli(clientCtx, txf, msgs...)
}

func BuildAndBroadcastMsgsWithCli(cliCtx client.Context, txf tx.Factory, msgs ...sdk.Msg) error {
	txBytes, err := GetSignedTxBytes(cliCtx, txf, msgs)
	if err != nil {
		return err
	}

	// just simulate
	if cliCtx.Simulate {
		fmt.Println("TxBytes", "0x"+hex.EncodeToString(txBytes))
		return nil
	}

	// broadcast to a Tendermint node
	resp, err := BroadcastTxBytes(cliCtx, txBytes, "")
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(resp)
}

// BuildAndBroadcastMsgs creates transaction and broadcasts it
func BuildAndBroadcastMsgs(cliCtx client.Context, txFactory tx.Factory, msgs []sdk.Msg) (res *sdk.TxResponse, err error) {
	txBytes, err := GetSignedTxBytes(cliCtx, txFactory, msgs)
	if err != nil {
		return &sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	return BroadcastTxBytes(cliCtx, txBytes, "")
}

// Get Sing
func GetSignedTxBytes(cliCtx client.Context, txf tx.Factory, msgs []sdk.Msg) ([]byte, error) {
	txf, err := PrepareTxBuilderFactory(cliCtx, txf)
	if err != nil {
		return nil, err
	}

	if cliCtx.Simulate {
		return nil, nil
	}

	txBuilder, err := tx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, err
	}

	if !cliCtx.SkipConfirm {
		out, err := cliCtx.TxConfig.TxJSONEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	// checking the fromName exists in CliCtx , then sign with msg with keyring
	fromName := cliCtx.GetFromName()
	if fromName != "" {
		err = tx.Sign(txf, cliCtx.GetFromName(), txBuilder)
		if err != nil {
			return nil, err
		}

		txBytes, err := cliCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}

		return txBytes, nil
	}

	// if there is no fromName in cliCtx we will sign the msg with private key
	signMode := txf.SignMode()
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   GetPubKeyForCosmos(),
		Data:     &sigData,
		Sequence: txf.Sequence(),
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	signData := authsign.SignerData{
		ChainID:       txf.ChainID(),
		AccountNumber: txf.AccountNumber(),
		Sequence:      txf.Sequence(),
	}

	sigV2Bytes, err := SignWithPrivKey(txf.SignMode(), signData, txBuilder, cliCtx.TxConfig, txf.Sequence())
	if err != nil {
		return nil, err
	}

	err = txBuilder.SetSignatures(sigV2Bytes)
	if err != nil {
		return nil, err
	}

	txBytes, err := cliCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

// Sign the msg with private key
func Sign(msg []byte) ([]byte, error) {
	return ethcrypto.Sign(ethcrypto.Keccak256Hash(msg).Bytes(), GetPrivKey().ToECDSA())
}

// SignWithPrivKey signs a given tx with the given private key, and returns the
// corresponding SignatureV2 if the signing is successful.
func SignWithPrivKey(
	signMode signing.SignMode, signerData authsigning.SignerData,
	txBuilder client.TxBuilder, txConfig client.TxConfig,
	accSeq uint64,
) (signing.SignatureV2, error) {
	var sigV2 signing.SignatureV2

	// Generate the bytes to be signed.
	signBytes, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		return sigV2, err
	}

	// Sign those bytes
	signature, err := Sign(signBytes)
	if err != nil {
		return sigV2, err
	}

	// Construct the SignatureV2 struct
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: signature,
	}

	sigV2 = signing.SignatureV2{
		PubKey:   GetPubKeyForCosmos(),
		Data:     &sigData,
		Sequence: accSeq,
	}

	return sigV2, nil
}

// BroadcastTxBytes sends request to tendermint using CLI
func BroadcastTxBytes(cliCtx client.Context, txBytes []byte, mode string) (res *sdk.TxResponse, err error) {
	Logger.Debug("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes), "txHash", hex.EncodeToString(tmTypes.Tx(txBytes).Hash()))
	if mode != "" {
		cliCtx.BroadcastMode = mode
	}
	return cliCtx.BroadcastTx(txBytes)
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilderFactory(cliCtx client.Context, txf tx.Factory) (tx.Factory, error) {
	from := cliCtx.GetFromAddress()
	if len(from[:]) == 0 {
		from = GetAddress()
	}

	if err := cliCtx.AccountRetriever.EnsureExists(cliCtx, from); err != nil {
		return txf, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := cliCtx.AccountRetriever.GetAccountNumberSequence(cliCtx, from)
		if err != nil {
			return txf, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	txf = txf.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
	return txf, nil
}

type sideTxSig struct {
	Address []byte
	Sig     []byte
}

// RecoverPubkey builds a StdSignature for given a StdSignMsg.
func recoverPubkey(msg []byte, sig []byte) ([]byte, error) {
	data := ethcrypto.Keccak256(msg)
	pubkey, err := secp256k1Crypto.RecoverPubkey(data, sig[:])
	if err != nil {
		return nil, err
	}
	ecdsa, err := ethcrypto.UnmarshalPubkey(pubkey)
	if err != nil {
		return nil, err
	}
	addr := tmcrypto.Address(ethcrypto.PubkeyToAddress(*ecdsa).Bytes())
	return addr, nil
}

// GetSideTxSigs returns sigs bytes from vote by tx hash
func GetSideTxSigs(txHash []byte, sideTxData []byte, unFilteredVotes []tmTypes.CommitSig) (sigs [][3]*big.Int, err error) {
	// TODO - Uncomment this if condition
	// side tx result with data
	// sideTxResultWithData := tmproto.SideTxResultWithData{
	// 	Result: &tmproto.SideTxResult{
	// 		TxHash: txHash,
	// 		Result: tmproto.SideTxResultType_YES,
	// 	},
	// 	Data: sideTxData,
	// }

	// TODO - Uncomment this if condition
	// draft signed data
	// signedData := sideTxResultWithData.GetData()

	sideTxSigs := make([]*sideTxSig, 0)
	for _, vote := range unFilteredVotes {
		// iterate through all side-tx results
		for _, sideTxResult := range vote.SideTxResults {
			// find side-tx result by tx-hash

			if bytes.Equal(sideTxResult.TxHash, txHash) &&
				len(sideTxResult.Sig) == 65 &&
				sideTxResult.Result == tmproto.SideTxResultType_YES {
				// validate sig
				// TODO - Uncomment this if condition
				// p, err := recoverPubkey(signedData, sideTxResult.Sig)
				if err == nil {
					// TODO - Uncomment this if condition
					// if it has valid sig, add it into side-tx sig array
					//if bytes.Equal(vote.ValidatorAddress.Bytes(), p) {
					sideTxSigs = append(sideTxSigs, &sideTxSig{
						Address: vote.ValidatorAddress.Bytes(),
						Sig:     sideTxResult.Sig,
					})
					//}
				}
			}
			// break
		}
	}

	if len(sideTxSigs) > 0 {
		// sort sigs by address
		sort.Slice(sideTxSigs, func(i, j int) bool {
			return bytes.Compare(sideTxSigs[i].Address, sideTxSigs[j].Address) < 0
		})

		// loop votes and append to sig to sigs
		for _, sideTxSig := range sideTxSigs {
			R, S, V, err := ethTypes.HomesteadSigner{}.SignatureValues(nil, sideTxSig.Sig)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, [3]*big.Int{R, S, V})
		}
	}

	return sigs, nil
}

// ValidateAndCompressPubKey validate and compress the pubkey
func ValidateAndCompressPubKey(pubkeyBytes []byte) ([]byte, error) {
	if len(pubkeyBytes) == UNCOMPRESSED_PUBKEY_SIZE {
		pubkeyBytes = AppendPubkeyPrefix(pubkeyBytes)
	}

	// check if key is uncompressed
	if len(pubkeyBytes) == UNCOMPRESSED_PUBKEY_SIZE_WITH_PREFIX {
		var err error
		pubkeyBytes, err = CompressPubKey(pubkeyBytes)
		if err != nil {
			return nil, fmt.Errorf("Invalid uncompressed pubkey %s", err)
		}
	}

	if len(pubkeyBytes) != COMPRESSED_PUBKEY_SIZE_WITH_PREFIX {
		return nil, fmt.Errorf("Invalid compressed pubkey")
	}

	return pubkeyBytes, nil
}
