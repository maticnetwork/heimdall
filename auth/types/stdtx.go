package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	_ sdk.Tx = (*StdTx)(nil)

	maxGasWanted = uint64((1 << 63) - 1)
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
	return nil
}

// GetSigners returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigners method for each Msg
// in the order they appear in tx.GetMsgs().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigners() []sdk.AccAddress {
	seen := map[string]bool{}
	var signers []sdk.AccAddress
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

// RLPTxDecoder decodes the txBytes to a StdTX
func RLPTxDecoder(pulp *Pulp) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		tx, err := pulp.DecodeBytes(txBytes)
		if err != nil {
			return nil, sdk.ErrTxDecode(err.Error())
		}

		return tx.(sdk.Tx), nil
	}
}

// RLPTxEncoder logic for RLP transaction encoding
func RLPTxEncoder(pulp *Pulp) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return pulp.EncodeToBytes(tx.(StdTx))
	}
}

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
