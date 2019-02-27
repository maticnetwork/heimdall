# Fetch git latest tag
LATEST_GIT_TAG:=$(shell git describe --tags $(git rev-list --tags --max-count=1))

dep:
	dep ensure -v
	mkdir -p vendor/github.com/tendermint vendor/github.com/ethereum
	git clone -b v0.11.0 --single-branch --depth 1 https://github.com/tendermint/iavl vendor/github.com/tendermint/iavl
	git clone -b v1.8.19 --single-branch --depth 1 https://github.com/ethereum/go-ethereum vendor/github.com/ethereum/go-ethereum

clean:
	rm -rf build

tests:
	go test  -v ./...

build: clean
	mkdir -p build
	go build -o build/heimdalld cmd/heimdalld/main.go
	go build -o build/heimdallcli cmd/heimdallcli/main.go
	go build -o build/bridge bridge/bridge.go

build-testnet:
	go build -o build/testnet cmd/testnet/main.go

contracts:
	abigen --abi=contracts/rootchain/rootchain.abi --pkg=rootchain --out=contracts/rootchain/rootchain.go
	abigen --abi=contracts/stakemanager/stakemanager.abi --pkg=stakemanager --out=contracts/stakemanager/stakemanager.go

init-heimdall:
	./build/heimdalld init

show-account-heimdall:
	./build/heimdalld show-account

show-node-id:
	./build/heimdalld tendermint show-node-id

run-heimdall:
	./build/heimdalld start 

start-heimdall:
	mkdir -p ./logs
	./build/heimdalld start > ./logs/heimdalld.log &

reset-heimdalld:
	./build/heimdalld unsafe-reset-all 

run-server:
	./build/heimdalld rest-server 

start-server:
	mkdir -p ./logs
	./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &

start:
	mkdir -p ./logs
	./build/heimdalld start > ./logs/heimdalld.log &
	./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &
	tail -f ./logs/heimdalld.log

run-bridge:
	./build/bridge start

start-bridge:
	mkdir -p logs
	./build/bridge start > ./logs/bridge.log

start-all: 
	mkdir -p ./logs
	bash docker/start-heimdall.sh

#
# docker commands
#

build-docker:
	@echo Fetching latest tag: $(LATEST_GIT_TAG)
	git checkout $(LATEST_GIT_TAG)
	docker build -t "maticnetwork/heimdall:$(LATEST_GIT_TAG)" -f docker/Dockerfile .

push-docker:
	@echo Pushing docker tag image: $(LATEST_GIT_TAG)
	docker push "maticnetwork/heimdall:$(LATEST_GIT_TAG)"

build-docker-develop:
	docker build -t "maticnetwork/heimdall:develop" -f docker/Dockerfile.develop .

.PHONY: contracts dep build
