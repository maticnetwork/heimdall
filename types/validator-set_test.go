package types

import (
	"encoding/hex"
	"math/rand"
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

func TestIncrementPriority(t *testing.T) {
	t.Parallel()

	v1 := &Validator{
		ID:          1,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 30,
		PubKey:      StringToPubkey("04b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("6C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v2 := &Validator{
		ID:          2,
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
		VotingPower: 10,
		PubKey:      StringToPubkey("03b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("3C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v4 := &Validator{
		ID:          4,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v5 := &Validator{
		ID:          5,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2t6e3d45a7age12c4b2158f79b9ffffffebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005CFF19C9D"),
	}

	v6 := &Validator{
		ID:          6,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7gte12c4b21t8f7gb95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4T29696E005C2319C9D"),
	}

	v7 := &Validator{
		ID:          7,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28ebe5ad5yc4g9be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E00GC23Y9C9D"),
	}

	v8 := &Validator{
		ID:          8,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace12c4b2158f79b95e4c28zbe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005CZ319C9D"),
	}

	vset1 := NewValidatorSet([]*Validator{
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
	})

	vset2 := vset1.Copy()

	var count int16 = 0

	vset1.IncrementProposerPriority(1)

	for i := 0; i < 100; i++ {
		if vset1.Proposer.Signer == v1.Signer {
			count++
		}

		rNum := rand.Intn(8)

		if rNum != 0 {
			vset1.IncrementProposerPriority(rNum)
		}

		vset1.IncrementProposerPriority(1)
	}

	t.Errorf("Count with Proposer Change at Staking %v", count)

	count = 0
	vset2.IncrementProposerPriority(1)

	for i := 0; i < 100; i++ {
		if vset2.Proposer.Signer == v1.Signer {
			count++
		}

		vset2.IncrementProposerPriority(1)
	}

	t.Errorf("Count with Proposer%v", count)

}

func TestValidatorSetUpdateWithRotation(t *testing.T) {
	t.Parallel()

	v1 := &Validator{
		ID:          1,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 100,
		PubKey:      StringToPubkey("04b12d8b2f6e3d45a7ace12c432158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("6C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v2 := &Validator{
		ID:          2,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 90,
		PubKey:      StringToPubkey("04914873c8d5935837ade39cbdabd6efb4d3d4064c5918da11e555bba0ab2c58fee95974a3222830cf73d257bdc18cfcd01765482108a48e68bc0b657618acb40e"),
		Signer:      HexToHeimdallAddress("9fB29AAc15b9A4B7F17c3385939b007540f4d791"),
	}

	v3 := &Validator{
		ID:          3,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 80,
		PubKey:      StringToPubkey("03b12d8b2f6e3d45a7ace12c4b2158f74b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("3C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	v4 := &Validator{
		ID:          4,
		StartEpoch:  0,
		EndEpoch:    0,
		VotingPower: 10,
		PubKey:      StringToPubkey("02b12d8b2f6e3d45a7ace42c4b2158f79b95e4c28ebe5ad54c439be9431d7fc9dc1164210bf6a5c3b8523528b931e772c86a307e8cff4b725e6b4a77d21417bf19"),
		Signer:      HexToHeimdallAddress("2C468CF8C9879006E22EC4029696E005C2319C9D"),
	}

	valSet := NewValidatorSet([]*Validator{
		v1,
		v2,
		v3,
		v4,
	})

	proposer := valSet.Proposer

	var newProposer *Validator
	var valChg *Validator

	//
	//Case 1 Increase Stake of Proposer
	//
	valChg = proposer.Copy()
	valChg.VotingPower = proposer.VotingPower + 100

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//
	//Case 2 Decrease Stake of Proposer
	//
	valChg = proposer.Copy()
	valChg.VotingPower = proposer.VotingPower - 5

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//
	//Case 3 Increase Validator Stake other than Proposer
	//
	valChg = v3.Copy()
	valChg.VotingPower = v3.VotingPower + 100

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//
	//Case 4 Decrease Validator Stake other then Proposer
	//
	valChg = v3.Copy()
	valChg.VotingPower = v3.VotingPower + 100

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//
	//Case 5 Increase Proposer Stake to Max
	//
	valChg = proposer.Copy()
	valChg.VotingPower = 1152921504606846975 - 500

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//
	//Case 6 Decrease Proposer Stake to Min
	//
	valChg = proposer.Copy()
	valChg.VotingPower = 1

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 7 Increase Proposer Stake to Outbid Next Highest Validator
	valSetTemp := valSet.Copy()
	valSetTemp.IncrementProposerPriority(1)

	valChg = valSetTemp.Proposer.Copy()
	valChg.VotingPower = 110
	proposerTemp := valChg

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", proposerTemp.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSetTemp, 10)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 8 Decrease Proposer Stake Below Lowest Validator
	valChg = proposer.Copy()
	valChg.VotingPower = 2

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 9 Increase Proposer Stake But Still Less Than Next Validator
	valSetTemp = valSet.Copy()
	valSetTemp.IncrementProposerPriority(2)

	valChg = valSetTemp.Proposer.Copy()
	valChg.VotingPower = 90
	proposerTemp = valChg

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", proposerTemp.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSetTemp)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSetTemp, 10)
	if newProposer.Signer != proposerTemp.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 10 Decrease Proposer Stake But Still More Than Next Validator
	valChg = proposer.Copy()
	valChg.VotingPower = 95

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 11 Proposer Exits
	valChg = proposer.Copy()
	valChg.VotingPower = 0

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != v2.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	//Case 12 Validator other than Proposer Exits
	valChg = v2.Copy()
	valChg.VotingPower = 0

	newProposer = UpdateValidatorSetWithoutRotation(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should not change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnExit(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeDecrease(t, valChg, valSet.Copy())
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

	newProposer = UpdateValidatorSetWithRotationOnStakeBelowThreshold(t, valChg, valSet.Copy(), 10)
	if newProposer.Signer != proposer.Signer {
		t.Errorf("Proposer should change OldProposer %v NewProposer%v", valSet.Proposer.ID, newProposer.ID)
	}

}

func UpdateValidatorSetWithoutRotation(t *testing.T, chgVal *Validator, valSet *ValidatorSet) *Validator {
	err := valSet.UpdateWithChangeSet([]*Validator{
		chgVal,
	})

	if err != nil {
		t.Error(err)
		return nil
	}

	return valSet.Proposer
}

func UpdateValidatorSetWithRotationOnExit(t *testing.T, chgVal *Validator, valSet *ValidatorSet) *Validator {
	err := valSet.UpdateWithChangeSet([]*Validator{
		chgVal,
	})

	if err != nil {
		t.Error(err)
		return nil
	}

	if valSet.Proposer.Signer == chgVal.Signer && chgVal.VotingPower == 0 {
		valSet.IncrementProposerPriority(1)
	}

	return valSet.Proposer
}

func UpdateValidatorSetWithRotationOnStakeDecrease(t *testing.T, chgVal *Validator, valSet *ValidatorSet) *Validator {
	var rotate bool

	if valSet.Proposer.Signer == chgVal.Signer && valSet.Proposer.VotingPower > chgVal.VotingPower {
		rotate = true
	}

	err := valSet.UpdateWithChangeSet([]*Validator{
		chgVal,
	})

	if err != nil {
		t.Error(err)
		return nil
	}

	if rotate {
		valSet.IncrementProposerPriority(1)
	}

	return valSet.Proposer
}

func UpdateValidatorSetWithRotationOnStakeBelowThreshold(t *testing.T, chgVal *Validator, valSet *ValidatorSet, threshold int64) *Validator {
	var rotate bool

	if valSet.Proposer.Signer == chgVal.Signer && valSet.Proposer.VotingPower > chgVal.VotingPower && chgVal.VotingPower < threshold {
		rotate = true
	}

	err := valSet.UpdateWithChangeSet([]*Validator{
		chgVal,
	})

	if err != nil {
		t.Error(err)
		return nil
	}

	if rotate {
		valSet.IncrementProposerPriority(1)
	}

	return valSet.Proposer
}
