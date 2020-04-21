package helpers

import (
	"errors"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// SimAppChainID hardcoded chainID for simulation
const (
	DefaultGenTxGas = 1000000
	SimAppChainID   = "simulation-app"
)

// GenTx generates a signed mock transaction.
func GenTx(msgs []sdk.Msg, feeAmt sdk.Coins, gas uint64, chainID string, accnums []uint64, seq []uint64, priv ...crypto.PrivKey) authTypes.StdTx {
	// fee := authTypes.StdFee{
	// 	Amount: feeAmt,
	// 	Gas:    gas,
	// }

	if len(msgs) == 0 {
		panic(errors.New("Msgs cannot be empty"))
	}

	sigs := make([]authTypes.StdSignature, len(priv))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))

	for i, p := range priv {
		// use a empty chainID for ease of testing
		sig, err := p.Sign(authTypes.StdSignBytes(chainID, accnums[i], seq[i], msgs[0], memo))
		if err != nil {
			panic(err)
		}

		sigs[i] = sig
	}

	return authTypes.NewStdTx(msgs[0], sigs[0], memo)
}
