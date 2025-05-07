#!/usr/bin/env bash

set -e

DOCKER_IMAGE="335650072995.dkr.ecr.eu-west-2.amazonaws.com/external/custody/cloudsign/app"
DOCKER_TAG="1.8.0"
ENCLAVE_OUTPUT_FILE="cloudsign-enclave-app.eif"
docker pull ${DOCKER_IMAGE}:${DOCKER_TAG}
nitro-cli build-enclave --docker-uri ${DOCKER_IMAGE}:${DOCKER_TAG} --output-file ~/cloudsign-enclave-app.eif