# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI](https://golangci.com/badges/github.com/maticnetwork/heimdall.svg)](https://golangci.com/r/github.com/maticnetwork/heimdall)


Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source 

Make sure you have go1.11+ already installed

### Install 
```bash 
$ make install
```
### Init-heimdall 
```bash 
$ heimdalld init
$ heimdalld init --chain=mainnet        Will init with genesis.json for mainnet
$ heimdalld init --chain=mumbai         Will init with genesis.json for mumbai
$ heimdalld init --chain=amoy           Will init with genesis.json for amoy
```
### Run-heimdall 
```bash 
$ heimdalld start
```
#### Usage
```
$ heimdalld start                       Will start for mainnet by default
$ heimdalld start --chain=mainnet       Will start for mainnet
$ heimdalld start --chain=mumbai        Will start for mumbai
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

Latest docs are [here](https://wiki.polygon.technology/docs/category/heimdall) 
