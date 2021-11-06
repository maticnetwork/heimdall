# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI](https://golangci.com/badges/github.com/maticnetwork/heimdall.svg)](https://golangci.com/r/github.com/maticnetwork/heimdall)


Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source 

Make sure your have go1.11+ already installed

### Install 
```bash 
$ make install
```
### Init-heimdall 
```bash 
$ heimdalld init
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
$ heimdalld start --chain=local         Will start for local with NewSelectionAlgoHeight = 0
```

### Run rest server
```bash 
$ heimdalld rest-server 
```


### Documentation 

Latest docs are [here](https://docs.matic.network/) 
