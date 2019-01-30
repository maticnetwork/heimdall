#!/usr/bin/env sh

# start processes
./build/heimdalld start > ./logs/heimdalld.log &
./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &
./build/bridge start > ./logs/bridge.log &

# tail logs
# tail ./logs/heimdalld.log ./logs/heimdalld-rest-server.log ./logs/bridge.log
