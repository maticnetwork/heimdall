package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v2"
)

// Ensure that different address types implement the interface
var _ yaml.Marshaler = HeimdallHash{}

// HeimdallHash represents heimdall address
type HeimdallHash common.Hash

// ZeroHeimdallHash represents zero address
var ZeroHeimdallHash = HeimdallHash{}

// EthHash get eth hash
func (aa HeimdallHash) EthHash() common.Hash {
	return common.Hash(aa)
}

// Equals returns boolean for whether two HeimdallHash are Equal
func (aa HeimdallHash) Equals(aa2 HeimdallHash) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty
func (aa HeimdallHash) Empty() bool {
	return bytes.Equal(aa.Bytes(), ZeroHeimdallHash.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa HeimdallHash) Marshal() ([]byte, error) {
	return aa.Bytes(), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *HeimdallHash) Unmarshal(data []byte) error {
	*aa = HeimdallHash(common.BytesToHash(data))
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa HeimdallHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa HeimdallHash) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *HeimdallHash) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*aa = HexToHeimdallHash(s)
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *HeimdallHash) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*aa = HexToHeimdallHash(s)
	return nil
}

// Bytes returns the raw address bytes.
func (aa HeimdallHash) Bytes() []byte {
	return aa[:]
}

// String implements the Stringer interface.
func (aa HeimdallHash) String() string {
	if aa.Empty() {
		return ""
	}

	return "0x" + hex.EncodeToString(aa.Bytes())
}

// Hex returns hex string
func (aa HeimdallHash) Hex() string {
	return aa.String()
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa HeimdallHash) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", aa.Bytes())))
	}
}

//
// hash utils
//

// BytesToHeimdallHash returns Address with value b.
func BytesToHeimdallHash(b []byte) HeimdallHash {
	return HeimdallHash(common.BytesToHash(b))
}

// HexToHeimdallHash returns Address with value b.
func HexToHeimdallHash(b string) HeimdallHash {
	return HeimdallHash(common.HexToHash(b))
}
