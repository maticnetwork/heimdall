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
	t.Parallel()

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

	v3 := &Validator{
		ID:          3,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 100,
		PubKey:      StringToPubkey("03b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("3C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v4 := &Validator{
		ID:          4,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 1,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	vset1 := NewValidatorSet([]*Validator{
		v1,
	})

	v1new := v1.Copy()
	v1new.VotingPower = 0

	err := vset1.UpdateWithChangeSet([]*Validator{
		v1new,
		v2.Copy(),
		v3.Copy(),
	})
	if err != nil {
		t.Error(err)
	}

	vset1.IncrementProposerPriority(1)

	if !vset1.GetProposer().Signer.Equals(v3.Signer) {
		t.Errorf("expected: %v, but got %v", v2.Signer, vset1.GetProposer().Signer)
	}

	err = vset1.UpdateWithChangeSet([]*Validator{
		v4.Copy(),
	})

	vset1Temp := vset1.Copy()

	vset1.RescalePriorities(PriorityWindowSizeFactor * vset1.TotalVotingPower())
	vset1.shiftByAvgProposerPriority()

	//shifting the proposer prioirity two times
	vset1.RescalePriorities(PriorityWindowSizeFactor * vset1.TotalVotingPower())
	vset1.shiftByAvgProposerPriority()

	for _, val := range vset1.Validators {
		address := val.Signer.Bytes()

		_, val2 := vset1Temp.GetByAddress(address)

		if val2 != nil && val.ProposerPriority != val2.ProposerPriority {
			t.Errorf("Proposer priority should not change when rescaling with same factor second time. ValOld Proposer Priority %v ValNew Proposer Priority %v", val2.ProposerPriority, val.ProposerPriority)
		}
	}
}

func TestUpdateChangesWithoutIncAccum(t *testing.T) {
	t.Parallel()

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

	v3 := &Validator{
		ID:          3,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 100,
		PubKey:      StringToPubkey("03b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("3C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v4 := &Validator{
		ID:          4,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 1,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	vset1 := NewValidatorSet([]*Validator{
		v1,
	})

	vset2 := vset1.Copy()

	v1new := v1.Copy()
	v1new.VotingPower = 0

	err := vset1.UpdateWithChangeSet([]*Validator{
		v1new,
		v2.Copy(),
		v3.Copy(),
		v4.Copy(),
	})

	if err != nil {
		t.Error(err)
	}

	vset1.IncrementProposerPriority(1)

	err = vset2.UpdateWithChangeSet([]*Validator{
		v1new,
		v2.Copy(),
		v3.Copy(),
		v4.Copy(),
	})

	if len(vset1.Validators) != len(vset2.Validators) {
		t.Errorf("Both the validator sets should have same length,VSet1 Length %v VSet2 Length %v", len(vset1.Validators), len(vset2.Validators))
	}

	for _, val := range vset1.Validators {
		address := val.Signer.Bytes()

		_, val2 := vset2.GetByAddress(address)

		if val2 == nil {
			t.Errorf("Validor for the particular address %v  should not be nil", address)
		}

		if val.VotingPower != val2.VotingPower || val.Nonce != val2.Nonce || val.ID != val2.ID || val.StartEpoch != val2.StartEpoch || val.EndEpoch != val2.EndEpoch || val.Jailed != val2.Jailed {
			t.Errorf("val and val2 should be equal in all the properties except proposer priority")
		}
	}
}
