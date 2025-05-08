#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh

sudo systemctl restart tacos.service
nitro-cli run-enclave \
  --eif-path ~/cloudsign-${CLOUDSIGN_VERSION}.eif \
  --cpu-count 2 \
  --memory 1024 \
  --enclave-cid 5 \
  --attach-console