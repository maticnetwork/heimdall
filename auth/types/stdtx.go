package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	_ sdk.Tx = (*StdTx)(nil)
)

// StdTx is a standard way to wrap a Msg with Fee and Signatures.
type StdTx struct {
	Msg       sdk.Msg      `json:"msg" yaml:"msg"`
	Signature StdSignature `json:"signature" yaml:"signature"`
	Memo      string       `json:"memo" yaml:"memo"`
}

// StdTxRaw is a standard way to wrap a RLP Msg with Fee and Signatures.
type StdTxRaw struct {
	Msg       rlp.RawValue
	Signature StdSignature
	Memo      string
}

// NewStdTx is function to get new std tx object
func NewStdTx(msg sdk.Msg, sig StdSignature, memo string) StdTx {
	return StdTx{
		Msg:       msg,
		Signature: sig,
		Memo:      memo,
	}
}

// GetMsgs returns the all the transaction's messages.
func (tx StdTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.Msg}
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	stdSigs := tx.GetSignatures()

	if tx.Signature.Empty() {
		return sdk.ErrNoSignatures("No signers")
	}

	if len(stdSigs) != 1 || len(stdSigs) != len(tx.GetSigners()) {
		return sdk.ErrUnauthorized("wrong number of signers")
	}

	return nil
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []sdk.AccAddress {
	var (
		signers []sdk.AccAddress
		seen    = map[string]bool{}
	)

	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	return signers
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string {
	return tx.Memo
}

// GetSignatures returns the signature of signers who signed the Msg.
func (tx StdTx) GetSignatures() []StdSignature {
	return []StdSignature{tx.Signature}
}

//
// Std signature
//

// StdSignature represents a sig
type StdSignature []byte

// Bytes returns the bytes
func (ss StdSignature) Bytes() []byte {
	return ss[:]
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (ss StdSignature) Marshal() ([]byte, error) {
	return ss.Bytes(), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (ss *StdSignature) Unmarshal(data []byte) error {
	*ss = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (ss StdSignature) MarshalJSON() ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(ss.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (ss StdSignature) MarshalYAML() (interface{}, error) {
	return ss.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (ss *StdSignature) UnmarshalJSON(data []byte) error {
	var s string
	if err := jsoniter.ConfigFastest.Unmarshal(data, &s); err != nil {
		return err
	}

	*ss = common.FromHex(s)

	return nil
}

// Empty checks is sig is empty
func (ss StdSignature) Empty() bool {
	return len(ss.Bytes()) == 0
}

// String implements the Stringer interface.
func (ss StdSignature) String() string {
	if ss.Empty() {
		return ""
	}

	return hexutil.Encode(ss)
}

//
// Std fee
//

// StdFee includes the amount of coins paid in fees and the maximum
// gas to be used by the transaction. The ratio yields an effective "gasprice",
// which must be above some miminum to be accepted into the mempool.
type StdFee struct {
	Amount sdk.Coins `json:"amount"`
	Gas    uint64    `json:"gas"`
}

// NewStdFee returns a new instance of StdFee
func NewStdFee(gas uint64, amount sdk.Coins) StdFee {
	return StdFee{
		Amount: amount,
		Gas:    gas,
	}
}

// Bytes for signing later
func (fee StdFee) Bytes() []byte {
	// normalize. XXX
	// this is a sign of something ugly
	// (in the lcd_test, client side its null,
	// server side its [])
	if len(fee.Amount) == 0 {
		fee.Amount = sdk.NewCoins()
	}

	bz, err := ModuleCdc.MarshalJSON(fee) // TODO
	if err != nil {
		panic(err)
	}

	return bz
}

// GasPrices returns the gas prices for a StdFee.
//
// NOTE: The gas prices returned are not the true gas prices that were
// originally part of the submitted transaction because the fee is computed
// as fee = ceil(gasWanted * gasPrices).
func (fee StdFee) GasPrices() sdk.DecCoins {
	return sdk.NewDecCoins(fee.Amount).QuoDec(sdk.NewDec(int64(fee.Gas)))
}

//
// Decoders
//

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = StdTx{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}

		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}
