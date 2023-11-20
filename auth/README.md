# Auth Module

## Table of Contents

* [Overview](#overview)
    * [Gas and Fees](#gas-and-fees)
    * [Types](#types)
    * [Parameters](#parameters)
* [Query Commands](#query-commands)

## Overview

The auth module is responsible for specifying the base transaction and account types for an application. It contains the ante handler, where all basic transaction validity checks (signatures, nonces, auxiliary fields) are performed, and exposes the account keeper, which allows other modules to read, write, and modify accounts.

### Gas and Fees

Fees serve two purposes for an operator of the network.

Fees limit the growth of the state stored by every full node and allow for general purpose censorship of transactions of little economic value. Fees are best suited as an anti-spam mechanism where validators are disinterested in the use of the network and identities of users.

Since Heimdall doesn't support custom contract or code for any transaction, it uses fixed cost transactions. For fixed cost transactions, the validator can top up their accounts on the Ethereum chain and get tokens on Heimdall using the Topup module.

### Types

Besides accounts (specified in State), the types exposed by the auth module are StdSignature, the combination of an optional public key and a cryptographic signature as a byte array, StdTx, a struct that implements the sdk.Tx interface using StdSignature, and StdSignDoc, a replay-prevention structure for StdTx which transaction senders must sign over.

#### StdSignature

A StdSignature is the types of a byte array.
```
// StdSignature represents a sig
type StdSignature []byte
```

#### StdTx

A StdTx is a struct that implements the sdk.Tx interface, and is likely to be generic enough to serve the purposes of many types of transactions.

```
type StdTx struct {
        Msg       sdk.Msg      `json:"msg" yaml:"msg"`
        Signature StdSignature `json:"signature" yaml:"signature"`
        Memo      string       `json:"memo" yaml:"memo"`
}
```

#### StdSignDoc

A StdSignDoc is a replay-prevention structure to be signed over, which ensures that any submitted transaction (which is simply a signature over a particular byte string) will only be executable once on a Heimdall.

```
// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
    ChainID       string          `json:"chain_id" yaml:"chain_id"`
    AccountNumber uint64          `json:"account_number" yaml:"account_number"`
    Sequence      uint64          `json:"sequence" yaml:"sequence"`
    Msg           json.RawMessage `json:"msg" yaml:"msg"`
    Memo          string          `json:"memo" yaml:"memo"`
}
```

#### Account

It manages addresses, coins and nonce for transactions. It also signs and validates transactions.

```
type BaseAccount struct {
        Address types.HeimdallAddress `json:"address" yaml:"address"`
        Coins types.Coins `json:"coins" yaml:"coins"`
        PubKey crypto.PubKey `json:"public_key" yaml:"public_key"`
        AccountNumber uint64 `json:"account_number" yaml:"account_number"`
        Sequence uint64 `json:"sequence" yaml:"sequence"`
}
```

### Parameters

The auth module contains the following parameters:

|Key                   |Type  |Default value     |
|----------------------|------|------------------|
|MaxMemoCharacters     |uint64|256               |
|TxSigLimit            |uint64|7                 |
|TxSizeCostPerByte     |uint64|10                |
|SigVerifyCostED25519  |uint64|590               |
|SigVerifyCostSecp256k1|uint64|1000              |
|DefaultMaxTxGas       |uint64|1000000           |
|DefaultTxFees         |string|"1000000000000000"|

## Query Commands

One can run the following query command from the auth module.

- `account` - Query account details of a given address
- `params` - Query auth module parameters

To know your account details, run the following command:

```
heimdalld show-account
```
Expected response:
```
{
    "address": "0x68243159a498cf20d945cf3E4250918278BA538E",
    "pub_key": "0x040a9f6879c7cdab7ecc67e157cda15e8b2ddbde107a04bc22d02f50032e393f6360a05e85c7c1ecd201ad30dfb886af12dd02b47e4463f6f0f6f94159dc9f10b8"
}
```

### CLI Commands

```
heimdallcli query auth account <address> --trust-node
```

```
heimdallcli query auth params
```

### REST Endpoints

```
curl http://localhost:1317/auth/accounts/<address>
```

```
curl http://localhost:1317/auth/params
```
