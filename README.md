# Heimdall

[![Go Report Card](https://goreportcard.com/badge/github.com/maticnetwork/heimdall)](https://goreportcard.com/report/github.com/maticnetwork/heimdall) [![CircleCI](https://circleci.com/gh/maticnetwork/heimdall/tree/master.svg?style=shield)](https://circleci.com/gh/maticnetwork/heimdall/tree/master) [![GolangCI](https://golangci.com/badges/github.com/maticnetwork/heimdall.svg)](https://golangci.com/r/github.com/maticnetwork/heimdall)


Validator node for Matic Network. It uses peppermint, customized [Tendermint](https://github.com/tendermint/tendermint).

### Install from source

Make sure your have go1.15+ already installed

### Install
```bash
$ make install
```

### Run-heimdall
```bash
$ heimdalld init --chain-id <chain-id> <moniker>
$ heimdalld start
```

### Run rest server
REST server runs in-process with heimdall node.
To enable REST server edit `app.toml` and enable `api.enable`.
For serving `swagger`, enable `api.swagger`.

### Documentation

Latest docs are [here](https://docs.matic.network/)
Test webhook trigger