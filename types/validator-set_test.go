package types

import (
	"encoding/hex"
	"testing"
)

// StringToPubkey converts string to Pubkey
func StringToPubkey(pubkeyStr string) PubKey {
	_pubkey, _ := hex.DecodeString(pubkeyStr)
	return NewPubKey(_pubkey)
}

func TestUpdateChanges(t *testing.T) {
	v1 := &Validator{
		ID:          1,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("04b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("6C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v2 := &Validator{
		ID:          1,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("04914873c8d5935837ade39cbdabd6efb3d3d4064c5918da11e555bba0ab2c58fee95974a3222830cf73d257bdc18cfcd01765482108a48e68bc0b657618acb40e"),
		Signer:      HexToHeimdallAddress("9fB29AAc15b9A4B7F17c3385939b007540f4d791"),
	}

	vset1 := NewValidatorSet([]*Validator{
		v1,
	})

	v1new := v1.Copy()
	v1new.VotingPower = 0
	vset1.UpdateWithChangeSet([]*Validator{
		v1new,
		v2.Copy(),
	})
	vset1.IncrementProposerPriority(1)

	if !vset1.GetProposer().Signer.Equals(v2.Signer) {
		t.Errorf("expected: %v, but got %v", v2.Signer, vset1.GetProposer().Signer)
	}
}
