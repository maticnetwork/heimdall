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

	tx := newStdTx(msg, sig, "")
	return tx
}

// NewTestTxWithFee creates new test tx
func NewTestTxWithFee(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, fee StdFee) sdk.Tx {
	signBytes := StdSignBytesWithFee(ctx.ChainID(), accNum, seq, fee, msg, "")
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := newStdTxWithFee(msg, sig, "", fee)
	return tx
}

// NewTestTxWithMemo create new test tx
func NewTestTxWithMemo(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, memo string) sdk.Tx {
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, msg, "")
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := newStdTx(msg, sig, memo)
	return tx
}

// NewTestTxWithMemoWithFee create new test tx
func NewTestTxWithMemoWithFee(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, memo string, fee StdFee) sdk.Tx {
	signBytes := StdSignBytesWithFee(ctx.ChainID(), accNum, seq, fee, msg, "")
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := newStdTxWithFee(msg, sig, memo, fee)
	return tx
}

// NewTestTxWithSignBytes creates tx with sign bytes
func NewTestTxWithSignBytes(msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, signBytes []byte, memo string) sdk.Tx {
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := newStdTx(msg, sig, memo)
	return tx
}

// NewTestTxWithSignBytesWithFee creates tx with sign bytes
func NewTestTxWithSignBytesWithFee(msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, signBytes []byte, memo string, fee StdFee) sdk.Tx {
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := newStdTxWithFee(msg, sig, memo, fee)
	return tx
}

func NewTestStdFee() StdFee {
	return NewStdFee(50000,
		sdk.NewCoins(sdk.NewInt64Coin("matic", 150)),
	)
}

// coins to more than cover the fee
func NewTestCoins() sdk.Coins {
	return sdk.Coins{
		sdk.NewInt64Coin("matic", 10000000),
	}
}
