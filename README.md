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
$ make run-docker-develop
```

**Initialize heimdall**

Once docker container is created and running you will be on container.

Run following command to initalize heimdall and create config files -

```bash
<docker-container>$ make init-heimdall
```

You can check your address and public key with following command (mounted in `~/.heimdalld` directory):

```bash
cat ~/.heimdalld/config/priv_validator.json
```

**Create heimdall-config.json**

Create `~/.heimdalld/config/heimdall-config.json` directory with following content -

```json
{
  "main_rpcurl": "https://kovan.infura.io",
  "matic_rpcurl": "https://testnet.matic.network",

  "stakemanager_address": "8b28d78eb59c323867c43b4ab8d06e0f1efa1573",
  "rootchain_address": "e022d867085b1617dc9fb04b474c4de580dccf1a",
  "stakemanager_address": "74aaffd9b2e6d1e9f9913fd1f6f93614d4a1108a",
  "priv_validator_path": "/root/.heimdalld/config/priv_validator.json",

  "tendermint_endpoint": "http://127.0.0.1:26657"
}
```

**Start heimdall**

Start heimdall from Docker container

```bash
<docker-container>$ $ make start
```

### Propose new checkpoint

```
POST http://localhost:1317/checkpoint/new
Content-Type: application/json
Content-Length: length
Accept-Language: en-us
Connection: Keep-Alive

{
  "root_hash": "0xd494377d4439a844214b565e1c211ea7154ca300b98e3c296f19fc9ada36db33",
  "start_block": 4733031,
  "end_block": 4733034,
  "proposer_address": "0x84f8a67E4d16bb05aBCa3d154091566921e0B5e9"
}
```

**Note: You must have Ethers in your account while submitting checkpoint.**

### Docker

**Build docker**

```bash
$ make build-docker
```

**Build docker for develop branch**

```bash
$ make build-docker-develop
```

### License

GNU General Public License v3.0
