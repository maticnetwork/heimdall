package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Tx = (*StdTxWithFee)(nil)

	maxGasWanted = uint64((1 << 63) - 1)
)

// StdTxWithFee is a standard way to wrap a Msg with Fee and Signatures.
type StdTxWithFee struct {
	StdTx
	Fee StdFee `json:"fee" yaml:"fee"`
}

// NewStdTxWithFee is function to get new std tx object
func NewStdTxWithFee(msg sdk.Msg, fee StdFee, sig StdSignature, memo string) StdTxWithFee {
	return StdTxWithFee{
		StdTx: StdTx{
			Msg:       msg,
			Signature: sig,
			Memo:      memo,
		},
		Fee: fee,
	}
}

// GetMsgsWithFee returns the all the transaction's messages.
func (tx StdTxWithFee) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.Msg}
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTxWithFee) ValidateBasic() sdk.Error {
	stdSigs := tx.GetSignatures()

	if tx.Fee.Gas > maxGasWanted {
		return sdk.ErrGasOverflow(fmt.Sprintf("invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted))
	}
	if tx.Fee.Amount.IsAnyNegative() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.Amount))
	}

	if len(stdSigs) == 0 {
		return sdk.ErrNoSignatures("no signers")
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
func (tx StdTxWithFee) GetSigners() []sdk.AccAddress {
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
func (tx StdTxWithFee) GetMemo() string {
	return tx.Memo
}

// GetSignatures returns the signature of signers who signed the Msg.
func (tx StdTxWithFee) GetSignatures() []StdSignature {
	return []StdSignature{tx.Signature}
}

//
// Decoders
//

// DefaultTxWithFeeDecoder logic for standard transaction decoding
func DefaultTxWithFeeDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = StdTxWithFee{}

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTxWithFee.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}

		return tx, nil
	}
}

// DefaultTxWithFeeEncoder logic for standard transaction encoding
func DefaultTxWithFeeEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}
