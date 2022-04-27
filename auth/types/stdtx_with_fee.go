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

// todo the last param can be make as optional with `...` 
func NewTx(blockNumber, hardForkBlockNumber uint64, msg sdk.Msg, sig StdSignature, memo string, fee StdFee) sdk.Tx {
	if blockNumber >= hardForkBlockNumber {
		return newStdTxWithFee(msg, sig, memo, fee)
	}

	return newStdTx(msg, sig, memo)
}

// StdTxWithFee is a standard way to wrap a Msg with Fee and Signatures.
type StdTxWithFee struct {
	StdTx
	Fee       StdFee       `json:"fee" yaml:"fee"`
}

// NewStdTxWithFee is function to get new std tx object
func newStdTxWithFee(msg sdk.Msg, sig StdSignature, memo string, fee StdFee) StdTxWithFee {
	return StdTxWithFee{
		StdTx: StdTx{
			Msg:       msg,
			Signature: sig,
			Memo:      memo,
		},
		Fee:       fee,
	}
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTxWithFee) ValidateBasic() sdk.Error {
	err := tx.StdTx.ValidateBasic()
	if err != nil {
		return err
	}

	if tx.Fee.Gas > maxGasWanted {
		return sdk.ErrGasOverflow(fmt.Sprintf("invalid gas supplied; %d > %d", tx.Fee.Gas, maxGasWanted))
	}
	if tx.Fee.Amount.IsAnyNegative() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.Amount))
	}

	return nil
}

//
// Decoders
//

// DefaultTxWithFeeDecoder logic for standard transaction decoding
func DefaultTxDecoder[T sdk.Tx](cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = new(T)

		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}

		// StdTxWithFee.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}

		return *tx, nil
	}
}

