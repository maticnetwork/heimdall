package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/tendermint/tendermint/crypto"
	// "github.com/tendermint/tendermint/crypto/secp256k1"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*crypto.PrivKey)(nil), nil)

}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	// registry.RegisterInterface("tendermint.crypto.Pubkey", (*tmcrypto.PubKey)(nil))
	registry.RegisterInterface("crypto.Pubkey", (*crypto.PubKey)(nil))
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
