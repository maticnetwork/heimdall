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
