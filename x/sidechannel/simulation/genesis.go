package simulation

import (
	"fmt"
	"math/rand"

	"github.com/maticnetwork/heimdall/x/sidechannel/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// RandomPastCommits returns random past commits value
func RandomPastCommits(r *rand.Rand, n int, txsN int, validatorsN int) []*types.PastCommit {
	result := make([]*types.PastCommit, n)
	for i := 0; i < n; i++ {
		txs := make([][]byte, txsN)
		for j := 0; j < txsN; j++ {
			s := fmt.Sprintf("test-transaction %v", j)
			txs[j] = []byte(s)
		}

		validators := make([]abci.Validator, validatorsN)
		for j := 0; j < validatorsN; j++ {
			validators[j] = abci.Validator{
				Address: []byte("validator" + string(rune(j))),
				Power:   r.Int63n(100000),
			}
		}

		result[i] = &types.PastCommit{
			Height: uint64(2) + r.Uint64(),
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
				Address: []byte("validator" + string(rune(i+1))),
				Power:   r.Int63n(100000),
			},
			SignedLastBlock: r.Int()%2 == 0,
		}
	}

	return abci.LastCommitInfo{
		Votes: votes,
	}
}
