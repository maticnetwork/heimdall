dep:
	dep ensure -v
	mkdir -p vendor/github.com/tendermint
	git clone -b v0.9.2 --single-branch --depth 1 https://github.com/tendermint/iavl vendor/github.com/tendermint/iavl

build:
	cd cmd/heimdalld && go build main.go && cd -
	cd cmd/heimdallcli && go build main.go && cd -

.PHONY: dep
