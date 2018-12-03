package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisAccount genesis account
type GenesisAccount struct {
	Address       common.Address `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      int64          `json:"sequence_number"`
	AccountNumber int64          `json:"account_number"`
}

// TODO not sure why int64 is used , if no specific reason , convert to uint64 can hold larger values
// GenesisValidator genesis validator
type GenesisValidator struct {
	Address    common.Address `json:"address"`
	StartEpoch uint64         `json:"start_epoch"`
	EndEpoch   uint64         `json:"end_epoch"`
	Power      uint64         `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey `json:"pub_key"`
	Signer     common.Address `json:"signer"`
}

// ToTmValidator converts genesis valdator validator to Tendermint validator
func (v *GenesisValidator) ToTmValidator() tmTypes.Validator {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	copy(pubkeyBytes[:], v.PubKey[:])

	return tmTypes.Validator{
		Address:     v.Signer.Bytes(),
		PubKey:      pubkeyBytes,
		VotingPower: int64(v.Power),
	}
}

// ToHeimdallValidator converts genesis validator validator to Heimdall validator
func (v *GenesisValidator) ToHeimdallValidator() hmTypes.Validator {
	return hmTypes.Validator{
		Address:    v.Address,
		PubKey:     helper.BytesToPubkey(v.PubKey[:]),
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
