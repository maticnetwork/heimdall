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
	ID         hmTypes.ValidatorID `json:"id"`
	StartEpoch uint64              `json:"start_epoch"`
	EndEpoch   uint64              `json:"end_epoch"`
	Power      uint64              `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey      `json:"pub_key"`
	Signer     common.Address      `json:"signer"`
}

// HeimdallValidator converts genesis validator validator to Heimdall validator
func (v *GenesisValidator) HeimdallValidator() hmTypes.Validator {
	return hmTypes.Validator{
		ID:         v.ID,
		PubKey:     v.PubKey,
		Power:      v.Power,
		StartEpoch: v.StartEpoch,
		EndEpoch:   v.EndEpoch,
		Signer:     v.Signer,
	}
}

// GenesisState to Unmarshal
type GenesisState struct {
	BufferedCheckpoint hmTypes.CheckpointBlockHeader   `json:"buffered_checkpoint"`
	CheckpointCache    bool                            `json:"checkpoint_cache"`
	CheckpointACKCache bool                            `json:"ack_cache"`
	LastNoACK          uint64                           `json:"last_no_ack"`
	AckCount           uint64                          `json:"ack_count"`
	GenValidators      []GenesisValidator              `json:"gen_validators"`
	Validators         []hmTypes.Validator             `json:"validators"`
	CurrentValSet      hmTypes.ValidatorSet            `json:"current_val_set"`
	GenTxs             []json.RawMessage               `json:"gentxs"`
	Accounts           []GenesisAccount                `json:"accounts"`
	Headers            []hmTypes.CheckpointBlockHeader `json:"headers"`
}
