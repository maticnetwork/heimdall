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

  "stakeManagerAddress": "0xb4ee6879ba231824651991c8f0a34af4d6bfca6a",
  "rootchainAddress": "0x168ea52f1fafe28d584f94357383d4f6fa8a749a",
  "childBlockInterval": 10000
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
#### Generate new checkpoint
```$bash 
curl -X GET \
  http://localhost:1317/checkpoint/<startBlock>/<endBlock> \
  -H 'cache-control: no-cache' \
  -H 'postman-token: 9b98abc0-bdaa-772a-38d4-7c29351b87a4'
```

#### Propose new checkpoint

```
curl -X POST \
  http://localhost:1317/checkpoint/new \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: ccb714f8-7197-3f72-d4a3-97c6d654f53b' \
  -d '{
	"proposer":"0x0cdf0edd304a8e1715d5043d0afe3d3322cc6e3b",
	"rootHash":"0xb2160bdf78e2dd513763b7767d510cad84eba764b9ec0e00a0110fadaf14b179",
	"startBlock":201,
	"endBlock":204
}'
```

#### Submit ACK for checkpoint
```$bash 
curl -X POST \
  http://localhost:1317/checkpoint/ack \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -H 'postman-token: b800a79b-158b-8590-ba2f-aedb1fd1942a' \
  -d '{
	"HeaderBlock" : 20000
}'

```

**Note: You must have Ethers in your account while submitting checkpoint.**

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
