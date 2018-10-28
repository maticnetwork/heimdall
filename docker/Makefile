build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build-develop:
	docker build -t "maticnetwork/tendermint:develop" -f Dockerfile ..

push-develop:
	docker push "maticnetwork/tendermint:develop"

.PHONY: build build_develop push push_develop