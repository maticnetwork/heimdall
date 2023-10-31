# Bank Module

## Table of Contents

* [Overview](#overview)
* [How does it work](#how-does-it-work)
* [How to send coins](#how-to-send-coins)
* [Query commands](#query-commands)

## Overview

The bank module is responsible for handling multi-asset coin transfers between accounts. It exposes several interfaces with varying capabilities for secure interaction with other modules which must alter user balances.

## How does it work

```
type MsgSend struct {
	FromAddress types.HeimdallAddress `json:"from_address"`
	ToAddress   types.HeimdallAddress `json:"to_address"`
	Amount      sdk.Coins             `json:"amount"`
}
```

[Handler](handler.go) for this transaction validates whether send is enabled or not

Once the event is validated by the Handler, it will send a particular amount of coins to the sender

## How to send coins

One can run the following transactions commands from the bank module :

* `send` - Send coin to an address.

### CLI commands

```
heimdallcli tx bank send [TO_ADDRESS] [AMOUNT] --chain-id <CHAIN_ID>
```

### REST endpoints

Rest endpoint creates a message which needs to be written to a file. Then the sender needs to sign and broadcast a transaction

```
curl -X POST http://localhost:1317/bank/accounts/<TO_ADDRESS>/transfers \
  -H 'Content-Type: application/json' \
  -d '{
	"base_req": {
		"chain_id": <CHAIN_ID>,
		"from": <FROM_ADDRESS>
	},
	"amount": [
		{
			"denom": "matic",
			"amount": <AMOUNT>
		}
	]
}' > <FILE>.json

heimdallcli tx sign <FILE>.json --chain-id <CHAIN_ID> > <FILE2>.json

heimdallcli tx broadcast <FILE2>.json
```

## Query commands

One can run the following query commands from the bank module :

* `balance` - Query for bank balance of an address.

### CLI commands

```
heimdallcli query bank balance [ADDRESS]
```

### REST endpoints

```
curl -X GET "localhost:1317/bank/balances/{ADDRESS}"
```
