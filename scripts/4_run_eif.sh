#!/usr/bin/env bash

sudo systemctl restart tacos.service
nitro-cli run-enclave \
  --eif-path ~/cloudsign-enclave-app.eif \
  --cpu-count 2 \
  --memory 1024 \
  --enclave-cid 5 \
  --attach-console