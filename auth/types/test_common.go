package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewTestTx creates new test tx
func NewTestTx(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, fee StdFee) sdk.Tx {
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, fee, msg, "")
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := NewStdTx(msg, fee, sig, "")
	return tx
}

// NewTestTxWithMemo create new test tx
func NewTestTxWithMemo(ctx sdk.Context, msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, memo string, fee StdFee) sdk.Tx {
	signBytes := StdSignBytes(ctx.ChainID(), accNum, seq, fee, msg, "")
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := NewStdTx(msg, fee, sig, memo)
	return tx
}

// NewTestTxWithSignBytes creates tx with sign bytes
func NewTestTxWithSignBytes(msg sdk.Msg, priv crypto.PrivKey, accNum uint64, seq uint64, signBytes []byte, memo string, fee StdFee) sdk.Tx {
	sig, err := priv.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	tx := NewStdTx(msg, fee, sig, memo)
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
