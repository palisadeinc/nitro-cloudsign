#!/usr/bin/env bash

set -e

# Check if socat is installed and get its path
SOCAT_PATH=$(command -v socat)
if [ -z "$SOCAT_PATH" ]; then
  echo "Error: socat is required but not installed. Please install socat and ensure it is in your PATH." >&2
  exit 1
fi

ln -s "$SOCAT_PATH" /usr/local/bin/socat # this is required for the service to start

sudo cp tacos /usr/bin/tacos
sudo chmod 0500 /usr/bin/tacos
sudo cp ./service/tacos.service /etc/systemd/system/tacos.service
sudo systemctl daemon-reload
sudo systemctl enable tacos
sudo systemctl start tacos
