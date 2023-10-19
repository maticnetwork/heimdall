# Supply Module

## Table of Contents

* [Overview](#overview)
* [How does it work](#how-does-it-work)
* [How to add an event](#how-to-add-an-event)
* [Query commands](#query-commands)

## Overview

The supply functionality passively tracks the total supply of coins within a chain,
provides a pattern for modules to hold/interact with Coins, and introduces the invariant check to verify a chain's total supply. The total Supply of the network is equal to the sum of all coins from the account. The total supply is updated every time a Coin is minted (eg: as part of the inflation mechanism) or burned (eg: due to slashing or if a governance proposal is vetoed).

## Query commands

One can run the following query commands from the clerk module :

* `total` - Query for a total supply of the chain. Takes params for checking the total supply of a particular coin

### CLI commands

```
heimdallcli query supply total [denom]
```
