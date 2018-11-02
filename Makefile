dep:
	dep ensure -v
	mkdir -p vendor/github.com/tendermint vendor/github.com/ethereum
	git clone -b v0.11.0 --single-branch --depth 1 https://github.com/tendermint/iavl vendor/github.com/tendermint/iavl
	git clone -b v1.8.17 --single-branch --depth 1 https://github.com/ethereum/go-ethereum vendor/github.com/ethereum/go-ethereum

clean:
	rm -rf build

build: clean
	mkdir -p build
	go build -o build/heimdalld cmd/heimdalld/main.go
	go build -o build/heimdallcli cmd/heimdallcli/main.go

contracts:
	# mkdir -p contracts/validatorset
	# abigen --abi=contracts/validatorset/validatorset.abi --pkg=validatorset --out=contracts/validatorset/validatorset.go

init-heimdall:
	./build/heimdalld init

run-heimdall:
	./build/heimdalld start

reset-heimdalld:
	./build/heimdalld unsafe-reset-all

rest-server:
	./build/heimdallcli rest-server

start:
	./build/heimdalld start > ./logs/heimdalld.log &
	./build/heimdallcli rest-server > ./logs/heimdallcli.log &
	tail -f ./logs/heimdalld.log ./logs/heimdallcli.log

#
# docker commands
#

build-docker:
	cd docker; make build

build-docker-develop:
	cd docker; make build-develop

run-docker-develop:
	docker run --name node0 -it \
		-v ~/.heimdalld:/root/.heimdalld \
		-v `pwd`/logs:/go/src/github.com/maticnetwork/heimdall/logs \
		-p 1317:1317 \
		"maticnetwork/heimdall:develop"

.PHONY: dep build
