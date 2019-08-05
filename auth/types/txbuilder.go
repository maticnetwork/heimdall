package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCrypto "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder     sdk.TxEncoder
	keybase       crkeys.Keybase
	accountNumber uint64
	sequence      uint64
	chainID       string
	memo          string
}

// NewTxBuilder returns a new initialized TxBuilder.
func NewTxBuilder(
	txEncoder sdk.TxEncoder,
	accNumber uint64,
	seq uint64,
	gas uint64,
	gasAdj float64,
	simulateAndExecute bool,
	chainID string,
	memo string,
) TxBuilder {

	return TxBuilder{
		txEncoder:     txEncoder,
		keybase:       nil,
		accountNumber: accNumber,
		sequence:      seq,
		chainID:       chainID,
		memo:          memo,
	}
}

// NewTxBuilderFromCLI returns a new initialized TxBuilder with parameters from
// the command line using Viper.
func NewTxBuilderFromCLI() TxBuilder {
	kb, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		panic(err)
	}
	txbldr := TxBuilder{
		keybase:       kb,
		accountNumber: uint64(viper.GetInt64(client.FlagAccountNumber)),
		sequence:      uint64(viper.GetInt64(client.FlagSequence)),
		chainID:       viper.GetString(client.FlagChainID),
		memo:          viper.GetString(client.FlagMemo),
	}

	return txbldr
}

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxEncoder() sdk.TxEncoder { return bldr.txEncoder }

// AccountNumber returns the account number
func (bldr TxBuilder) AccountNumber() uint64 { return bldr.accountNumber }

// Sequence returns the transaction sequence
func (bldr TxBuilder) Sequence() uint64 { return bldr.sequence }

// Keybase returns the keybase
func (bldr TxBuilder) Keybase() crkeys.Keybase { return bldr.keybase }

// ChainID returns the chain id
func (bldr TxBuilder) ChainID() string { return bldr.chainID }

// Memo returns the memo message
func (bldr TxBuilder) Memo() string { return bldr.memo }

// WithTxEncoder returns a copy of the context with an updated codec.
func (bldr TxBuilder) WithTxEncoder(txEncoder sdk.TxEncoder) TxBuilder {
	bldr.txEncoder = txEncoder
	return bldr
}

// WithChainID returns a copy of the context with an updated chainID.
func (bldr TxBuilder) WithChainID(chainID string) TxBuilder {
	bldr.chainID = chainID
	return bldr
}

// WithKeybase returns a copy of the context with updated keybase.
func (bldr TxBuilder) WithKeybase(keybase crkeys.Keybase) TxBuilder {
	bldr.keybase = keybase
	return bldr
}

// WithSequence returns a copy of the context with an updated sequence number.
func (bldr TxBuilder) WithSequence(sequence uint64) TxBuilder {
	bldr.sequence = sequence
	return bldr
}

// WithMemo returns a copy of the context with an updated memo.
func (bldr TxBuilder) WithMemo(memo string) TxBuilder {
	bldr.memo = strings.TrimSpace(memo)
	return bldr
}

// WithAccountNumber returns a copy of the context with an account number.
func (bldr TxBuilder) WithAccountNumber(accnum uint64) TxBuilder {
	bldr.accountNumber = accnum
	return bldr
}

// BuildSignMsg builds a single message to be signed from a TxBuilder given a
// set of messages. It returns an error if a fee is supplied but cannot be
// parsed.
func (bldr TxBuilder) BuildSignMsg(msgs []sdk.Msg) (StdSignMsg, error) {
	if bldr.chainID == "" {
		return StdSignMsg{}, fmt.Errorf("chain ID required but not specified")
	}

	return StdSignMsg{
		ChainID:       bldr.chainID,
		AccountNumber: bldr.accountNumber,
		Sequence:      bldr.sequence,
		Memo:          bldr.memo,
		Msg:           msgs[0], // allow only one message
	}, nil
}

// Sign transaction with default node key
func (bldr TxBuilder) Sign(privKey secp256k1.PrivKeySecp256k1, msg StdSignMsg) ([]byte, error) {
	sig, err := MakeSignature(privKey, msg)
	if err != nil {
		return nil, err
	}

	return bldr.txEncoder(NewStdTx(msg.Msg, sig, msg.Memo))
}

// SignWithPassphrase signs a transaction given a name, passphrase, and a single message to
// signed. An error is returned if signing fails.
func (bldr TxBuilder) SignWithPassphrase(name, passphrase string, msg StdSignMsg) ([]byte, error) {
	sig, err := MakeSignatureWithKeybase(bldr.keybase, name, passphrase, msg)
	if err != nil {
		return nil, err
	}

	return bldr.txEncoder(NewStdTx(msg.Msg, sig, msg.Memo))
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a set of messages.
func (bldr TxBuilder) BuildAndSign(privKey secp256k1.PrivKeySecp256k1, msgs []sdk.Msg) ([]byte, error) {
	stdMsg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	return bldr.Sign(privKey, stdMsg)
}

// BuildAndSignWithPassphrase builds a single message to be signed, and signs a transaction
// with the built message given a name, passphrase, and a set of messages.
func (bldr TxBuilder) BuildAndSignWithPassphrase(name, passphrase string, msgs []sdk.Msg) ([]byte, error) {
	stdMsg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	return bldr.SignWithPassphrase(name, passphrase, stdMsg)
}

// BuildTxForSim creates a StdSignMsg and encodes a transaction with the
// StdSignMsg with a single empty StdSignature for tx simulation.
func (bldr TxBuilder) BuildTxForSim(msgs []sdk.Msg) ([]byte, error) {
	signMsg, err := bldr.BuildSignMsg(msgs)
	if err != nil {
		return nil, err
	}

	// the ante handler will populate with a sentinel pubkey
	sig := StdSignature{}
	return bldr.txEncoder(NewStdTx(signMsg.Msg, sig, signMsg.Memo))
}

// SignStdTxWithPassphrase appends a signature to a StdTx and returns a copy of it. If append
// is false, it replaces the signatures already attached with the new signature.
func (bldr TxBuilder) SignStdTxWithPassphrase(name, passphrase string, stdTx StdTx, appendSig bool) (signedStdTx StdTx, err error) {
	if bldr.chainID == "" {
		return StdTx{}, fmt.Errorf("chain ID required but not specified")
	}

	stdSignature, err := MakeSignatureWithKeybase(bldr.keybase, name, passphrase, StdSignMsg{
		ChainID:       bldr.chainID,
		AccountNumber: bldr.accountNumber,
		Sequence:      bldr.sequence,
		Msg:           stdTx.GetMsgs()[0],
		Memo:          stdTx.GetMemo(),
	})
	if err != nil {
		return
	}

	signedStdTx = NewStdTx(stdTx.GetMsgs()[0], stdSignature, stdTx.GetMemo())
	return
}

// MakeSignature builds a StdSignature for given a StdSignMsg.
func MakeSignature(privKey secp256k1.PrivKeySecp256k1, msg StdSignMsg) (sig StdSignature, err error) {
	return ethCrypto.Sign(msg.Bytes(), privKey[:])
}

// MakeSignatureWithKeybase builds a StdSignature given keybase, key name, passphrase, and a StdSignMsg.
func MakeSignatureWithKeybase(
	keybase crkeys.Keybase,
	name string,
	passphrase string,
	msg StdSignMsg,
) (sig StdSignature, err error) {
	if keybase == nil {
		keybase, err = keys.NewKeyBaseFromHomeFlag()
		if err != nil {
			return
		}
	}

	sigBytes, _, err := keybase.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return
	}
	return sigBytes, nil
}
