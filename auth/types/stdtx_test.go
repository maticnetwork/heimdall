package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	priv = secp256k1.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

func TestStdTx(t *testing.T) {
	msg := sdk.NewTestMsg(addr)
	sig := StdSignature{}

	tx := NewStdTx(msg, sig, "")
	require.Equal(t, msg, tx.GetMsgs()[0])
	require.Equal(t, sig, tx.GetSignatures()[0])

	feePayer := tx.GetSigners()[0]
	require.Equal(t, addr, feePayer)
}

func TestTxValidateBasic(t *testing.T) {
	ctx := sdk.NewContext(nil, abci.Header{ChainID: "mychainid"}, false, log.NewNopLogger())

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()

	// msg and signatures
	msg1 := sdk.NewTestMsg(addr1)
	tx := NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0))

	err := tx.ValidateBasic()
	require.NoError(t, err)
}

func TestDefaultTxEncoder(t *testing.T) {
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	encoder := DefaultTxEncoder(cdc)

	msg := sdk.NewTestMsg(addr)
	tx := NewStdTx(msg, StdSignature{}, "")

	cdcBytes, err := cdc.MarshalBinaryLengthPrefixed(tx)

	require.NoError(t, err)
	encoderBytes, err := encoder(tx)

	require.NoError(t, err)
	require.Equal(t, cdcBytes, encoderBytes)
}
