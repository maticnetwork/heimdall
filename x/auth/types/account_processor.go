package types

import (
	"github.com/maticnetwork/heimdall/x/auth/exported"
)

// AccountProcessor is an interface to process account as per module
type AccountProcessor func(*GenesisAccount, *BaseAccount) exported.Account
