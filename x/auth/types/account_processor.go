package types

// AccountProcessor is an interface to process account as per module
type AccountProcessor func(*GenesisAccount, *BaseAccount) AccountI
