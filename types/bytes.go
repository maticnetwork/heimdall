package types

import (
	"bytes"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// HexBytes the main purpose of HexBytes is to enable HEX-encoding for json/encoding.
type HexBytes []byte

// Equals returns boolean for whether two AccAddresses are Equal
func (bz HexBytes) Equals(bz2 HexBytes) bool {
	if bz.Empty() && bz2.Empty() {
		return true
	}

	return bytes.Equal(bz.Bytes(), bz2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty
func (bz HexBytes) Empty() bool {
	return len(bz) == 0
}

// Marshal needed for protobuf compatibility
func (bz HexBytes) Marshal() ([]byte, error) {
	return bz.Bytes(), nil
}

// Unmarshal needed for protobuf compatibility
func (bz *HexBytes) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

// MarshalJSON this is the point of Bytes.
func (bz HexBytes) MarshalJSON() ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(bz.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (bz HexBytes) MarshalYAML() (interface{}, error) {
	return bz.String(), nil
}

// UnmarshalJSON this is the point of Bytes.
func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := jsoniter.ConfigFastest.Unmarshal(data, &s); err != nil {
		return err
	}

	*bz = common.FromHex(s)

	return nil
}

// UnmarshalYAML unmarshals from YAML assuming Bech32 encoding.
func (bz *HexBytes) UnmarshalYAML(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}

	*bz = common.FromHex(s)

	return nil
}

// Bytes return bytes
func (bz HexBytes) Bytes() []byte {
	return bz
}

func (bz HexBytes) String() string {
	return hexutil.Encode(bz)
}

// Format format bytes
func (bz HexBytes) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(bz.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", bz)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", bz.Bytes())))
	}
}

//
// Utils
//

// BytesToHexBytes returns HexBytes with value b.
func BytesToHexBytes(b []byte) HexBytes {
	return HexBytes(b)
}

// HexToHexBytes returns hex bytes
func HexToHexBytes(b string) HexBytes {
	return HexBytes(common.FromHex(b))
}
