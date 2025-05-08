#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh
aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 335650072995.dkr.ecr.eu-west-2.amazonaws.com
docker pull 335650072995.dkr.ecr.eu-west-2.amazonaws.com/external/custody/cloudsign-nitro/app:${CLOUDSIGN_VERSION}