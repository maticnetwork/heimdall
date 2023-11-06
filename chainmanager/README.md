# Chainmanager Module

## Table of Contents

* [Overview](#overview)
* [Query commands](#query-commands)

## Overview

The chainmanager module is responsible for fetching the chainmanager params. These params include contract address of mainchain (Ethereum) and maticchain (Bor), chain ids, mainchain and maticchain confirmation blocks

## Query commands

One can run the following query commands from the chainmanager module :

* `params` - Fetch the parameters associated to chainmanager module.

### CLI commands

```
heimdallcli query chainmanager params
```

### REST endpoints

```
curl localhost:1317/chainmanager/params
```
