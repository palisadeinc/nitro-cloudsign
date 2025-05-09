#!/usr/bin/env bash

set -e
set -u
set -o pipefail

source ./variables.sh

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

# Check if servitor_config.env exists and load it if present
CONFIG_FILE="servitor_config.env"
if [ -f "${CONFIG_FILE}" ]; then
    echo "Found existing configuration file: ${CONFIG_FILE}"
    source "${CONFIG_FILE}"
    echo "Loaded configuration from ${CONFIG_FILE}"
    echo "You can press Enter to use these values or enter new ones:"
    echo "  PAIRING_KEY: ${PAIRING_KEY:-<not set>}"
    echo "  DB_DATA_SOURCE: ${DB_DATA_SOURCE:-<not set>}"
    echo "  TSM_DB_DATA_SOURCE: ${TSM_DB_DATA_SOURCE:-<not set>}"
    echo "  SIGN_MODE_ACTIVE: ${SIGN_MODE_ACTIVE:-<not set>}"
    echo "  DB_ENCRYPTION_KEY_REF: ${DB_ENCRYPTION_KEY_REF:-<not set>}"
else
    echo "No existing configuration found. Please enter the required information:"
    PAIRING_KEY=""
    DB_DATA_SOURCE=""
    TSM_DB_DATA_SOURCE=""
    SIGN_MODE_ACTIVE="true"
    DB_ENCRYPTION_KEY_REF=""
fi

TEMP_INPUT=""

read -p "Pairing key [${PAIRING_KEY:-<not set>}]: " TEMP_INPUT
if [ -n "${TEMP_INPUT}" ]; then
    PAIRING_KEY="${TEMP_INPUT}"
fi

TEMP_INPUT=""
read -p "DB Data Source (eg postgres://user:password@host:port/db?sslmode=require) [${DB_DATA_SOURCE:-<not set>}]: " TEMP_INPUT
if [ -n "${TEMP_INPUT}" ]; then
    DB_DATA_SOURCE="${TEMP_INPUT}"
fi

TEMP_INPUT=""
read -p "TSM DB Data Source (eg postgres://user:password@host:port/db?sslmode=require) [${TSM_DB_DATA_SOURCE:-<not set>}]: " TEMP_INPUT
if [ -n "${TEMP_INPUT}" ]; then
    TSM_DB_DATA_SOURCE="${TEMP_INPUT}"
fi

TEMP_INPUT=""
read -p "Sign Mode Active (true/false) [${SIGN_MODE_ACTIVE:-<not set>}]: " TEMP_INPUT
if [ -n "${TEMP_INPUT}" ]; then
    SIGN_MODE_ACTIVE="${TEMP_INPUT}"
fi

TEMP_INPUT=""
read -p "DB Encryption Key Ref [${DB_ENCRYPTION_KEY_REF:-<not set>}]: " TEMP_INPUT
if [ -n "${TEMP_INPUT}" ]; then
    DB_ENCRYPTION_KEY_REF="${TEMP_INPUT}"
fi

# Write configuration to a file for future reference
CONFIG_FILE="servitor_config.env"
echo "Writing configuration to ${CONFIG_FILE}"
cat > ${CONFIG_FILE} << EOF
PAIRING_KEY=${PAIRING_KEY}
DB_DATA_SOURCE=${DB_DATA_SOURCE}
TSM_DB_DATA_SOURCE=${TSM_DB_DATA_SOURCE}
SIGN_MODE_ACTIVE=${SIGN_MODE_ACTIVE}
DB_ENCRYPTION_KEY_REF=${DB_ENCRYPTION_KEY_REF}
EOF
chmod 600 ${CONFIG_FILE}
echo "Configuration saved to ${CONFIG_FILE}"



echo "Starting servitor... Log file: ${SERVITOR_LOG_FILE}"
PAIRING_KEY="${PAIRING_KEY}" DB_DATA_SOURCE="${DB_DATA_SOURCE}" TSM_DB_DATA_SOURCE="${TSM_DB_DATA_SOURCE}" SIGN_MODE_ACTIVE="${SIGN_MODE_ACTIVE}" DB_ENCRYPTION_KEY_REF="${DB_ENCRYPTION_KEY_REF}" servitor > ${SERVITOR_LOG_FILE}

echo "Servitor stopped."