package simulation

import (
	"fmt"
	"math/rand"

	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/sidechannel/types"
)

const validator = "validator"

// RandomPastCommits returns random past commits value
func RandomPastCommits(r *rand.Rand, n int, txsN int, validatorsN int) []types.PastCommit {
	result := make([]types.PastCommit, n)

	for i := 0; i < n; i++ {
		txs := make([]tmTypes.Tx, txsN)

		for j := 0; j < txsN; j++ {
			s := fmt.Sprintf("test-transaction %v", j)
			txs[j] = []byte(s)
		}

		validators := make([]abci.Validator, validatorsN)
		for j := 0; j < validatorsN; j++ {
			validators[j] = abci.Validator{
				Address: []byte(validator + fmt.Sprintf("%d", j)),
				Power:   r.Int63n(100000),
			}
		}

		result[i] = types.PastCommit{
			Height: 2 + r.Int63n(10000),
			Txs:    txs,
		}
	}

	return result
}

// RandomLastCommitInfo returns random last commit info
func RandomLastCommitInfo(r *rand.Rand, votesN int) abci.LastCommitInfo {
	votes := make([]abci.VoteInfo, votesN)

	for i := 0; i < votesN; i++ {
		votes[i] = abci.VoteInfo{
			Validator: abci.Validator{
				Address: []byte(validator + fmt.Sprintf("%d", i+1)),
				Power:   r.Int63n(100000),
			},
			SignedLastBlock: r.Int()%2 == 0,
		}
	}

	return abci.LastCommitInfo{
		Votes: votes,
	}
}
