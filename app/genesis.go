package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisAccount genesis account
type GenesisAccount struct {
	Address       common.Address `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      int64          `json:"sequence_number"`
	AccountNumber int64          `json:"account_number"`
}

// GenesisValidator genesis validator
type GenesisValidator struct {
	Address    common.Address `json:"address"`
	StartEpoch uint64         `json:"start_epoch"`
	EndEpoch   uint64         `json:"end_epoch"`
	Power      uint64         `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey `json:"pub_key"`
	Signer     common.Address `json:"signer"`
}

// ToHeimdallValidator converts genesis validator validator to Heimdall validator
func (v *GenesisValidator) ToHeimdallValidator() hmTypes.Validator {
	return hmTypes.Validator{
		Address:    v.Address,
		PubKey:     v.PubKey,
		Power:      v.Power,
		StartEpoch: v.StartEpoch,
		EndEpoch:   v.EndEpoch,
		Signer:     v.Address,
	}
}

// GenesisState to Unmarshal
type GenesisState struct {
	Accounts   []GenesisAccount   `json:"accounts"`
	Validators []GenesisValidator `json:"validators"`
	GenTxs     []json.RawMessage  `json:"gentxs"`
}
