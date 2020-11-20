# Fetch git latest tag
LATEST_GIT_TAG:=$(shell git describe --tags $(git rev-list --tags --max-count=1))
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/maticnetwork/heimdall/version.Name=heimdall \
		  -X github.com/maticnetwork/heimdall/version.ServerName=heimdalld \
		  -X github.com/maticnetwork/heimdall/version.ClientName=heimdallcli \
		  -X github.com/maticnetwork/heimdall/version.Version=$(VERSION) \
		  -X github.com/maticnetwork/heimdall/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.Name=heimdall \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=heimdalld \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=heimdallcli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

clean:
	rm -rf build

tests:
	# go test  -v ./...

	go test -v ./app/ ./auth/ ./clerk/ ./sidechannel/ ./bank/ ./chainmanager/ ./topup/ ./checkpoint/ ./staking/ ./bor/ ./gov/ -cover -coverprofile=cover.out


build: clean
	mkdir -p build
	go build -o build/heimdalld ./cmd/heimdalld
	go build -o build/heimdallcli ./cmd/heimdallcli
	go build -o build/bridge bridge/bridge.go
	@echo "====================================================\n==================Build Successful==================\n===================================================="

install:
	go install $(BUILD_FLAGS) ./cmd/heimdalld
	go install $(BUILD_FLAGS) ./cmd/heimdallcli
	go install $(BUILD_FLAGS) bridge/bridge.go

contracts:
	abigen --abi=contracts/rootchain/rootchain.abi --pkg=rootchain --out=contracts/rootchain/rootchain.go
	abigen --abi=contracts/stakemanager/stakemanager.abi --pkg=stakemanager --out=contracts/stakemanager/stakemanager.go
	abigen --abi=contracts/slashmanager/slashmanager.abi --pkg=slashmanager --out=contracts/slashmanager/slashmanager.go
	abigen --abi=contracts/statereceiver/statereceiver.abi --pkg=statereceiver --out=contracts/statereceiver/statereceiver.go
	abigen --abi=contracts/statesender/statesender.abi --pkg=statesender --out=contracts/statesender/statesender.go
	abigen --abi=contracts/stakinginfo/stakinginfo.abi --pkg=stakinginfo --out=contracts/stakinginfo/stakinginfo.go
	abigen --abi=contracts/validatorset/validatorset.abi --pkg=validatorset --out=contracts/validatorset/validatorset.go
	abigen --abi=contracts/erc20/erc20.abi --pkg=erc20 --out=contracts/erc20/erc20.go


init-heimdall:
	./build/heimdalld init

show-account-heimdall:
	./build/heimdalld show-account

show-node-id:
	./build/heimdalld tendermint show-node-id

run-heimdall:
	./build/heimdalld start

start-heimdall:
	mkdir -p ./logs &
	./build/heimdalld start > ./logs/heimdalld.log &

reset-heimdall:
	./build/heimdalld unsafe-reset-all
	./build/bridge purge-queue
	rm -rf ~/.heimdalld/bridge

run-server:
	./build/heimdalld rest-server

start-server:
	mkdir -p ./logs &
	./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &

start:
	mkdir -p ./logs
	bash docker/start.sh

run-bridge:
	./build/bridge start --all

start-bridge:
	mkdir -p logs &
	./build/bridge start --all > ./logs/bridge.log &

start-all:
	mkdir -p ./logs
	bash docker/start-heimdall.sh

#
# Code quality
#

LINT_COMMAND := $(shell command -v golangci-lint 2> /dev/null)
lint:
ifndef LINT_COMMAND
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.23.8
endif
	golangci-lint run

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

.PHONY: contracts build
