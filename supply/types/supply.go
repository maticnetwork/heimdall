package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/maticnetwork/heimdall/types"
)

// Supply represents a struct that passively keeps track of the total supply amounts in the network
type Supply struct {
	Total types.Coins `json:"total" yaml:"total"` // total supply of tokens registered on the chain
}

// NewSupply creates a new Supply instance
func NewSupply(total types.Coins) Supply { return Supply{total} }

// DefaultSupply creates an empty Supply
func DefaultSupply() Supply { return NewSupply(types.NewCoins()) }

// Inflate adds coins to the total supply
func (supply *Supply) Inflate(amount types.Coins) {
	supply.Total = supply.Total.Add(amount)
}

// Deflate subtracts coins from the total supply
func (supply *Supply) Deflate(amount types.Coins) {
	supply.Total = supply.Total.Sub(amount)
}

// String returns a human readable string representation of a supplier.
func (supply Supply) String() string {
	b, err := yaml.Marshal(supply)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// ValidateBasic validates the Supply coins and returns error if invalid
func (supply Supply) ValidateBasic() error {
	if !supply.Total.IsValid() {
		return fmt.Errorf("invalid total supply: %s", supply.Total.String())
	}
	return nil
}
