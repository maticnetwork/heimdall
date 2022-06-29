package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTestTx creates new test tx
func NewTestTx(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64) sdk.Tx {
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, msg, "")

	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := NewStdTx(msg, sig, "")

	return tx
}

// NewTestTxWithMemo create new test tx
func NewTestTxWithMemo(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, memo string) sdk.Tx {
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, msg, "")

	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := NewStdTx(msg, sig, memo)

	return tx
}
