package helper

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/bits"
	"os"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/common"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ZeroHash represents empty hash
var ZeroHash = common.Hash{}

// ZeroAddress represents empty address
var ZeroAddress = common.Address{}

// ZeroPubKey represents empty pub key
var ZeroPubKey = hmTypes.PubKey{}

// GetFromAddress get from address
func GetFromAddress(cliCtx context.CLIContext) types.HeimdallAddress {
	fromAddress := cliCtx.GetFromAddress()
	if !fromAddress.Empty() {
		return types.AccAddressToHeimdallAddress(fromAddress)
	}

	return types.BytesToHeimdallAddress(GetAddress())
}

// Paginate returns the correct starting and ending index for a paginated query,
// given that client provides a desired page and limit of objects and the handler
// provides the total number of objects. If the start page is invalid, non-positive
// values are returned signaling the request is invalid.
//
// NOTE: The start page is assumed to be 1-indexed.
func Paginate(numObjs, page, limit, defLimit int) (start, end int) {
	if page == 0 {
		// invalid start page
		return -1, -1
	} else if limit == 0 {
		limit = defLimit
	}

	start = (page - 1) * limit
	end = limit + start

	if end >= numObjs {
		end = numObjs
	}

	if start >= numObjs {
		// page is out of bounds
		return -1, -1
	}

	return start, end
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

		address := validator.Signer.Bytes()
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

// GetPkObjects from crypto priv key
func GetPkObjects(privKey crypto.PrivKey) (secp256k1.PrivKeySecp256k1, secp256k1.PubKeySecp256k1) {
	var privObject secp256k1.PrivKeySecp256k1
	var pubObject secp256k1.PubKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privKey.Bytes(), &privObject)
	cdc.MustUnmarshalBinaryBare(privObject.PubKey().Bytes(), &pubObject)
	return privObject, pubObject
}

// GetPubObjects returns PubKeySecp256k1 public key
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

// GetSigs returns sigs bytes from vote
func GetSigs(unFilteredVotes []*tmTypes.CommitSig) (sigs []byte) {
	votes := make([]*tmTypes.CommitSig, 0)
	for _, item := range unFilteredVotes {
		if item != nil {
			votes = append(votes, item)
		}
	}

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

// GetTxEncoder returns tx encoder
func GetTxEncoder() sdk.TxEncoder {
	return authTypes.RLPTxEncoder(authTypes.GetPulpInstance())
}

// GetTxDecoder returns tx decoder
func GetTxDecoder() sdk.TxDecoder {
	return authTypes.RLPTxDecoder(authTypes.GetPulpInstance())
}

// GetStdTxBytes get tx bytes
func GetStdTxBytes(cliCtx context.CLIContext, tx authTypes.StdTx) ([]byte, error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder())
	return txBldr.GetStdTxBytes(tx)
}

// BroadcastMsgs creates transaction and broadcasts it
func BroadcastMsgs(cliCtx context.CLIContext, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder())
	return BuildAndBroadcastMsgs(cliCtx, txBldr, msgs)
}

// BroadcastTx broadcasts transaction
func BroadcastTx(cliCtx context.CLIContext, tx authTypes.StdTx, mode string) (res sdk.TxResponse, err error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder())

	var txBytes []byte
	txBytes, err = txBldr.GetStdTxBytes(tx)
	if err == nil {
		res, err = BroadcastTxBytes(cliCtx, txBytes, mode)
	}

	return
}

// BroadcastMsgsWithCLI creates message and sends tx
// Used from cli- waits till transaction is included in block
func BroadcastMsgsWithCLI(cliCtx context.CLIContext, msgs []sdk.Msg) error {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder())

	if cliCtx.GenerateOnly {
		return PrintUnsignedStdTx(cliCtx, txBldr, msgs)
	}

	return BuildAndBroadcastMsgsWithCLI(cliCtx, txBldr, msgs)
}

// BuildAndBroadcastMsgs creates transaction and broadcasts it
func BuildAndBroadcastMsgs(cliCtx context.CLIContext, txBldr authTypes.TxBuilder, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBytes, err := GetSignedTxBytes(cliCtx, txBldr, msgs)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// broadcast to a Tendermint node
	return BroadcastTxBytes(cliCtx, txBytes, "")
}

// BuildAndBroadcastMsgsWithCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func BuildAndBroadcastMsgsWithCLI(cliCtx context.CLIContext, txBldr authTypes.TxBuilder, msgs []sdk.Msg) error {
	txBytes, err := GetSignedTxBytesWithCLI(cliCtx, txBldr, msgs)
	if err != nil {
		return err
	}

	// just simulate
	if cliCtx.Simulate {
		fmt.Println("TxBytes", "0x"+hex.EncodeToString(txBytes))
		return nil
	}

	// broadcast to a Tendermint node
	res, err := BroadcastTxBytes(cliCtx, txBytes, BroadcastSync) // wait until tx included in block
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(res)
}

// GetSignedTxBytes returns signed tx bytes
func GetSignedTxBytes(cliCtx context.CLIContext, txBldr authTypes.TxBuilder, msgs []sdk.Msg) ([]byte, error) {
	txBldr, err := PrepareTxBuilder(cliCtx, txBldr)
	if err != nil {
		return nil, err
	}

	fromName := cliCtx.GetFromName()
	if fromName == "" {
		return txBldr.BuildAndSign(GetPrivKey(), msgs)
	}

	if cliCtx.Simulate {
		return nil, nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return nil, err
		}

		var json []byte
		if viper.GetBool(client.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return nil, err
	}
	// build and sign the transaction
	return txBldr.BuildAndSignWithPassphrase(fromName, passphrase, msgs)
}

