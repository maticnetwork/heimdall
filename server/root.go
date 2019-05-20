package server

import (
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmLog "github.com/tendermint/tendermint/libs/log"
	checkpoint "github.com/maticnetwork/heimdall/checkpoint/rest"
	staking "github.com/maticnetwork/heimdall/staking/rest"
	"github.com/cosmos/cosmos-sdk/client/lcd"
)

// ServeCommands will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommands(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			rs := lcd.NewRestServer(cdc)
			createHandler(cdc)
			logger := tmLog.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
			err := rs.Start(viper.GetString(client.FlagListenAddr),
				viper.GetInt(client.FlagMaxOpenConnections))

			logger.Info("REST server started")
			return err
		},
	}
	cmd.Flags().String(client.FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	cmd.Flags().String(client.FlagCORS, "", "Set the domains that can make CORS requests (* for all)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(client.FlagMaxOpenConnections, 1000, "The number of maximum open connections")

	return cmd
}

func createHandler(cdc *codec.Codec) http.Handler {
	r := mux.NewRouter()

	cliCtx := context.NewCLIContext().WithCodec(cdc)
	cliCtx.TrustNode = true
	//keys.RegisterRoutes(r, true)
	rpc.RegisterRoutes(cliCtx, r)
	tx.RegisterRoutes(cliCtx, r, cdc)

	// Addded rest commands to adding transction !
	checkpoint.RegisterRoutes(cliCtx, r, cdc)
	staking.RegisterRoutes(cliCtx, r, cdc)
	return r
}
