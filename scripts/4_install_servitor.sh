#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh

wget https://github.com/palisadeinc/nitro-cloudsign/releases/download/${SERVITOR_VERSION}/servitor-linux-amd64-${SERVITOR_VERSION}
chmod +x servitor-linux-amd64-${SERVITOR_VERSION}
sudo mv servitor-linux-amd64-${SERVITOR_VERSION} /usr/local/bin/servitor
