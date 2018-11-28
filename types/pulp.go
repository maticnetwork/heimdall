package types

import (
	"encoding/hex"
	"reflect"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	PulpHashLength int = 4
)

// Pulp codec for RLP
type Pulp struct {
	typeInfos map[string]reflect.Type
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
	p.typeInfos = make(map[string]reflect.Type)
	return p
}

// GetPulpHash returns string hash
func GetPulpHash(name string) []byte {
	return crypto.Keccak256([]byte(name))[:PulpHashLength]
}

// RegisterConcrete should be used to register concrete types that will appear in
// interface fields/elements to be encoded/decoded by pulp.
func (p *Pulp) RegisterConcrete(msg sdk.Msg) {
	rtype := reflect.TypeOf(msg)
	name := msg.Route()
	p.typeInfos[hex.EncodeToString(GetPulpHash(name))] = rtype
}

// GetMsgTxInstance get new instance associated with base tx
func (p *Pulp) GetMsgTxInstance(hash []byte) interface{} {
	rtype := p.typeInfos[hex.EncodeToString(hash[:PulpHashLength])]
	return reflect.New(rtype).Elem().Interface().(sdk.Msg)
}

// EncodeToBytes encodes msg to bytes
func (p *Pulp) EncodeToBytes(msg sdk.Msg) ([]byte, error) {
	name := msg.Route()
	txBytes, err := rlp.EncodeToBytes(msg)
	if err != nil {
		return nil, err
	}

	return append(GetPulpHash(name), txBytes[:]...), nil
}

// DecodeBytes decodes bytes to msg
func (p *Pulp) DecodeBytes(data []byte) (interface{}, error) {
	rtype := p.typeInfos[hex.EncodeToString(data[:PulpHashLength])]
	msg := reflect.New(rtype).Interface()
	err := rlp.DecodeBytes(data[PulpHashLength:], msg)
	if err != nil {
		return nil, err
	}

	// change pointer to non-pointer
	vptr := reflect.New(reflect.TypeOf(msg).Elem()).Elem()
	vptr.Set(reflect.ValueOf(msg).Elem())
	return vptr.Interface(), nil
}
