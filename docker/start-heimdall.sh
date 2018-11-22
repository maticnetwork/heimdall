#!/usr/bin/env sh

# start redis
redis-server /go/src/github.com/maticnetwork/heimdall/docker/redis.conf

# start processes
./build/heimdalld start > ./logs/heimdalld.log &
./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &
./build/bridge start > ./logs/bridge.log &

# tail logs
tail -f ./logs/heimdalld.log ./logs/heimdalld-rest-server.log ./logs/bridge.log
