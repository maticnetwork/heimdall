dep:
	dep ensure -v
	mkdir -p vendor/github.com/tendermint vendor/github.com/ethereum
	git clone -b v0.9.2 --single-branch --depth 1 https://github.com/tendermint/iavl vendor/github.com/tendermint/iavl
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

run-heimdall:
	./build/heimdalld start

rest-server:
	./build/heimdallcli rest-server

start:
	mkdir -p logs
	./build/heimdalld start > ./logs/heimdalld.log &
	./build/heimdallcli rest-server > ./logs/heimdallcli.log &

#
# docker commands
#

build-docker:
	cd docker; make build

build-docker-develop:
	cd docker; make build-develop

.PHONY: dep build
