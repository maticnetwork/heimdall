package server

import (
	ctx "context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/go-kit/log"
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
	hmRest "github.com/maticnetwork/heimdall/types/rest"

	// unnamed import of statik for swagger UI support
	"github.com/maticnetwork/heimdall/server/gRPC"
	_ "github.com/maticnetwork/heimdall/server/statik"
)

const shutdownTimeout = 10 * time.Second
const FlagGrpcAddr = "grpc-addr"
const FlagRPCReadHeaderTimeout = "read-header-timeout"

func StartRestServer(mainCtx ctx.Context, cdc *codec.Codec, registerRoutesFn func(ctx client.CLIContext, mux *mux.Router), restCh chan struct{}) error {
	// init vars for the Light Client Rest server
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	router := mux.NewRouter()
	logsWriter := helper.GetLogsWriter(helper.GetConfig().LogsWriterFile)
	logger := tmLog.NewTMLogger(log.NewSyncWriter(logsWriter)).With("module", "rest-server")

	registerRoutesFn(cliCtx, router)

	// server configuration
	cfg := rpcserver.DefaultConfig()
	cfg.MaxOpenConnections = viper.GetInt(client.FlagMaxOpenConnections)

	if viper.GetUint(client.FlagRPCReadTimeout) != 0 {
		readTimeOut := viper.GetUint(client.FlagRPCReadTimeout)
		cfg.ReadTimeout = time.Duration(readTimeOut) * time.Second
	}

	if viper.GetUint(client.FlagRPCWriteTimeout) != 0 {
		writeTimeOut := viper.GetUint(client.FlagRPCWriteTimeout)
		cfg.WriteTimeout = time.Duration(writeTimeOut) * time.Second
	}

	listenAddr := viper.GetString(client.FlagListenAddr)

	// this uses net.Listener underneath
	// which doesn't block, it runs in background
	// in other means it simply spawns a socket connection in OS level
	// and returns with the details we use to proxy orders to that socket
	listener, err := rpcserver.Listen(listenAddr, cfg)
	if err != nil {
		logger.Error("RPC could not listen: %v ", err)
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

	// Setup gRPC server
	gRPCLogger := tmLog.NewTMLogger(log.NewSyncWriter(logsWriter)).With("module", "gRPC-server")
	if err := gRPC.SetupGRPCServer(mainCtx, cdc, viper.GetString(FlagGrpcAddr), gRPCLogger); err != nil {
		return err
	}

	g.Go(func() error {
		// wait for os interrupt, then close Listener
		<-gCtx.Done()
		logger.Info("Shutting down heimdall rest server...")
		return listener.Close()
	})
	// wait here
	if err := g.Wait(); err != nil {
		if err != http.ErrServerClosed {
			logger.Error("Cannot start REST server.", "Error", err)
			return err
		}
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
		rww := &rpcserver.ResponseWriterWrapper{
			Status:         -1,
			ResponseWriter: w,
		}
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

	readHeaderTimeout := viper.GetUint(FlagRPCReadHeaderTimeout)
	if readHeaderTimeout == 0 {
		readHeaderTimeout = uint(cfg.ReadTimeout)
	}

	s := &http.Server{
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			select {
			case <-ctx.Done():
				fmt.Println("graceful handler exit")
				rw.WriteHeader(http.StatusOK)
				return
			default:
				recoverHandler(rw, r)
			}

		}),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: time.Duration(readHeaderTimeout) * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
		BaseContext: func(_ net.Listener) ctx.Context {
			return shutdownCtx
		},
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		return s.Serve(listener)
	})

	g.Go(func() error {
		// wait for interrupt signal coming from mainCtx
		// and then go to server shutdown
		<-shutdownCtx.Done()
		ctx, cancel := ctx.WithTimeout(ctx.Background(), shutdownTimeout)
		defer cancel()

		// nolint: contextcheck
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

	DecorateWithRestFlags(cmd)

	return cmd
}

// function is called whenever is the reste server flags has to be added to command
func DecorateWithRestFlags(cmd *cobra.Command) {
	cmd.Flags().String(client.FlagListenAddr, "tcp://0.0.0.0:1317", "The address for the server to listen on")
	cmd.Flags().Bool(client.FlagTrustNode, true, "Trust connected full node (don't verify proofs for responses)")
	cmd.Flags().Int(client.FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Uint(client.FlagRPCReadTimeout, 10, "The RPC read timeout (in seconds)")
	cmd.Flags().Uint(FlagRPCReadHeaderTimeout, 10, "The RPC header read timeout (in seconds)")
	cmd.Flags().Uint(client.FlagRPCWriteTimeout, 10, "The RPC write timeout (in seconds)")
	// heimdall specific flags for rest server start
	cmd.Flags().String(client.FlagChainID, "", "The chain ID to connect to")
	cmd.Flags().String(client.FlagNode, helper.DefaultTendermintNode, "Address of the node to connect to")
	// heimdall specific flags for gRPC server start
	cmd.Flags().String(FlagGrpcAddr, "0.0.0.0:3132", "The address for the gRPC server to listen on")
}

// RegisterRoutes register routes of all modules
func RegisterRoutes(ctx client.CLIContext, mux *mux.Router) {
	registerSwaggerUI(mux)

	rpc.RegisterRPCRoutes(ctx, mux)
	tx.RegisterRoutes(ctx, mux)

	// Register the status endpoint here (as it's generic)
	mux.HandleFunc("/status", statusHandlerFn(ctx)).Methods("GET")

	// auth.RegisterRoutes(rs.CliCtx, rs.Mux)
	// bank.RegisterRoutes(rs.CliCtx, rs.Mux)

	// checkpoint.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// staking.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// bor.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)
	// clerk.RegisterRoutes(rs.CliCtx, rs.Mux, rs.Cdc)

	// register rest routes
	app.ModuleBasics.RegisterRESTRoutes(ctx, mux)
}

func registerSwaggerUI(mux *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	mux.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", staticServer))
}

func statusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := cliCtx.Client.Status()
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, status.SyncInfo)
	}
}
