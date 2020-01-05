package server

import (
	"net/http"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-kit/kit/log"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/app"
	tx "github.com/maticnetwork/heimdall/client/tx"
	"github.com/maticnetwork/heimdall/helper"

	// unnamed import of statik for swagger UI support
	_ "github.com/maticnetwork/heimdall/server/statik"
)

// ServeCommands will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommands(cdc *codec.Codec, registerRoutesFn func(*lcd.RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			helper.InitHeimdallConfig("")

			rs := lcd.NewRestServer(cdc)
			registerRoutesFn(rs)
			logger := tmLog.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
			err := rs.Start(
				viper.GetString(client.FlagListenAddr),
				viper.GetInt(client.FlagMaxOpenConnections),
				0,
				0,
			)

			logger.Info("REST server started")
			return err
		},
	}
	cmd.Flags().String(client.FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	cmd.Flags().Bool(client.FlagTrustNode, true, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, "tcp://localhost:26657", "Address of the node to connect to")
	cmd.Flags().Int(client.FlagMaxOpenConnections, 1000, "The number of maximum open connections")

	return cmd
}

// RegisterRoutes register routes of all modules
func RegisterRoutes(rs *lcd.RestServer) {
	registerSwaggerUI(rs)

	rpc.RegisterRPCRoutes(rs.CliCtx, rs.Mux)
	tx.RegisterRoutes(rs.CliCtx, rs.Mux)

	// auth.RegisterRoutes(rs.CliCtx, rs.Mux)
	// bank.RegisterRoutes(rs.CliCtx, rs.Mux)

	// checkpoint.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// bor.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// clerk.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)

	// register rest routes
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	// list all paths
	// rs.Mux.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// 	t, err := route.GetPathTemplate()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	r, err := route.GetMethods()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	fmt.Println(strings.Join(r, ","), t)
	// 	return nil
	// })
}

func registerSwaggerUI(rs *lcd.RestServer) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	rs.Mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}
