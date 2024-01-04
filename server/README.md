# Server

## Overview

The `server` module in Heimdall is responsible for running a REST and gRPC server that exposes multiple endpoints for interacting with the Heimdall node. The REST module handles endpoints from every other module and also root endpoints. The gRPC server handles only some specific endpoints related to bor-heimdall communication. The implementation for the gRPC server is in the gRPC folder.

### REST

Every module in heimdall has their own rest endpoints, All of these endpoints are registered in the `server/rest.go` file via `app.ModuleBasics.RegisterRESTRoutes` and it also handle the root endpoints. The `server/rest.go` file also contains the `StartServer` function which starts the REST server.

### gRPC

The gRPC server is specifically used for communication between bor and heimdall. The implementation for the gRPC server is in the `server/grpc` folder. The `server/gRPC/gRPC.go` file contains the `StartServer` function which starts the gRPC server.

## Usage

To start the server, run the following command

```bash
heimdalld start-server
```

The `start-server` command is added into the heimdall binary and takes folowing flags

```bash
      --chain-id string            The chain ID to connect to
      --grpc-addr string           The address for the gRPC server to listen on (default "0.0.0.0:3132")
      --laddr string               The address for the server to listen on (default "tcp://0.0.0.0:1317")
      --max-open int               The number of maximum open connections (default 1000)
      --node string                Address of the node to connect to (default "tcp://localhost:26657")
      --read-header-timeout uint   The RPC header read timeout (in seconds) (default 10)
      --read-timeout uint          The RPC read timeout (in seconds) (default 10)
      --trust-node                 Trust connected full node (don't verify proofs for responses) (default true)
      --write-timeout uint         The RPC write timeout (in seconds) (default 10)
```
