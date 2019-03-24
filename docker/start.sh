#!/usr/bin/env sh

# start processes
./build/heimdalld start > ./logs/heimdalld.log &
./build/heimdalld rest-server > ./logs/heimdalld-rest-server.log &

# tail logs
tail -f ./logs/heimdalld.log ./logs/heimdalld-rest-server.log 