docker run --name hm100 -it -d   \
		-v ~/mytestnet/node0:/root/.heimdalld \
		-v ~/mytestnet/logs0:/go/src/github.com/maticnetwork/heimdall/logs \
		-p 1307:1317 \
		"maticnetwork/tendermint:develop"
docker run --name hm111 -it -d -P \
      -v ~/mytestnet/node1:/root/.heimdalld \
      -v ~/mytestnet/logs1:/go/src/github.com/maticnetwork/heimdall/logs \
      -p 1317:1317 \
      "maticnetwork/tendermint:develop"

docker run --name hm122 -it -d \
		-v ~/mytestnet/node2:/root/.heimdalld \
		-v ~/mytestnet/logs2:/go/src/github.com/maticnetwork/heimdall/logs \
		-p 1327:1317 \
		"maticnetwork/tendermint:develop"

  docker run --name hm133 -it -d  \
      -v ~/mytestnet/node3:/root/.heimdalld \
      -v ~/mytestnet/logs3:/go/src/github.com/maticnetwork/heimdall/logs \
      -p 1337:1317 \
      "maticnetwork/tendermint:develop"
