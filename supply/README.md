# Supply Module

## Table of Contents

* [Overview](#overview)
* [Query commands](#query-commands)

## Overview

The supply functionality passively tracks the total supply of coins within a chain,
provides a pattern for modules to hold/interact with coins, and introduces the invariant check to verify a chain's total supply. The total supply of the network is equal to the sum of all coins from the account.

## Query commands

One can run the following query commands from the bank module :

* `total` - Query for a total supply of the chain with different coins. Takes params for checking the total supply of a particular coin

### CLI commands

```
heimdallcli query supply total [denom]
```
