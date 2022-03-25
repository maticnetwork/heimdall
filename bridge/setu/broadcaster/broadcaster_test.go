package broadcaster

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/maticnetwork/heimdall/app"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Parallel test - to check BroadcastToHeimdall syncronization
func TestBroadcastToHeimdall(t *testing.T) {
	t.Parallel()
	cdc := app.MakeCodec()
	// cli context

	tendermintNode := "http://localhost:26657"
	viper.Set(helper.NodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	// cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastSync
	// cliCtx.TrustNode = true

	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	_txBroadcaster := NewTxBroadcaster(cdc)

	testData := []checkpointTypes.MsgCheckpoint{
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 0, EndBlock: 63, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 0, EndBlock: 63, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 0, EndBlock: 63, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 0, EndBlock: 63, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
	}

	for index, test := range testData {
		t.Run(string(rune(index)), func(t *testing.T) {
			// create and send checkpoint message
			msg := checkpointTypes.NewMsgCheckpointBlock(
				test.Proposer,
				test.StartBlock,
				test.EndBlock,
				test.RootHash,
				test.AccountRootHash,
				"1234",
			)

			response, err := http.Get(_txBroadcaster.cliCtx.NodeURI + "/unconfirmed_txs")
			if err != nil {
				t.Fatal(err)
			}

			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			var Response struct {
				Result struct {
					Total string   `json:"total"`
					Txs   []string `json:"txs"`
				} `json:"result"`
			}

			err = json.Unmarshal(body, &Response)
			if err != nil {
				t.Fatal(err)
			}

			if len(Response.Result.Txs) != 0 {
				for _, txn := range Response.Result.Txs {
					txBytes, err := base64.StdEncoding.DecodeString(txn)
					if err != nil {
						t.Fatal(err)
					}
					decodedTxn, err := helper.GetTxDecoder(cdc)(txBytes)
					if err != nil {
						t.Fatal(err)
					}

					txnMsg := decodedTxn.GetMsgs()[0]
					if txnMsg.Type() == msg.Type() {

					}
				}
			}

			err = _txBroadcaster.BroadcastToHeimdall(msg)
			assert.Empty(t, err, "Error broadcasting tx to heimdall", err)
		})
	}
}
