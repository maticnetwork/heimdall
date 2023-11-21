# Governance Module

## Table of Contents

- [Overview](#overview)
- [Commands](#commands)
  - [Query Commands](#query-commands)
  - [Tx Commands](#tx-commands)
  - [Run Using CLI](#run-using-cli)
  - [Run Using REST](#run-using-rest)

## Overview

Heimdall governance works exactly the same as Cosmos-sdk's gov module. In this system, holders of the native staking token of the chain can vote on proposals on a 1 token = 1 vote basis. Here is a list of features the module currently supports:

- Proposal submission: Validators can submit proposals with a deposit. Once the minimum deposit is reached, proposal enters voting period. Valdiators that deposited on proposals can recover their deposits once the proposal is rejected or accepted.
- Vote: Validators can vote on proposals that reached MinDeposit.

There are deposit period and voting period as params in gov module. Minimum deposit has to be achieved before deposit period ends, otherwise proposal will be automatically rejected.

Once minimum deposits reached within deposit period, voting period starts. In voting period, all validators should vote their choices for the proposal. After voting period ends, gov/endblocker.go executes tally function and accepts or rejects proposal based on tally_params — quorum, threshold and veto.

There are different types of proposals that can be implemented in Heimdall. As of now, it supports only the Param change proposal.

### Param change proposal

Using this type of proposal, validators can change any params in any module of Heimdall.

Example: change minimum tx_fees for the transaction in auth module. When the proposal gets accepted, it automatically changes the params in Heimdall state. No extra TX is needed.

## Commands

One can run the following commands from the governance module:

### Query Commands

- `proposal` - Query details of a single proposal
- `proposals` - Query proposals with optional filters
- `vote` - Query details of a single vote
- `votes` - Query votes on a proposal
- `deposit` - Query details of a deposit
- `deposits` - Query deposits on a proposal
- `tally` - Get the tally of a proposal vote
- `params` - Query the parameters of the governance process
- `param` - Query the parameters (voting|tallying|deposit) of the governance process
- `proposer` - Query the proposer of a governance proposal

### Tx Commands

- `submit-proposal` - Submit a proposal along with an initial deposit
- `deposit` - Deposit tokens for an active proposal
- `vote` - Vote for an active proposal with options: yes/no/no_with_veto/abstain

### Run Using CLI

```
heimdallcli query gov proposal [proposal-id]
```
```
heimdallcli query gov proposals
```
```
heimdallcli query gov vote [proposal-id] [voter-id]
```
```
heimdallcli query gov votes [proposal-id]
```
```
heimdallcli query gov deposit [proposal-id] [depositer-addr]
```
```
heimdallcli query gov deposits [proposal-id]
```
```
heimdallcli query gov tally [proposal-id]
```
```
heimdallcli query gov params
```
```
heimdallcli query gov param [param-type]
```
```
heimdallcli query gov proposer [proposal-id]
```
```
heimdallcli tx gov submit-proposal --proposal proposal.json --from key
```
where `proposal.json` will have
``` 
{
  "title": "Auth Param Change",
  "description": "Update max tx gas",
  "changes": [
    {
      "subspace": "auth",
      "key": "MaxTxGas",
      "value": "2000000"
    }
  ],
  "deposit": [
    {
      "denom": "matic",
      "amount": "1000000000000000000"
    }
  ]
}
```
```
heimdallcli tx gov deposit [proposal-id] [stake]
```
```
heimdallcli tx gov vote [proposal-id] [option] --validator-id 1  --chain-id <heimdall-chain-id>
```
with options from `yes/no/no_with_veto/abstain`, Like
```
heimdallcli tx gov vote 1 "Yes" --validator-id 1  --chain-id <heimdal-chain-id>
```

### Run Using REST

```
curl "localhost:1317/gov/proposals"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/votes"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/votes/{voter-id}"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/deposits"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/deposits/{depositer-addr}"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/tally"
```
```
curl "localhost:1317/gov/params"
```
```
curl "localhost:1317/gov/proposals/{proposal-id}/proposer"
```



