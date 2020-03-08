package helper

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/types"
)

func TestUnpackSigAndVotes(t *testing.T) {
	// Signer Address List for below SubmitHeaderBlock Transaction Payload
	signerAddresses := []string{"a03d8f5af7413e4fd5a37fde9286e390ef8f3c07", "b1bf4473c6b1918a6e37408e1c14df81281411a8", "ba754e3893adb3cabc0afe7932b4b5a3cee3f3ab"}
	signatureLen := 65

	payload := "ec83d3ba000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000000000030ef8f6865696d64616c6c2d39337251774b84766f7465820cb0800294907eb68cd3480777e3fde8897fb1373de6e982cc0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c3c65be37c47302d110fbe6453ef84cc69ecc744d9dbd6fa98d68a621541e2e8291758b962f41ed31a1f8c25db2631da2f49fe20ff1e2d683dd0d802fcef6928d5000622073cfbc99994cd06d7a7a8b01e453b57495010d6eb312a68b00ca6f581d0729f494a385bd1fbe2fd8df3da706fa85a54694ab0d9a4177555048aaa7b3371005c6dd42d128482e603c5adc7cfca0f0c730b49d6bd8ba750307d497c21097c922479f51151c53b35f20a1a1cb4790afa95e470b8f70b0318726d6175b9055b340000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000045f84394b1bf4473c6b1918a6e37408e1c14df81281411a883543ff8835441f7a07d2842c3044740cfe1e1a5f782bfd3b91de0c634e9933524b5e3daacc854f49b845d94796f000000000000000000000000000000000000000000000000000000"
	decodedPayload, err := hex.DecodeString(payload)
	require.Empty(t, err, "Error while decoding payload")

	abi, err := getABI(string(rootchain.RootchainABI))
	require.Empty(t, err, "Error while getting RootChainABI")

	// Unpacking the Payload
	voteSignBytes, inputSigs, txData, err := UnpackSigAndVotes(decodedPayload, abi)
	require.Empty(t, err, "Error while unpacking payload")
	t.Log("voteSignBytes", hex.EncodeToString(voteSignBytes))
	t.Log("inputSigs", hex.EncodeToString(inputSigs))
	t.Log("txData", hex.EncodeToString(txData))
	t.Log("Signatures Count", len(inputSigs)/signatureLen)

	// Validating the Unpack Output
	for i, j := 0, 0; i < len(signerAddresses); i, j = i+1, j+signatureLen {
		pubKey, err := authTypes.RecoverPubkey(voteSignBytes, inputSigs[j:j+signatureLen])
		require.Empty(t, err, "Error while recovering pubkey from signature. voteSignBytes = %v, Signature=%v ", voteSignBytes, hex.EncodeToString(inputSigs[i:i+signatureLen]))
		pKey := types.NewPubKey(pubKey)
		signerAddress := pKey.Address().Bytes()
		t.Log("Pubkey Recovered", hex.EncodeToString(pubKey))
		t.Log("Signer Address", hex.EncodeToString(signerAddress))
		require.Equal(t, signerAddresses[i], hex.EncodeToString(signerAddress), "Signer Address Doesn't match")
	}
}

func TestGetPowerFromAmount(t *testing.T) {
	scenarios1 := map[string]string{
		"48000000000000000000000": "48000",
		"10000000000000000000000": "10000",
		"1000000000000000000000":  "1000",
		"4800000000000000000000":  "4800",
		"480000000000000000000":   "480",
		"20000000000000000000":    "20",
		"10000000000000000000":    "10",
		"1000000000000000000":     "1",
	}

	for k, v := range scenarios1 {
		bv, _ := big.NewInt(0).SetString(k, 10)
		p, err := GetPowerFromAmount(bv)
		require.Nil(t, err, "Error must be null for input %v, output %v", k, v)
		require.Equal(t, p.String(), v, "Power must match")
	}
}
