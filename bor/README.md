# Bor Module

## Table of Contents

* [Preliminary terminology](#preliminary-terminology)
* [Overview](#overview)
* [How does it work](#how-does-it-work)
	* [How to propose a span](#how-to-propose-a-span)
* [Query commands](#query-commands)


## Preliminary terminology

* A `side-transaction` is a normal heimdall transaction but the data with which the message is composed needs to be voted on by the validators since the data is obscure to the consensus protocol itself and it has no way of validating the data's correctness.
* A `sprint` comprises of 16 bor blocks (configured in [bor](https://github.com/maticnetwork/launch/blob/fe86ba6cd16e5c36067a5ae49c0bad62ce8b1c3f/mainnet-v1/sentry/validator/bor/genesis.json#L26C18-L28)).
* A `span` comprises of 400 sprints in bor (check heimdall's bor [params](https://heimdall-api.polygon.technology/bor/params) endpoint ).

## Overview

The validators on bor chain produce blocks in sprints and spans. Hence, it is imperative for the protocol to formalise the validators who will be producers in a range of blocks (`span`). The `bor` module in heimdall facilitates this by pseudo-randomly selecting validators who will producing blocks (producers) from the current validator set. The bor chain fetches and persists this information before the next span begins. `bor` module is a crucial component in heimdall since the PoS chain "liveness" depends on it.

## How does it work

A `Span` is defined by the data structure:

```
type Span struct {
	ID                uint64       `json:"span_id" yaml:"span_id"`
	StartBlock        uint64       `json:"start_block" yaml:"start_block"`
	EndBlock          uint64       `json:"end_block" yaml:"end_block"`
	ValidatorSet      ValidatorSet `json:"validator_set" yaml:"validator_set"`
	SelectedProducers []Validator  `json:"selected_producers" yaml:"selected_producers"`
	ChainID           string       `json:"bor_chain_id" yaml:"bor_chain_id"`
}
```
where ,

* `ID` means the id of the span, calculated by monotonically incrementing the ID of the previous span.
* `StartBlock` corresponds to the block in bor from which the given span would commence.
* `EndBlock` corresponds to the block in bor at which the given span would conclude.
* `ValidatorSet` defines the set of active validators.
* `SelectedProducers` are the validators selected to produce blocks in bor from the validator set.
* `ChainID` corresponds to bor chain ID.

A validator on heimdall can construct a span proposal message:

```
type MsgProposeSpan struct {
    ID         uint64                  `json:"span_id"`
    Proposer   hmTypes.HeimdallAddress `json:"proposer"`
    StartBlock uint64                  `json:"start_block"`
    EndBlock   uint64                  `json:"end_block"`
    ChainID    string                  `json:"bor_chain_id"`
    Seed       common.Hash             `json:"seed"`
}
```

Upon broadcasting the message, it is initially checked by `HandleMsgProposeSpan` handler for basic sanity (verify whether the proposed span is in continuity, appropriate span duration, correct chain ID, etc.). Since this is a side-transaction, the validators then vote on the data present in `MsgProposeSpan` on the basis of its correctness. All these checks are done in `SideHandleMsgSpan` (verifying `Seed`, span continuity, etc) and if correct, the validator would vote `YES`.
Finally, if there are 2/3+ `YES` votes, the `PostHandleMsgEventSpan` persists the proposed span in the state via the keeper :  

```
err := k.FreezeSet(ctx, msg.ID, msg.StartBlock, msg.EndBlock, msg.ChainID, msg.Seed)
if err != nil {
	k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
	return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
}
```

`FreezeSet` internally invokes `SelectNextProducers`, which pseudo-randomly picks producers from the validator set, leaning more towards validators with higher voting power based on stake:

```
// select next producers
newProducers, err := k.SelectNextProducers(ctx, seed)
if err != nil {
	return err
}
```

and then initialises and stores the span:

```
newSpan := hmTypes.NewSpan(
	id,
	startBlock,
	endBlock,
	k.sk.GetValidatorSet(ctx),
	newProducers,
	borChainID,
)

return k.AddNewSpan(ctx, newSpan)
```

### How to propose a span

A validator can leverage the CLI to propose a span like so :

```
heimdallcli tx bor propose-span --proposer <VALIDATOR ADDRESS> --start-block <BOR_START_BLOCK> --span-id <SPAN_ID> --bor-chain-id <BOR_CHAIN_ID>
```

Or the REST server : 

```
curl -X POST "localhost:1317/bor/propose-span?bor-chain-id=<BOR_CHAIN_ID>&start-block=<BOR_START_BLOCK>&span-id=<SPAN_ID>"
```

## Query commands

One can run the following query commands from the bor module :

* `span` - Query the span corresponding to the given span id.
* `latest span` - Query the latest span.
* `params` - Fetch the parameters associated to bor module.
* `spanlist` - Fetch span list.
* `next-span-seed` - Query the seed for the next span.
* `propose-span` - Print the `propose-span` command.

### CLI commands

```
heimdallcli query bor span --span-id=<SPAN_ID>
```

```
heimdallcli query bor latest-span
```

```
heimdallcli query bor params
```

```
heimdallcli query bor spanlist --page=<PAGE_NUM> --limit=<LIMIT>
```

```
heimdallcli query bor next-span-seed
```

```
heimdallcli query bor propose-span --proposer <VALIDATOR ADDRESS> --start-block <BOR_START_BLOCK> --span-id <SPAN_ID> --bor-chain-id <BOR_CHAIN_ID>
```

### REST endpoints

```
curl localhost:1317/bor/span/<SPAN_ID>
```

```
curl localhost:1317/bor/latest-span
```

```
curl localhost:1317/bor/params
```

```
curl localhost:1317/bor/next-span-seed
```

```
curl "localhost:1317/bor/prepare-next-span?span_id=<SPAN_ID>&start_block=<BOR_START_BLOCK>&chain_id="<BOR_CHAIN_ID>""
```