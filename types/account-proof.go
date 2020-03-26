package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DividendAccount contains ID, merkle proof, leaf index in merketree
type DividendAccountProof struct {
	ID    DividendAccountID `json:"ID"`
	Proof []byte            `json:"accountProof"`
	Index uint64            `json:"index"`
}

func NewDividendAccountProof(id DividendAccountID, proof []byte, index uint64) DividendAccountProof {
	return DividendAccountProof{
		ID:    id,
		Proof: proof,
		Index: index,
	}
}

func (ap *DividendAccountProof) String() string {
	if ap == nil {
		return "nil-DividendAccountProof"
	}

	return fmt.Sprintf("DividendAccount{%v %v %v}",
		ap.ID,
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
