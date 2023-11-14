# Topup

Heimdall Topup is an amount which will be used to pay fees on Heimdall chain.

There are two ways to topup your account:

1. When new validator joins, they can mention a `topup` amount as top-up in addition to the staked amount, which will be moved as balance on Heimdall chain to pays fees on Heimdall.
2. A user can directly call the top-up function on the staking smart contract on Ethereum to increase top-up balance on Heimdall.

## Messages

### MsgTopup

`MsgTopup` transaction is responsible for minting balance to an address on Heimdall based on Ethereum chain's `TopUpEvent` on staking manager contract.

Handler for this transaction processes top-up and increases the balance only once for any given `msg.TxHash` and `msg.LogIndex`. It throws `Older invalid tx found` error, if trying to process the top-up more than once.

Here is the structure for the top-up transaction message:

```go
type MsgTopup struct {
	FromAddress types.HeimdallAddress `json:"from_address"`
	User        types.HeimdallAddress `json:"user"`
	Fee         sdk.Int               `json:"fee"`
	TxHash      types.HeimdallHash    `json:"tx_hash"`
	LogIndex    uint64                `json:"log_index"`
	BlockNumber uint64                `json:"block_number"`
}
```

### MsgWithdrawFee

`MsgWithdrawFee` transaction is responsible for withdrawing balance from Heimdall to Ethereum chain. A Validator can withdraw any amount from Heimdall.

Handler processes the withdraw by deducting the balance from the given validator and prepares the state to send the next checkpoint. The next possible checkpoint will contain the withdraw related state for the specific validator.

Handler gets validator information based on `ValidatorAddress` and processes the withdraw. 

```go
// MsgWithdrawFee - high-level transaction of the fee coin withdrawal module
type MsgWithdrawFee struct {
	UserAddress types.HeimdallAddress `json:"from_address"`
	Amount      sdk.Int               `json:"amount"`
}
```

## CLI Commands

### Topup fee

```bash
heimdallcli tx topup fee --fee-amount <fee-amount> --log-index <log-index>  --tx-hash <transaction-hash> --user <validator ID> --block-number <block-number>
```

### Withdraw fee

```bash
heimdallcli tx topup withdraw --amount=<withdraw-amount>
```

To check reflected topup on account run following command

```bash
heimdallcli query auth account <validator-address> --trust-node
```

## REST APIs

### Topup fee

```bash
curl -X POST "http://localhost/topup/fee" -H "accept: application/json" -d "{
  "block_number": 0,
  "fee": "string",
  "log_index": 0,
  "tx_hash": "string",
  "user": "string"
}"
```

### Withdraw fee

```bash
curl -X POST "http://localhost/topup/withdraw" -H "accept: application/json" -d "{
  "amount": "string",
}"
```