#!/usr/bin/env sh

if [ ! -f $HEIMDALL_DIR/config/heimdall-config.toml ]; then
    heimdalld --home=$HEIMDALL_DIR init
fi;

# start processes
heimdalld --home=$HEIMDALL_DIR start &
heimdalld --home=$HEIMDALL_DIR rest-server &
sleep 100
bridge --home=$HEIMDALL_DIR start --all

exit $?