// GetSignedTxBytesWithCLI returns signed tx bytes
func GetSignedTxBytesWithCLI(cliCtx context.CLIContext, txBldr authTypes.TxBuilder, msgs []sdk.Msg) ([]byte, error) {
	txBldr, err := PrepareTxBuilder(cliCtx, txBldr)
	if err != nil {
		return nil, err
	}

	fromName := cliCtx.GetFromName()
	if fromName == "" {
		return txBldr.BuildAndSign(GetPrivKey(), msgs)
	}

	if cliCtx.Simulate {
		return nil, nil
	}

	if !cliCtx.SkipConfirm {
		stdSignMsg, err := txBldr.BuildSignMsg(msgs)
		if err != nil {
			return nil, err
		}

		var json []byte
		if viper.GetBool(client.FlagIndentResponse) {
			json, err = cliCtx.Codec.MarshalJSONIndent(stdSignMsg, "", "  ")
			if err != nil {
				panic(err)
			}
		} else {
			json = cliCtx.Codec.MustMarshalJSON(stdSignMsg)
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", json)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf)
		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	passphrase, err := keys.GetPassphrase(fromName)
	if err != nil {
		return nil, err
	}

	return txBldr.BuildAndSignWithPassphrase(fromName, passphrase, msgs)
}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilder(cliCtx context.CLIContext, txBldr authTypes.TxBuilder) (authTypes.TxBuilder, error) {
	from := cliCtx.GetFromAddress()
	if len(from[:]) == 0 {
		from = GetAddress()
	}

	// get heimdall address
	fhAddress := types.BytesToHeimdallAddress(from)

	accGetter := authTypes.NewAccountRetriever(cliCtx)
	if err := accGetter.EnsureExists(fhAddress); err != nil {
		return txBldr, err
	}

	txbldrAccNum, txbldrAccSeq := txBldr.AccountNumber(), txBldr.Sequence()
	// TODO: (ref #1903) Allow for user supplied account number without
	// automatically doing a manual lookup.
	if txbldrAccNum == 0 || txbldrAccSeq == 0 {
		num, seq, err := authTypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(fhAddress)
		if err != nil {
			return txBldr, err
		}

		if txbldrAccNum == 0 {
			txBldr = txBldr.WithAccountNumber(num)
		}
		if txbldrAccSeq == 0 {
			txBldr = txBldr.WithSequence(seq)
		}
	}

	return txBldr, nil
}

// PrintUnsignedStdTx builds an unsigned StdTx and prints it to os.Stdout.
func PrintUnsignedStdTx(cliCtx context.CLIContext, txBldr authTypes.TxBuilder, msgs []sdk.Msg) error {
	stdTx, err := buildUnsignedStdTxOffline(txBldr, cliCtx, msgs)
	if err != nil {
		return err
	}

	json, err := cliCtx.Codec.MarshalJSON(stdTx)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(cliCtx.Output, "%s\n", json)
	return nil
}

// SignStdTx appends a signature to a StdTx and returns a copy of it. If appendSig
// is false, it replaces the signatures already attached with the new signature.
// Don't perform online validation or lookups if offline is true.
func SignStdTx(
	cliCtx context.CLIContext, stdTx authTypes.StdTx, appendSig bool, offline bool,
) (authTypes.StdTx, error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder())

	var signedStdTx authTypes.StdTx

	fromName := cliCtx.GetFromName()
	var addr []byte
	if fromName == "" {
		addr = GetAddress()
	} else {
		info, err := txBldr.Keybase().Get(fromName)
		if err != nil {
			return signedStdTx, err
		}

		addr = info.GetPubKey().Address().Bytes()
	}

	if !offline {
		var err error
		txBldr, err = populateAccountFromState(txBldr, cliCtx, addr)
		if err != nil {
			return signedStdTx, err
		}
	}

	if fromName != "" {
		passphrase, err := keys.GetPassphrase(fromName)
		if err != nil {
			return signedStdTx, err
		}

		// with passpharse
		return txBldr.SignStdTxWithPassphrase(fromName, passphrase, stdTx, appendSig)
	}

	return txBldr.SignStdTx(GetPrivKey(), stdTx, appendSig)
}

// ReadStdTxFromFile and decode a StdTx from the given filename.  Can pass "-" to read from stdin.
func ReadStdTxFromFile(cdc *amino.Codec, filename string) (stdTx authTypes.StdTx, err error) {
	var bytes []byte
	if filename == "-" {
		bytes, err = ioutil.ReadAll(os.Stdin)
	} else {
		bytes, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return
	}

	if err = cdc.UnmarshalJSON(bytes, &stdTx); err != nil {
		return
	}
	return
}

// BroadcastTxBytes sends request to tendermint using CLI
func BroadcastTxBytes(cliCtx context.CLIContext, txBytes []byte, mode string) (sdk.TxResponse, error) {
	Logger.Debug("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes), "txHash", hex.EncodeToString(tmhash.Sum(txBytes[4:])))
	if mode != "" {
		cliCtx.BroadcastMode = mode
	}
	return cliCtx.BroadcastTx(txBytes)
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

// AppendBytes appends bytes
func AppendBytes(data ...[]byte) []byte {
	var result []byte
	for _, v := range data {
		result = append(result, v[:]...)
	}
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

func populateAccountFromState(txBldr authTypes.TxBuilder, cliCtx context.CLIContext, addr []byte) (authTypes.TxBuilder, error) {
	// get account getter
	accGetter := authTypes.NewAccountRetriever(cliCtx)

	// key
	key := hmTypes.BytesToHeimdallAddress(addr)

	// ensure account exists
	if err := accGetter.EnsureExists(key); err != nil {
		return txBldr, err
	}

	acc, err := accGetter.GetAccount(key)
	if err != nil {
		return txBldr, err
	}

	accNum := acc.GetAccountNumber()
	accSeq := acc.GetSequence()

	return txBldr.WithAccountNumber(accNum).WithSequence(accSeq), nil
}

func isTxSigner(user sdk.AccAddress, signers []sdk.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}

	return false
}

func buildUnsignedStdTxOffline(txBldr authTypes.TxBuilder, cliCtx context.CLIContext, msgs []sdk.Msg) (stdTx authTypes.StdTx, err error) {
	stdSignMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		return stdTx, err
	}

	return authTypes.NewStdTx(stdSignMsg.Msg, nil, stdSignMsg.Memo), nil
}

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

// ToBytes32 is a convenience method for converting a byte slice to a fix
// sized 32 byte array. This method will truncate the input if it is larger
// than 32 bytes.
func ToBytes32(x []byte) [32]byte {
	var y [32]byte
	copy(y[:], x)
	return y
}

// GetReceiptLogData get receipt log data
func GetReceiptLogData(log *ethTypes.Log) []byte {
	var result []byte
	for i, topic := range log.Topics {
		if i > 0 {
			result = append(result, topic.Bytes()...)
		}
	}

	return append(result, log.Data...)
}

// GetPowerFromAmount returns power from amount -- note that this will polute amount object
func GetPowerFromAmount(amount *big.Int) (*big.Int, error) {
	decimals18 := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)
	if amount.Cmp(decimals18) == -1 {
		return nil, errors.New("amount must be more than 1 token")
	}

	return amount.Div(amount, decimals18), nil
}

// GetAmountFromString converts string to its big Int
func GetAmountFromString(amount string) (*big.Int, error) {
	amountInDecimals, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, errors.New("cannot convert string to big int")
	}
	return amountInDecimals, nil
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

// EventByID looks up a event by the topic id
func EventByID(abiObject *abi.ABI, sigdata []byte) *abi.Event {
	for _, event := range abiObject.Events {
		if bytes.Equal(event.Id().Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}
