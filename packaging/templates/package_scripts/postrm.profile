#!/bin/bash
#
###################
# Remove heimdall profile installation
###################
sudo rm /var/lib/heimdall/config/heimdall-config.toml
sudo rm /var/lib/heimdall/config/config.toml
sudo systemctl daemon-reload