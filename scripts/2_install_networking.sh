#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh

sudo yum install -y socat

# Check if socat is installed and get its path
SOCAT_PATH=$(command -v socat)
if [ -z "$SOCAT_PATH" ]; then
  echo "Error: socat is required but not installed. Please install socat and ensure it is in your PATH." >&2
  exit 1
fi

sudo ln -s "$SOCAT_PATH" /usr/local/bin/socat # this is required for the service to start

wget -O tacos https://github.com/palisadeinc/nitro-cloudsign/raw/refs/tags/${TACOS_VERSION}/tools/tacos
sudo cp tacos /usr/bin/tacos
sudo chmod 0500 /usr/bin/tacos

wget -O tacos.service https://raw.githubusercontent.com/palisadeinc/nitro-cloudsign/refs/tags/${TACOS_VERSION}/tools/service/tacos.service
sudo cp tacos.service /etc/systemd/system/tacos.service
sudo systemctl daemon-reload
sudo systemctl enable tacos
sudo systemctl start tacos