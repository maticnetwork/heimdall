#!/bin/sh
set -e

PKG="heimdall"

if ! getent passwd $PKG >/dev/null ; then
    adduser --disabled-password --disabled-login --shell /usr/sbin/nologin --quiet --system --no-create-home --home /nonexistent $PKG
    echo "Created system user $PKG"
fi

