package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DividendAccountProof contains ID, merkle proof, leaf index in merketree
type DividendAccountProof struct {
	User  HeimdallAddress `json:"user"`
	Proof HexBytes        `json:"accountProof"`
	Index uint64          `json:"index"`
}

// NewDividendAccountProof generate proof for new dividend account
func NewDividendAccountProof(user HeimdallAddress, proof HexBytes, index uint64) DividendAccountProof {
	return DividendAccountProof{
		User:  user,
		Proof: proof,
		Index: index,
	}
}

func (ap *DividendAccountProof) String() string {
	if ap == nil {
		return "nil-DividendAccountProof"
	}

	return fmt.Sprintf("DividendAccount{%v %v %v}",
		ap.User,
		ap.Proof,
		ap.Index)
}

// MarshallDividendAccountProof - amino Marshall DividendAccountProof
func MarshallDividendAccountProof(cdc *codec.Codec, dividendAccountProof DividendAccountProof) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(dividendAccountProof)
	if err != nil {
		return bz, err
	}

	return bz, nil
}

// UnMarshallDividendAccountProof - amino Unmarshall DividendAccountProof
func UnMarshallDividendAccountProof(cdc *codec.Codec, value []byte) (DividendAccountProof, error) {

	var dividendAccountProof DividendAccountProof
	err := cdc.UnmarshalBinaryBare(value, &dividendAccountProof)
	if err != nil {
		return dividendAccountProof, err
	}
	return dividendAccountProof, nil
}
