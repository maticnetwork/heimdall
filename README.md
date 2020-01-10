# Heimdall

Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source

Make sure your have go1.11+ already installed

### Install

```bash
$ make install
```

### Run-heimdall

```bash
$ heimdalld start
```

### Run rest server

```bash
$ heimdalld rest-server
```

### Documentation

Latest docs are [here](https://docs.matic.network/)

### Sign raw tx using REST APIs

Here is an example to sign `send` tx:

```bash
$ curl -X POST http://localhost:1317/bank/accounts/0xc6B47758074Baa00493514102fA63f91A5bC20BB/transfers \
  -H 'Content-Type: application/json' \
  -d '{
	"base_req": {
		"chain_id": "heimdall-4qiRRM",
		"from": "0x6c468CF8c9879006E22EC4029696E005C2319C9D"
	},
	"amount": [
		{
			"denom": "vetic",
			"amount": "10"
		}
	]
}' > raw.json

$ ./build/heimdallcli tx sign raw.json --chain-id heimdall-P5rXwg > signed.json
$ ./build/heimdallcli tx broadcast signed.json
```

To encode signed json tx:

```bash
$ ./build/heimdallcli tx encode signed.json
```
