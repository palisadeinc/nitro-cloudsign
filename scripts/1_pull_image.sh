#!/usr/bin/env bash

set -e

source ./var_version.sh
aws ecr get-login-password --region eu-west-2 | docker login --username AWS --password-stdin 335650072995.dkr.ecr.eu-west-2.amazonaws.com
docker pull 335650072995.dkr.ecr.eu-west-2.amazonaws.com/external/custody/cloudsign-nitro/app:${VERSION}