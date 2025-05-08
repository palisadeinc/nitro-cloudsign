#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh

read -p "Pairing key: " PAIRING_KEY
read -p "DB Data Source (eg postgres://user:password@host:port/db): " DB_DATA_SOURCE
read -p "TSM DB Data Source (eg postgres://user:password@host:port/db): " TSM_DB_DATA_SOURCE

# Ensure the log file directory exists and is writable
LOG_DIR=$(dirname "${SERVITOR_LOG_FILE}")
if [ ! -d "${LOG_DIR}" ]; then
    mkdir -p "${LOG_DIR}"
    if [ $? -ne 0 ]; then
        echo "Error: Could not create log directory ${LOG_DIR}"
        exit 1
    fi
fi

if [ ! -w "${LOG_DIR}" ]; then
    echo "Error: Log directory ${LOG_DIR} is not writable."
    exit 1
fi

echo "Starting servitor... Log file: ${SERVITOR_LOG_FILE}"
PAIRING_KEY=${PAIRING_KEY} DB_DATA_SOURCE=${DB_DATA_SOURCE} TSM_DB_DATA_SOURCE=${TSM_DB_DATA_SOURCE} servitor > ${SERVITOR_LOG_FILE}

echo "Servitor stopped."