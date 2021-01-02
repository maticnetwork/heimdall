package cmd

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	ethcrypto "github.com/maticnetwork/bor/crypto"
)

// HeimdallPubkeyCmd debug heimdall pub key
func HeimdallPubkeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "heimdall-pubkey [pubkey]",
		Short: "Debug pubkey",
		Long: fmt.Sprintf(`Debug pubkey

Example:
$ %s debug pubkey AxgWnTKXilBxQFfhKYuzdePur084I7BSkU+gIXVerusZ
			`, version.AppName),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// clientCtx := client.GetClientContextFromCmd(cmd)

			pk, err := getPubKeyFromString(args[0])
			if err != nil {
				return err
			}

			var compressed []byte = pk
			var uncompressed *ecdsa.PublicKey

			// if uncompressed, get compressed
			if len(pk) == 65 {
				uncompressed, err = ethcrypto.UnmarshalPubkey(pk)
				if err != nil {
					panic(err)
				}
				compressed = ethcrypto.CompressPubkey(uncompressed)
			} else {
				uncompressed, err = ethcrypto.DecompressPubkey(pk)
				if err != nil {
					panic(err)
				}
			}

			// marshal uncompressed
			uncompressedBytes := ethcrypto.FromECDSAPub(uncompressed)

			// secp256k1Obj
			secp256k1Obj := secp256k1.PubKey(compressed)

			cmd.Println("Address:", secp256k1Obj.Address())
			cmd.Println("Pubkey (base64):", string(base64.StdEncoding.EncodeToString(secp256k1Obj.Bytes())))
			cmd.Printf("Pubkey (compressed): %X\n", secp256k1Obj.Bytes())
			cmd.Printf("Pubkey (uncompressed): %X\n", uncompressedBytes)

			return nil
		},
	}
}

// getPubKeyFromString returns a Tendermint PubKey (ethsecp256k1) bytes by attempting
// to decode the pubkey string from hex and base64. If all
// encodings fail, an error is returned.
func getPubKeyFromString(pkstr string) ([]byte, error) {
	pkstr = strings.Replace(pkstr, "0x", "", 1)
	bz, err := hex.DecodeString(pkstr)
	if err == nil {
		return bz, nil
	}

	bz, err = base64.StdEncoding.DecodeString(pkstr)
	if err == nil {
		return bz, nil
	}

	return nil, fmt.Errorf("pubkey '%s' invalid; expected hex, base64, or bech32 of correct size", pkstr)
}
