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
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
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

// GetVoteSigs returns sigs bytes from vote
func GetVoteSigs(unFilteredVotes []*tmTypes.CommitSig) (sigs []byte) {
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

type sideTxSig struct {
	Address []byte
	Sig     []byte
}

// GetSideTxSigs returns sigs bytes from vote by tx hash
func GetSideTxSigs(txHash []byte, sideTxData []byte, unFilteredVotes []*tmTypes.CommitSig) (sigs [][3]*big.Int, err error) {
	// side tx result with data
	sideTxResultWithData := tmTypes.SideTxResultWithData{
		SideTxResult: tmTypes.SideTxResult{
			TxHash: txHash,
			Result: int32(abci.SideTxResultType_Yes),
		},
		Data: sideTxData,
	}

	// draft signed data
	signedData := sideTxResultWithData.GetBytes()

	sideTxSigs := make([]*sideTxSig, 0)
	for _, vote := range unFilteredVotes {
		if vote != nil {
			// iterate through all side-tx results
			for _, sideTxResult := range vote.SideTxResults {
				// find side-tx result by tx-hash
				if bytes.Equal(sideTxResult.TxHash, txHash) &&
					len(sideTxResult.Sig) == 65 &&
					sideTxResult.Result == int32(abci.SideTxResultType_Yes) {
					// validate sig
					var pk secp256k1.PubKeySecp256k1
					if p, err := authTypes.RecoverPubkey(signedData, sideTxResult.Sig); err == nil {
						copy(pk[:], p[:])

						// if it has valid sig, add it into side-tx sig array
						if bytes.Equal(vote.ValidatorAddress.Bytes(), pk.Address().Bytes()) {
							sideTxSigs = append(sideTxSigs, &sideTxSig{
								Address: vote.ValidatorAddress.Bytes(),
								Sig:     sideTxResult.Sig,
							})
						}
					}
				}
				// break
			}
		}
	}

	// Nothing to do with sigs, Just a type check in latest geth code
	dummyLegacyTxn := ethTypes.NewTransaction(0, common.Address{}, nil, 0, nil, nil)

	if len(sideTxSigs) > 0 {
		// sort sigs by address
		sort.Slice(sideTxSigs, func(i, j int) bool {
			return bytes.Compare(sideTxSigs[i].Address, sideTxSigs[j].Address) < 0
		})

		// loop votes and append to sig to sigs
		for _, sideTxSig := range sideTxSigs {
			R, S, V, err := ethTypes.HomesteadSigner{}.SignatureValues(dummyLegacyTxn, sideTxSig.Sig)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, [3]*big.Int{R, S, V})
		}
	}

	return sigs, nil
}

// GetVoteBytes returns vote bytes
func GetVoteBytes(unFilteredVotes []*tmTypes.CommitSig, chainID string) []byte {
	var vote *tmTypes.CommitSig
	for _, item := range unFilteredVotes {
		if item != nil {
			vote = item
			break
		}
	}

	// if vote not found, return empty bytes
	if vote == nil {
		return []byte{}
	}

	v := tmTypes.Vote(*vote)
	// sign bytes for vote
	return v.SignBytes(chainID)
}

// GetTxEncoder returns tx encoder
func GetTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return authTypes.DefaultTxEncoder(cdc)
}

// GetTxDecoder returns tx decoder
func GetTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return authTypes.DefaultTxDecoder(cdc)
}

// GetStdTxBytes get tx bytes
func GetStdTxBytes(cliCtx context.CLIContext, tx authTypes.StdTx) ([]byte, error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder(cliCtx.Codec))
	return txBldr.GetStdTxBytes(tx)
}

// BroadcastMsgs creates transaction and broadcasts it
func BroadcastMsgs(cliCtx context.CLIContext, msgs []sdk.Msg) (sdk.TxResponse, error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder(cliCtx.Codec))
	return BuildAndBroadcastMsgs(cliCtx, txBldr, msgs)
}

// BroadcastTx broadcasts transaction
func BroadcastTx(cliCtx context.CLIContext, tx authTypes.StdTx, mode string) (res sdk.TxResponse, err error) {
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder(cliCtx.Codec))

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
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder(cliCtx.Codec))

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
	txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(GetTxEncoder(cliCtx.Codec))

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
	Logger.Debug("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes), "txHash", hex.EncodeToString(tmTypes.Tx(txBytes).Hash()))
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
	innerPrefix = []byte{1}

	leftPrefix  = []byte{0}
	rightPrefix = []byte{1}
)

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

// GetAmountFromPower returns amount from power
func GetAmountFromPower(power int64) (*big.Int, error) {
	pow := big.NewInt(0).SetInt64(power)
	decimals18 := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)
	return pow.Mul(pow, decimals18), nil
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
		if bytes.Equal(event.ID.Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}

// GetHeimdallServerEndpoint returns heimdall server endpoint
func GetHeimdallServerEndpoint(endpoint string) string {
	u, _ := url.Parse(GetConfig().HeimdallServerURL)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

// FetchFromAPI fetches data from any URL
func FetchFromAPI(cliCtx cliContext.CLIContext, URL string) (result rest.ResponseWithHeight, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// response
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		// unmarshall data from buffer
		var response rest.ResponseWithHeight
		if err := cliCtx.Codec.UnmarshalJSON(body, &response); err != nil {
			return result, err
		}
		return response, nil
	}

	Logger.Debug("Error while fetching data from URL", "status", resp.StatusCode, "URL", URL)
	return result, fmt.Errorf("Error while fetching data from url: %v, status: %v", URL, resp.StatusCode)
}
