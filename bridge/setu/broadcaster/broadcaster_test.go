package broadcaster

import (
	"fmt"
	"os"
	"testing"

	httpClient "github.com/tendermint/tendermint/rpc/client/http"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
	checkpointTypes "github.com/maticnetwork/heimdall/x/checkpoint/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Parallel test - to check BroadcastToHeimdall synchronization
func TestBroadcastToHeimdall(t *testing.T) {
	//viperConfig := viper.New()
	t.Parallel()
	cdc, _ := app.MakeCodecs()
	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set("log_level", "info")

	viper.Set(flags.FlagNode, tendermintNode)
	viper.Set(flags.FlagHome, os.ExpandEnv("$HOME/.heimdalld"))

	rootDir := viper.GetString(flags.FlagHome)
	fmt.Println("rootdir ", rootDir)
	//cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	//client.Con
	//cliCtx.BroadcastMode = client.BroadcastSync
	//cliCtx.TrustNode = true
	cdc, _ = app.MakeCodecs()
	// encoding
	encoding := app.MakeEncodingConfig()
	// cli context
	cliCtx := client.Context{}.WithJSONMarshaler(cdc)
	chainID := helper.GetGenesisDoc().ChainID
	_httpClient, _ := httpClient.New(helper.GetConfig().TendermintRPCUrl, "/websocket")

	err := helper.InitHeimdallConfig()
	require.NoError(t, err)

	cliCtx = cliCtx.WithNodeURI(helper.GetConfig().TendermintRPCUrl).
		WithClient(_httpClient).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithInterfaceRegistry(encoding.InterfaceRegistry).
		WithTxConfig(encoding.TxConfig).
		WithFromAddress(helper.GetAddress()).
		WithChainID(chainID).
		WithSkipConfirmation(true)

	cliCtx.BroadcastMode = flags.BroadcastAsync

	// *pflag.FlagSet

	cmd := cobra.Command{}

	_txBroadcaster := NewTxBroadcaster(cliCtx, cdc, cmd.Flags())

	testData := []checkpointTypes.MsgCheckpoint{
		{
			Proposer: helper.GetAddressStr(), StartBlock: 0, EndBlock: 63,
			RootHash:        hmCommon.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e").String(),
			AccountRootHash: hmCommon.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d").String(),
		},
		{
			Proposer: helper.GetAddressStr(), StartBlock: 64, EndBlock: 1024,
			RootHash:        hmCommon.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e").String(),
			AccountRootHash: hmCommon.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d").String(),
		},
		{
			Proposer: helper.GetAddressStr(), StartBlock: 1025, EndBlock: 2048,
			RootHash:        hmCommon.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e").String(),
			AccountRootHash: hmCommon.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d").String(),
		},
		{
			Proposer: helper.GetAddressStr(), StartBlock: 2049, EndBlock: 3124,
			RootHash:        hmCommon.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e").String(),
			AccountRootHash: hmCommon.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d").String(),
		},
	}

	for index, test := range testData {
		t.Run(string(rune(index)), func(t *testing.T) {
			// create and send checkpoint message
			accAddr, err := sdk.AccAddressFromHex(test.Proposer)
			require.NoError(t, err)
			msg := checkpointTypes.NewMsgCheckpointBlock(
				accAddr,
				test.StartBlock,
				test.EndBlock,
				hmCommon.HexToHeimdallHash(test.RootHash),
				hmCommon.HexToHeimdallHash(test.AccountRootHash),
				"15001",
			)

			err = _txBroadcaster.BroadcastToHeimdall(&msg)
			fmt.Println("Err is ", err)
			assert.Empty(t, err, "Error broadcasting tx to heimdall", err)
		})
	}
}
