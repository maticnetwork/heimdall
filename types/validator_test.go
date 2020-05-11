package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// valInput struct is used to seed data for testing
// if the need arises it can be ported to the main build
type valInput struct {
	id         ValidatorID
	startEpoch uint64
	endEpoch   uint64
	nonce      uint64
	power      int64
	pubKey     PubKey
	signer     HeimdallAddress
}

func TestNewValidator(t *testing.T) {

	// valCase created so as to pass it to assertPanics func,
	// ideally would like to get rid of this and pass the function directly

	tc := []struct {
		in  valInput
		out *Validator
		msg string
	}{
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToHeimdallAddress([]byte("12345678909876543210")),
			},
			out: &Validator{Signer: BytesToHeimdallAddress([]byte("12345678909876543210"))},
			msg: "testing for exact HeimdallAddress",
		},
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToHeimdallAddress([]byte("1")),
			},
			out: &Validator{Signer: BytesToHeimdallAddress([]byte("1"))},
			msg: "testing for small HeimdallAddress",
		},
		{
			in: valInput{
				id:     ValidatorID(uint64(0)),
				signer: BytesToHeimdallAddress([]byte("123456789098765432101")),
			},
			out: &Validator{Signer: BytesToHeimdallAddress([]byte("123456789098765432101"))},
			msg: "testing for excessively long HeimdallAddress, max length is supposed to be 20",
		},
	}
	for _, c := range tc {
		out := NewValidator(c.in.id, c.in.startEpoch, c.in.endEpoch, 1, c.in.power, c.in.pubKey, c.in.signer)
		assert.Equal(t, c.out, out)
	}
}

// TestSortValidatorByAddress am populating only the signer as that is the only value used in sorting
func TestSortValidatorByAddress(t *testing.T) {
	tc := []struct {
		in  []Validator
		out []Validator
		msg string
	}{
		{
			in: []Validator{
				Validator{Signer: BytesToHeimdallAddress([]byte("3"))},
				Validator{Signer: BytesToHeimdallAddress([]byte("2"))},
				Validator{Signer: BytesToHeimdallAddress([]byte("1"))},
			},
			out: []Validator{
				Validator{Signer: BytesToHeimdallAddress([]byte("1"))},
				Validator{Signer: BytesToHeimdallAddress([]byte("2"))},
				Validator{Signer: BytesToHeimdallAddress([]byte("3"))},
			},
			msg: "reverse sorting of validator objects",
		},
	}
	for i, c := range tc {
		out := SortValidatorByAddress(c.in)
		assert.Equal(t, c.out, out, fmt.Sprintf("i: %v, case: %v", i, c.msg))
	}
}

func TestValidateBasic(t *testing.T) {
	neg1, uNeg1 := uint64(1), uint64(0)
	uNeg1 = uNeg1 - neg1
	tc := []struct {
		in  Validator
		out bool
		msg string
	}{
		{
			in:  Validator{StartEpoch: 1, EndEpoch: 5, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToHeimdallAddress([]byte("3"))},
			out: true,
			msg: "Valid basic validator test",
		},
		{
			in:  Validator{StartEpoch: 1, EndEpoch: 5, Nonce: 0, PubKey: NewPubKey([]byte("")), Signer: BytesToHeimdallAddress([]byte("3"))},
			out: false,
			msg: "Invalid PubKey \"\"",
		},
		{
			in:  Validator{StartEpoch: 1, EndEpoch: 5, Nonce: 0, PubKey: ZeroPubKey, Signer: BytesToHeimdallAddress([]byte("3"))},
			out: false,
			msg: "Invalid PubKey",
		},

		//		{
		//			in:  Validator{StartEpoch: uNeg1, EndEpoch: 5, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToHeimdallAddress([]byte("3"))},
		//			out: false,
		//			msg: "Invalid StartEpoch",
		//		},
		{
			// do we allow for endEpoch to be smaller than startEpoch ??
			in:  Validator{StartEpoch: 1, EndEpoch: uNeg1, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToHeimdallAddress([]byte("3"))},
			out: false,
			msg: "Invalid endEpoch",
		},
		{
			// in:  Validator{StartEpoch: 1, EndEpoch: 1, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: HeimdallAddress(BytesToHeimdallAddress([]byte(string(""))))},
			in:  Validator{StartEpoch: 1, EndEpoch: 1, Nonce: 0, PubKey: NewPubKey([]byte("nonZeroTestPubKey")), Signer: BytesToHeimdallAddress([]byte(""))},
			out: false,
			msg: "Invalid Signer",
		},
		{
			in:  Validator{},
			out: false,
			msg: "Valid basic validator test",
		},
	}
	for _, c := range tc {
		out := c.in.ValidateBasic()
		assert.Equal(t, c.out, out, c.msg)
	}
}
