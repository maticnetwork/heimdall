package types

const (
	// AddrLen defines a valid address length
	AddrLen = 20
	// MainPrefix defines the  prefix main prefix
	MainPrefix = "matic"

	// PrefixAccount is the prefix for account keys
	PrefixAccount = "acc"
	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// PrefixAddress is the prefix for addresses
	PrefixAddress = "addr"

	// PrefixAccAddr defines the  prefix of an account's address
	PrefixAccAddr = MainPrefix
	// PrefixAccPub defines the  prefix of an account's public key
	PrefixAccPub = MainPrefix + PrefixPublic
	// PrefixValAddr defines the  prefix of a validator's operator address
	PrefixValAddr = MainPrefix + PrefixValidator + PrefixOperator
	// PrefixValPub defines the  prefix of a validator's operator public key
	PrefixValPub = MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// PrefixConsAddr defines the  prefix of a consensus node address
	PrefixConsAddr = MainPrefix + PrefixValidator + PrefixConsensus
	// PrefixConsPub defines the  prefix of a consensus node public key
	PrefixConsPub = MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)
