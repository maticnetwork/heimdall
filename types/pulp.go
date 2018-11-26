package types

import (
	"encoding/hex"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	PulpHashLength int = 4
)

// Pulp codec for RLP
type Pulp struct {
	typeInfos map[string]func() sdk.Msg
}

var once sync.Once
var pulp *Pulp

// GetPulpInstance gets new pulp codec
func GetPulpInstance() *Pulp {
	once.Do(func() {
		pulp = NewPulp()
	})
	return pulp
}

// NewPulp creates new pulp codec
func NewPulp() *Pulp {
	p := &Pulp{}
	p.typeInfos = make(map[string]func() sdk.Msg)
	return p
}

// GetPulpHash returns string hash
func GetPulpHash(name string) []byte {
	return crypto.Keccak256([]byte(name))[:PulpHashLength]
}

// RegisterConcrete should be used to register concrete types that will appear in
// interface fields/elements to be encoded/decoded by pulp.
func (p *Pulp) RegisterConcrete(val func() sdk.Msg) {
	name := reflect.TypeOf(val()).String()
	p.typeInfos[hex.EncodeToString(GetPulpHash(name))] = val
}

// GetMsgTxInstance get new instance associated with base tx
func (p *Pulp) GetMsgTxInstance(hash []byte) sdk.Msg {
	return p.typeInfos[hex.EncodeToString(hash[:PulpHashLength])]()
}

// EncodeToBytes encodes msg to bytes
func (p *Pulp) EncodeToBytes(msg sdk.Msg) ([]byte, error) {
	name := reflect.TypeOf(msg).String()
	txBytes, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return nil, err
	}

	return append(GetPulpHash("*"+name), txBytes[:]...), nil
}

// DecodeBytes decodes bytes to msg
func (p *Pulp) DecodeBytes(data []byte, msg sdk.Msg) error {
	err := rlp.DecodeBytes(data[PulpHashLength:], msg)
	if err != nil {
		return err
	}

	return nil
}
