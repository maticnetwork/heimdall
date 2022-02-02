#!/usr/bin/env sh

if [ ! -f $HEIMDALL_DIR/config/heimdall-config.toml ]; then
    heimdalld --home=$HEIMDALL_DIR init
fi;

if [ "$1" = 'bridge' ]; then
    shift
    exec bridge --home=$HEIMDALL_DIR "$@"
fi

if [ "$1" = 'heimdallcli' ]; then
    shift
    exec heimdallcli --home=$HEIMDALL_DIR "$@"
fi

exec heimdalld --home=$HEIMDALL_DIR "$@"
