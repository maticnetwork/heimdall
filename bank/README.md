# Bank Module

## Table of Contents

* [Overview](#overview)
* [How does it work](#how-does-it-work)
* [Query commands](#query-commands)
* [Transaction commands](#transaction-commands)

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

Once the event is validated by the Handler, It will go send coins a particular amount of coins to the sender

## Query commands

One can run the following query commands from the bank module :

* `balance` - Query for bank balance of an address.

### CLI commands

```
heimdallcli query bank balance [address]
```

### REST endpoints

```
curl -X GET "localhost:1317/bank/balances/{address}"
```

## Transaction commands

One can run the following transactions commands from the bank module :

* `send` - Send coin to an address.

### CLI commands

```
heimdallcli tx bank send [to_address] [amount]
```

### REST endpoints

```
curl -X POST "localhost:1317/bank/accounts/{address}/transfers"
```
