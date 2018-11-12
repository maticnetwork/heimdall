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

### Installation with Docker

**Run Docker container**

Create and run docker container with mounted directory -

```bash
$ docker run --name matic-heimdall -P -it \
    -v ~/.heimdalld:/root/.heimdalld \
    -v (pwd)/logs:/go/src/github.com/maticnetwork/heimdall/logs \
    "maticnetwork/heimdall:<tag-name>" \
    bash
```

Note: Do not forget to replace `<tag-name>` with actual tag-name.

**Initialize heimdall**

Once docker container is created and running you will be on container.

Run following command to initalize heimdall and create config files -

```bash
$ docker exec -it matic-heimdall bash
<docker-container>$ cd /go/src/github.com/maticnetwork/heimdall
<docker-container>$ make init-heimdall
```

**Create heimdall-config.json**

Create `~/.heimdalld/config/heimdall-config.json` directory with following content -

```json
{
  "mainRPCUrl": "https://kovan.infura.io",
  "maticRPCUrl": "https://testnet.matic.network",

  "stakeManagerAddress": "8b28d78eb59c323867c43b4ab8d06e0f1efa1573",
  "rootchainAddress": "e022d867085b1617dc9fb04b474c4de580dccf1a"
}
```

You can check your address and public key with following command:

```bash
$ docker exec -it matic-heimdall sh -c "make show-account-heimdall"
```

**Start heimdall**

Start heimdall from Docker container

```bash
$ docker exec -it matic-heimdall sh -c "make start-all"
```

### Propose new checkpoint

```
POST http://localhost:1317/checkpoint/new
Content-Type: application/json
Content-Length: length
Accept-Language: en-us
Connection: Keep-Alive

{
  "rootHash": "0xd494377d4439a844214b565e1c211ea7154ca300b98e3c296f19fc9ada36db33",
  "startBlock": 4733031,
  "endBlock": 4733034
}
```

**CURL command**

```bash
$ curl -X POST \
  http://localhost:1317/checkpoint/new \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
    "rootHash": "0x5ba1680c5f5d5da8c7e3c08ba5d168c69da7a7104cf4beab94f7c0c955551f35",
    "startBlock": 4733040,
    "endBlock": 4733050
  }'
```

**Note: You must have Ethers in your account while submitting checkpoint.**

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
