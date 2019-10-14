#!/usr/bin/env sh

# start processes
heimdalld start > ./logs/heimdalld.log &
heimdalld rest-server > ./logs/heimdalld-rest-server.log &

# tail logs
tail -f ./logs/heimdalld.log ./logs/heimdalld-rest-server.log 