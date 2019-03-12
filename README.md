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

Start heimdall

```bash
$ make run-heimdall
```

Start rest-server

```$bash
$ make run-server
```

Start bridge

```$bash
$ make run-bridge
```

### Installation with Docker

**Run Docker container**

Create and run docker container with mounted directory -

```bash
$ docker run --name matic-heimdall -p 1317:1317 -p 26656:26656 -p 26657:26657 -it \
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

Run following command from host to initalize heimdall and create config files -

```bash
$ docker exec -it matic-heimdall sh -c "make init-heimdall"
```

**Modify heimdall-config.json**

Modify `~/.heimdalld/config/heimdall-config.json` file with latest contract addresses and URL's

> Example heimdall-config.json

```json
{
  "mainRPCUrl": "https://ropsten.infura.io",
  "maticRPCUrl": "https://testnet.matic.network",
  "stakeManagerAddress": "0xd0d82149efb003eb8afd602a3c3a1532898ea1af",
  "rootchainAddress": "0x4463d704416dccf1781231c484e2aedd7dc9da43",
  "childBlockInterval": 10000,
  "checkpointerPollInterval": 60000,
  "syncerPollInterval": 30000,
  "noackPollInterval": 15000000000,
  "avgCheckpointLength": 256,
  "maxCheckpointLength": 1024,
  "noackWaitTime": 300000000000,
  "checkpointBufferTime": 256000000000
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
$ make tests
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
