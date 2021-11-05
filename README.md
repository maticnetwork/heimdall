# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI](https://golangci.com/badges/github.com/maticnetwork/heimdall.svg)](https://golangci.com/r/github.com/maticnetwork/heimdall)


Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source 

Make sure your have go1.11+ already installed

### Install 
```bash 
$ make install network=mainnet
```
#### Usage 
```
# make install					Will generate for mainnet by default
# make install network=mainnet			Will generate for mainnet
# make install network=mumbai			Will generate for mumbai
# make install network=local			Will generate for local with NewSelectionAlgoHeight = 0
# make install network=anythingElse		Will generate for mainnet by default
```
### Run-heimdall 
```bash 
$ heimdalld start
```

### Run rest server

```bash 
$ heimdalld rest-server 
```

### Develop using Docker

You can build and run Heimdall using the included Dockerfile in the root directory:

```bash
docker build -t heimdall .
docker run heimdall
```

### Documentation 

Latest docs are [here](https://docs.matic.network/) 
