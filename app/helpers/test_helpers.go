package helpers

import (
	"errors"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/maticnetwork/heimdall/types/simulation"
)

// SimAppChainID hardcoded chainID for simulation
const (
	DefaultGenTxGas = 1000000
	SimAppChainID   = "simulation-app"
)

// GenTx generates a signed mock transaction.
func GenTx(gen client.TxConfig, msgs []sdk.Msg, feeAmt sdk.Coins, gas uint64, chainID string, accnums []uint64, seq []uint64, priv ...cryptotypes.PrivKey) (sdk.Tx, error) {
	// fee := authTypes.StdFee{
	// 	Amount: feeAmt,
	// 	Gas:    gas,
	// }

	if len(msgs) == 0 {
		panic(errors.New("Msgs cannot be empty"))
	}

	sigs := make([]signing.SignatureV2, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode := gen.SignModeHandler().DefaultMode()

	for i, p := range priv {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
			},
			Sequence: seq[i],
		}
	}
	tx := gen.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = tx.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// // 2nd round: once all signer infos are set, every signer can sign.
	// for i, p := range priv {
	// 	signerData := authsign.SignerData{
	// 		ChainID:       chainID,
	// 		AccountNumber: accNums[i],
	// 		Sequence:      accSeqs[i],
	// 	}
	// 	signBytes, err := gen.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sig, err := p.Sign(signBytes)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
	// 	err = tx.SetSignatures(sigs...)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	return tx.GetTx(), nil
}
