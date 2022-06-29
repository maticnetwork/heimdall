package types

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"
)

var (
	priv = secp256k1.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

func TestStdTx(t *testing.T) {
	t.Parallel()

	msg := sdk.NewTestMsg(addr)
	sig := StdSignature{}

	tx := NewStdTx(msg, sig, "")
	require.Equal(t, msg, tx.GetMsgs()[0])
	require.Equal(t, sig, tx.GetSignatures()[0])

	feePayer := tx.GetSigners()[0]
	require.Equal(t, addr, feePayer)
}

func TestTxValidateBasic(t *testing.T) {
	t.Parallel()

	ctx := sdk.NewContext(nil, abci.Header{ChainID: "mychainid"}, false, log.NewNopLogger())

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()

	// msg and signatures
	msg1 := sdk.NewTestMsg(addr1)
	tx := NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0))

	require.NotNil(t, msg1)

	err := tx.ValidateBasic()
	require.Nil(t, err)
	require.NoError(t, err)
	require.NotPanics(t, func() { msg1.GetSignBytes() })
}

func TestDefaultTxEncoder(t *testing.T) {
	t.Parallel()

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

func TestTxDecode(t *testing.T) {
	t.Parallel()

	tx, err := base64.StdEncoding.DecodeString("wWhvHPg6AQHY1wEBlP+zHe/ZNZTQii57ULFjrJulHewY2NcBAZT/sx3v2TWU0Ioue1CxY6ybpR3sGICEXTLzJQ==")
	require.NoError(t, err)

	expected := "c1686f1cf83a0101d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec18d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec1880845d32f325"
	require.Equal(t, expected, hex.EncodeToString(tx), "Tx encoding should match")
}

func TestTxHash(t *testing.T) {
	t.Parallel()

	txStr := "AANQR/im+GCUHE8PBUoNahQVOC3A/YPGU1GIsiCAggP/oAUa5K2J62X6bWX065hIawNsvuv3z2qU4ObSU8l7Mgm0oMCuqfNQzHmirstq75vRV+hkFczlWh9VjSGNn8JQCo3YhF5C8VG4QZGyoPc937dVz4DrkdYdDRwnigW0qiIE+yMVS/Drcdt9FXol4Tzegb+1qIQbP+EXUnnFLFAuaeUF7A3Rs8WajjUBgA=="
	txHashStr := "b4560c30b12ebae71977373bcca2b0b553ae510efc4b167b4ebe7925f6e98557"

	txBz, err := base64.StdEncoding.DecodeString(txStr)
	require.NoError(t, err)

	var tx tmTypes.Tx = txBz

	require.Equal(t, txHashStr, hex.EncodeToString(tx.Hash()))
}
