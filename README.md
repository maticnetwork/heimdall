# Heimdall

Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Installation for development

Install all dependencies and tools

```bash
$ make dep
```

Build heimdall

```bash
$ make build
```

Start heimdall process

```bash
$ make run-heimdall
```

Start rest-server

```$bash
$ make rest-server
```

### Installation with Docker

**Run Docker container**

Create and run docker container with mounted directory -

```bash
$ docker run --name heimdall -p 1317:1317 -p 26656:26656 -p 26657:26657 -it \
    -v ~/.heimdalld:/root/.heimdalld \
    -v $(pwd)/logs:/go/src/github.com/maticnetwork/heimdall/logs \
    "maticnetwork/heimdall:<tag-name>" \
    bash
```

Note: Do not forget to replace `<tag-name>` with actual tag-name.

**Initialize heimdall**

Once docker container is created and running you will be on container.<br>
You can run make commands directly on the container.

`OR`

Run following commands from host to initalize heimdall and create config files -

```bash
$ docker exec -it matic-heimdall sh -c "make init-heimdall"
```

**Modify heimdall-config.json**

Modify `~/.heimdalld/config/heimdall-config.json` file with latest contract addresses and URL's like below -

```json
{
  "mainRPCUrl": "https://kovan.infura.io",
  "maticRPCUrl": "https://testnet.matic.network",

  "stakeManagerAddress": "0xb4ee6879ba231824651991c8f0a34af4d6bfca6a",
  "rootchainAddress": "0x168ea52f1fafe28d584f94357383d4f6fa8a749a",
  "childBlockInterval": 10000
}
```

You can check your address and public key with following command:

```bash
$ docker exec -it matic-heimdall sh -c "make show-account-heimdall"
```

You can also check your node ID with the following command:

```bash
$ docker exec -it matic-heimdall sh -c "make show-node-id"
```

**Adding Peers**

You can add peers separated by commas at `~/.heimdalld/config/config.toml` under `persistent_peers`
With the format `NodeID@IP:PORT` or `NodeID@DOMAIN:PORT`

**Start heimdall**

Start heimdall and other necessary services from host

```bash
$ docker exec -it matic-heimdall sh -c "make start-all"
```

Logs can be found under `./logs`

### Run Tests

You can run tests found in tests directory to make sure everything is working as expected after making changes

```$bash
$ go test -run <TestCaseName>/<SubTestName>
```

> Please add -v flag to see test logs

##### Example

```$bash
$ go test -v -run TestValUpdates/add
$ go test -v -run TestValidator
```

### Docker (Only for developers)

#### For develop

```bash
$ make build-docker-develop
```

#### For releases

```bash
$ make build-docker
```

**Push docker image to docker hub (Only for internal team)**

```bash
$ make push-docker
```

### License

GNU General Public License v3.0
