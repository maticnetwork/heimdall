#!/usr/bin/env sh

if [ "$1" = 'heimdallcli' ]; then
    shift
    exec heimdallcli --home=$HEIMDALL_DIR "$@"
fi

exec heimdalld --home=$HEIMDALL_DIR "$@"
