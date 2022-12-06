package types

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	crkeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethCrypto "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder          sdk.TxEncoder
	keybase            crkeys.Keybase
	accountNumber      uint64
	sequence           uint64
	gas                uint64
	gasAdjustment      float64
	simulateAndExecute bool
	chainID            string
	memo               string
	fees               sdk.Coins
	gasPrices          sdk.DecCoins
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
	fees sdk.Coins,
	gasPrices sdk.DecCoins,
) TxBuilder {
	return TxBuilder{
		txEncoder:          txEncoder,
		keybase:            nil,
		accountNumber:      accNumber,
		sequence:           seq,
		gas:                gas,
		gasAdjustment:      gasAdj,
		simulateAndExecute: simulateAndExecute,
		chainID:            chainID,
		memo:               memo,
		fees:               fees,
		gasPrices:          gasPrices,
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
		keybase:            kb,
		accountNumber:      uint64(viper.GetInt64(client.FlagAccountNumber)),
		sequence:           uint64(viper.GetInt64(client.FlagSequence)),
		gas:                client.GasFlagVar.Gas,
		gasAdjustment:      viper.GetFloat64(client.FlagGasAdjustment),
		simulateAndExecute: client.GasFlagVar.Simulate,
		chainID:            viper.GetString(client.FlagChainID),
		memo:               viper.GetString(client.FlagMemo),
	}

	return txbldr
}

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxEncoder() sdk.TxEncoder { return bldr.txEncoder }

// AccountNumber returns the account number
func (bldr TxBuilder) AccountNumber() uint64 { return bldr.accountNumber }

// Sequence returns the transaction sequence
func (bldr TxBuilder) Sequence() uint64 { return bldr.sequence }

// Gas returns the gas for the transaction
func (bldr TxBuilder) Gas() uint64 { return bldr.gas }

// GasAdjustment returns the gas adjustment
func (bldr TxBuilder) GasAdjustment() float64 { return bldr.gasAdjustment }

// Keybase returns the keybase
func (bldr TxBuilder) Keybase() crkeys.Keybase { return bldr.keybase }

// SimulateAndExecute returns the option to simulate and then execute the transaction
// using the gas from the simulation results
func (bldr TxBuilder) SimulateAndExecute() bool { return bldr.simulateAndExecute }

// ChainID returns the chain id
func (bldr TxBuilder) ChainID() string { return bldr.chainID }

// Memo returns the memo message
func (bldr TxBuilder) Memo() string { return bldr.memo }

// Fees returns the fees for the transaction
func (bldr TxBuilder) Fees() sdk.Coins { return bldr.fees }

// GasPrices returns the gas prices set for the transaction, if any.
func (bldr TxBuilder) GasPrices() sdk.DecCoins { return bldr.gasPrices }

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

// WithGas returns a copy of the context with an updated gas.
func (bldr TxBuilder) WithGas(gas uint64) TxBuilder {
	bldr.gas = gas
	return bldr
}

// WithFees returns a copy of the context with an updated fee.
func (bldr TxBuilder) WithFees(fees string) TxBuilder {
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		panic(err)
	}

	bldr.fees = parsedFees

	return bldr
}

// WithGasPrices returns a copy of the context with updated gas prices.
func (bldr TxBuilder) WithGasPrices(gasPrices string) TxBuilder {
	parsedGasPrices, err := sdk.ParseDecCoins(gasPrices)
	if err != nil {
		panic(err)
	}

	bldr.gasPrices = parsedGasPrices

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

// SignStdTx appends a signature to a StdTx and returns a copy of it. If append
// is false, it replaces the signatures already attached with the new signature.
func (bldr TxBuilder) SignStdTx(privKey secp256k1.PrivKeySecp256k1, stdTx StdTx, appendSig bool) (signedStdTx StdTx, err error) {
	if bldr.chainID == "" {
		return StdTx{}, fmt.Errorf("chain ID required but not specified")
	}

	signMsg := StdSignMsg{
		ChainID:       bldr.chainID,
		AccountNumber: bldr.accountNumber,
		Sequence:      bldr.sequence,
		Memo:          stdTx.Memo,
		Msg:           stdTx.Msg, // allow only one message
	}

	sig, err := MakeSignature(privKey, signMsg)
	if err != nil {
		return
	}

	signedStdTx = NewStdTx(signMsg.Msg, sig, signMsg.Memo)

	return
}

// GetStdTxBytes get tx bytes
func (bldr TxBuilder) GetStdTxBytes(stdTx StdTx) (result []byte, err error) {
	return bldr.txEncoder(stdTx)
}

// MakeSignature builds a StdSignature for given a StdSignMsg.
func MakeSignature(privKey secp256k1.PrivKeySecp256k1, msg StdSignMsg) (sig StdSignature, err error) {
	data := crypto.Keccak256(msg.Bytes())
	return ethCrypto.Sign(data, privKey[:])
}

// RecoverPubkey builds a StdSignature for given a StdSignMsg.
func RecoverPubkey(msg []byte, sig []byte) ([]byte, error) {
	data := crypto.Keccak256(msg)
	return ethCrypto.RecoverPubkey(data, sig[:])
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
