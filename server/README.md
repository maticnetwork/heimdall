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

The `start-server` command is added into the heimdall binary and takes following flags

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
## Swagger UI

The REST server also exposes a swagger UI which can be used to interact with the REST endpoints. The swagger UI can be accessed at `http://localhost:1317/swagger-ui/` after starting the server.

Installation

1. go get github.com/rakyll/statik
2. go get -u github.com/go-swagger/go-swagger/cmd/swagger #For downloading the Go Swagger to create the spec using the swagger comments.


Steps to follow

1. Add the Swagger Comments to the API added or updated using documentation at https://goswagger.io/use/spec.html.
2. Run GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models  from the root directory.
3. cd maticnetwork/heimdall/server
4. Replace the Swagger.yaml file inside swagger-ui directory with the swagger.yaml newly generated in root directly in step 2
5. cd maticnetwork/heimdall/server && statik -src=./swagger-ui
6. cd maticnetwork/heimdall && make build
7. cd maticnetwork/heimdall && make run-server


Steps to follow for updated swagger-ui without using go-swagger

1. Add the Swagger Comments to the API added or updated using documentation at https://goswagger.io/use/spec.html.
2. Run GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models  from the root directory.
3. Copy zip file of source code from https://github.com/swagger-api/swagger-ui/releases
4. Unzip the zip file. Copy the contents of dist/ from the zip to heimdall/server/swagger-ui/
5. Convert the heimdall/server/swagger-ui/swagger.yaml to JSON format and place it in the same directory as the swagger.yaml file.
6. In heimdall/server/swagger-ui/swagger-initializer.js change `url: "./swagger.json"`,
7. cd maticnetwork/heimdall/server && statik -src=./swagger-ui

Visit http://localhost:1317/swagger-ui/ 

Reference
- https://github.com/rakyll/statik
