#!/usr/bin/env bash

set -e

source ./var_version.sh
nitro-cli build-enclave --docker-uri 335650072995.dkr.ecr.eu-west-2.amazonaws.com/external/custody/cloudsign-nitro/app:${VERSION} --output-file ~/cloudsign.eif