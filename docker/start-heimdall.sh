#!/usr/bin/env sh

# start processes
./build/heimdalld start > ./logs/heimdalld.log &
./build/heimdallcli rest-server > ./logs/heimdallcli.log &
./build/bridge start > ./logs/bridge.log &

# tail logs
tail -f ./logs/heimdalld.log ./logs/heimdallcli.log ./logs/bridge.log
