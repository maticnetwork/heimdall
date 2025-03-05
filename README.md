# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI Lint](https://github.com/maticnetwork/heimdall/actions/workflows/ci.yml/badge.svg)](https://github.com/maticnetwork/heimdall/actions)

Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source

Make sure you have Go v1.20+ already installed.

### Install

```bash
$ make install
```

### Init Heimdall

```bash
$ heimdalld init
$ heimdalld init --chain=mainnet        Will init with genesis.json for mainnet
$ heimdalld init --chain=amoy           Will init with genesis.json for amoy
```

### Run Heimdall

```bash
$ heimdalld start
```

#### Usage

```bash
$ heimdalld start                       Will start for mainnet by default
$ heimdalld start --chain=mainnet       Will start for mainnet
$ heimdalld start --chain=amoy          Will start for amoy
$ heimdalld start --chain=local         Will start for local with NewSelectionAlgoHeight = 0
```

### Run rest server

```bash
$ heimdalld rest-server
```

### Run bridge

```bash
$ heimdalld bridge
```

### Develop using Docker

You can build and run Heimdall using the included Dockerfile in the root directory:

```bash
docker build -t heimdall .
docker run heimdall
```

### Documentation

Latest docs are [here](https://docs.polygon.technology/pos/).
