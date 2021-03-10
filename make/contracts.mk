.PHONY: contracts
contracts:
	abigen --abi=contracts/rootchain/rootchain.abi --pkg=rootchain --out=contracts/rootchain/rootchain.go
	abigen --abi=contracts/stakemanager/stakemanager.abi --pkg=stakemanager --out=contracts/stakemanager/stakemanager.go
	abigen --abi=contracts/slashmanager/slashmanager.abi --pkg=slashmanager --out=contracts/slashmanager/slashmanager.go
	abigen --abi=contracts/statereceiver/statereceiver.abi --pkg=statereceiver --out=contracts/statereceiver/statereceiver.go
	abigen --abi=contracts/statesender/statesender.abi --pkg=statesender --out=contracts/statesender/statesender.go
	abigen --abi=contracts/stakinginfo/stakinginfo.abi --pkg=stakinginfo --out=contracts/stakinginfo/stakinginfo.go
	abigen --abi=contracts/validatorset/validatorset.abi --pkg=validatorset --out=contracts/validatorset/validatorset.go
	abigen --abi=contracts/erc20/erc20.abi --pkg=erc20 --out=contracts/erc20/erc20.go
