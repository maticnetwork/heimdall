package helper

import (
	"encoding/json"
	"fmt"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

type validatorPretty struct {
	Address cmn.HexBytes `json:"address"`
	Power   int64        `json:"power"`
}

func ValidatorsToString(vs []abci.Validator) string {
	s := make([]validatorPretty, len(vs))
	for i, v := range vs {
		s[i] = validatorPretty{
			Address: v.Address,
			Power:   v.Power,
		}
	}
	b, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

// If more or equal than 1/3 of total voting power changed in one block, then
// a light client could never prove the transition externally. See
// ./lite/doc.go for details on how a light client tracks validators.
func updateValidators(currentSet *types.ValidatorSet, abciUpdates []abci.ValidatorUpdate) error {
	updates, err := types.PB2TM.ValidatorUpdates(abciUpdates)
	if err != nil {
		return err
	}

	// these are tendermint types now
	for _, valUpdate := range updates {
		if valUpdate.VotingPower < 0 {
			return fmt.Errorf("Voting power can't be negative %v", valUpdate)
		}

		address := valUpdate.Address
		_, val := currentSet.GetByAddress(address)
		if valUpdate.VotingPower == 0 {
			// remove val
			_, removed := currentSet.Remove(address)
			if !removed {
				return fmt.Errorf("Failed to remove validator %X", address)
			}
		} else if val == nil {
			// add val
			added := currentSet.Add(valUpdate)
			if !added {
				return fmt.Errorf("Failed to add new validator %v", valUpdate)
			}
		} else {
			// update val
			updated := currentSet.Update(valUpdate)
			if !updated {
				return fmt.Errorf("Failed to update validator %X to %v", address, valUpdate)
			}
		}
	}
	return nil
}
