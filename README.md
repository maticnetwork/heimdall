# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI](https://golangci.com/badges/github.com/maticnetwork/heimdall.svg)](https://golangci.com/r/github.com/maticnetwork/heimdall)


Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source 

Make sure your have go1.11+ already installed

### Install

```bash 
$ make process-template
```
```
make process-template							Will generate for mainnet by default
make process-template network=mainnet			Will generate for mainnet
make process-template network=mumbai			Will generate for mumbai
make process-template network=local             Will generate for local with NewSelectionAlgoHeight = 0
make process-template network=anythingElse      Will generate for mainnet by default
```

```bash 
$ make install 
```  

### Run-heimdall 
```bash 
$ heimdalld start
```

### Run rest server

```bash 
$ heimdalld rest-server 
```


### Documentation 

Latest docs are [here](https://docs.matic.network/) 
