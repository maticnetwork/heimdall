package server

import (
	ctx "context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmLog "github.com/tendermint/tendermint/libs/log"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	"golang.org/x/sync/errgroup"

	"github.com/maticnetwork/heimdall/app"
	tx "github.com/maticnetwork/heimdall/client/tx"
	"github.com/maticnetwork/heimdall/helper"

	// unnamed import of statik for swagger UI support
	_ "github.com/maticnetwork/heimdall/server/statik"
)

const shutdownTimeout = 10 * time.Second

type options struct {
	listenAddr   string
	maxOpen      int
	readTimeout  uint
	writeTimeout uint
}

func StartRestServer(mainCtx ctx.Context, cdc *codec.Codec, registerRoutesFn func(ctx client.CLIContext, mux *mux.Router), restCh chan struct{}) error {
	// init vars for the Light Client Rest server
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	router := mux.NewRouter()
	logger := tmLog.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "rest-server")
	registerRoutesFn(cliCtx, router)

	// server configuration
	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = viper.GetInt(client.FlagMaxOpenConnections)
	cfg.ReadTimeout = time.Duration(0) * time.Second
	cfg.WriteTimeout = time.Duration(0) * time.Second
	listenAddr := viper.GetString(client.FlagListenAddr)

	// this uses net.Listener underneath
	// which doesn't block, it runs in background
	// in other means it simply spawns a socket connection in OS level
	// and returns with the details we use to proxy orders to that socket
	listener, err := rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		// TODO: log here
		return err
	}
	// no err? -> signal here that server is open for business
	close(restCh)

	logger.Info(
		fmt.Sprintf(
			"Starting application REST service (chain-id: %q)...",
			viper.GetString(flags.FlagChainID),
		),
	)

	g, gCtx := errgroup.WithContext(mainCtx)
	// start serving
	g.Go(func() error {
		return startRPCServer(mainCtx, listener, router, logger, cfg)
	})

	g.Go(func() error {
		// wait for os interrupt, then close Listener
		<-gCtx.Done()
		return listener.Close()
	})
	// wait here
	if err := g.Wait(); err != nil {
		logger.Error("Cannot start REST server.", "Error", err)
		return err
	}

	return nil
}

// this is borrowed from maticnetwork rpcserver
type maxBytesHandler struct {
	h http.Handler
	n int64
}

func (h maxBytesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, h.n)
	h.h.ServeHTTP(w, r)
}

func recoverAndLog(handler http.Handler, logger tmLog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to remember the status
		rww := &rpcserver.ResponseWriterWrapper{-1, w}
		begin := time.Now()

		rww.Header().Set("X-Server-Time", fmt.Sprintf("%v", begin.Unix()))

		defer func() {
			// Send a 500 error if a panic happens during a handler.
			// Without this, Chrome & Firefox were retrying aborted ajax requests,
			// at least to my localhost.
			if e := recover(); e != nil {

				// If RPCResponse
				if res, ok := e.(rpctypes.RPCResponse); ok {
					rpcserver.WriteRPCResponseHTTP(rww, res)
				} else {
					// For the rest,
					logger.Error(
						"Panic in RPC HTTP handler", "err", e, "stack",
						string(debug.Stack()),
					)
					rpcserver.WriteRPCResponseHTTPError(rww, http.StatusInternalServerError, rpctypes.RPCInternalError(rpctypes.JSONRPCStringID(""), e.(error)))
				}
			}

			// Finally, log.
			durationMS := time.Since(begin).Nanoseconds() / 1000000
			if rww.Status == -1 {
				rww.Status = 200
			}
			logger.Info("Served RPC HTTP response",
				"method", r.Method, "url", r.URL,
				"status", rww.Status, "duration", durationMS,
				"remoteAddr", r.RemoteAddr,
			)
		}()

		handler.ServeHTTP(rww, r)
	}
}

func startRPCServer(shutdownCtx ctx.Context, listener net.Listener, handler http.Handler, logger tmLog.Logger, cfg *rpcserver.Config) error {
	logger.Info(fmt.Sprintf("Starting RPC HTTP server on %s", listener.Addr()))
	recoverHandler := recoverAndLog(maxBytesHandler{h: handler, n: cfg.MaxBodyBytes}, logger)

	s := &http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			for {
				select {
				case <-ctx.Done():
					fmt.Println("gracefull handler exit")
					rw.WriteHeader(http.StatusOK)
					return
				default:
					recoverHandler(rw, r)
				}
			}
		}),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
		BaseContext: func(_ net.Listener) ctx.Context {
			return shutdownCtx
		},
	}

	g, gCtx := errgroup.WithContext(shutdownCtx)
	g.Go(func() error {
		return s.Serve(listener)
	})

	g.Go(func() error {
		// wait for interrupt signal comming from mainCtx
		// and then go to server shutdown
		<-gCtx.Done()
		ctx, cancel := ctx.WithTimeout(ctx.Background(), shutdownTimeout)
		defer cancel()

		return s.Shutdown(ctx)
	})

	if err := g.Wait(); err != nil {
		logger.Info("RPC HTTP server stopped", "err", err)
		return err
	}

	return nil
}

// ServeCommands will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommands(shutdownCtx ctx.Context, cdc *codec.Codec, registerRoutesFn func(ctx client.CLIContext, mux *mux.Router)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			helper.InitHeimdallConfig("")
			restCh := make(chan struct{}, 1)
			err := StartRestServer(shutdownCtx, cdc, registerRoutesFn, restCh)
			return err
		},
	}
	cmd.Flags().String(client.FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	cmd.Flags().Bool(client.FlagTrustNode, true, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, helper.DefaultTendermintNode, "Address of the node to connect to")
	cmd.Flags().Int(client.FlagMaxOpenConnections, 1000, "The number of maximum open connections")

	return cmd
}

// RegisterRoutes register routes of all modules
func RegisterRoutes(ctx client.CLIContext, mux *mux.Router) {
	registerSwaggerUI(mux)

	rpc.RegisterRPCRoutes(ctx, mux)
	tx.RegisterRoutes(ctx, mux)

	// auth.RegisterRoutes(rs.CliCtx, rs.Mux)
	// bank.RegisterRoutes(rs.CliCtx, rs.Mux)

	// checkpoint.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// bor.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// clerk.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)

	// register rest routes
	app.ModuleBasics.RegisterRESTRoutes(ctx, mux)

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

func registerSwaggerUI(mux *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}
	staticServer := http.FileServer(statikFS)
	mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}

// Check locally if rest server port has been opened
func restServerHealthCheck(restCh chan struct{}) {
	address := viper.GetString(client.FlagListenAddr)
	for {
		conn, err := net.Dial("tcp", address[6:])
		if err != nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		if conn != nil {
			defer conn.Close()
		}

		close(restCh)
		break
	}
}
