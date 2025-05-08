#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh
nitro-cli build-enclave --docker-uri 335650072995.dkr.ecr.eu-west-2.amazonaws.com/external/custody/cloudsign-nitro/app:${CLOUDSIGN_VERSION} --output-file ~/cloudsign-${CLOUDSIGN_VERSION}.eif