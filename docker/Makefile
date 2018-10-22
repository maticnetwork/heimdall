build:
	@sh -c "'$(CURDIR)/build.sh'"

push:
	@sh -c "'$(CURDIR)/push.sh'"

build_develop:
	docker build -t "maticnetwork/tendermint:develop" -f Dockerfile.develop ..

build_testing:
	docker build --tag maticnetwork/testing -f ./Dockerfile.testing ..

push_develop:
	docker push "maticnetwork/tendermint:develop"

.PHONY: build build_develop push push_develop