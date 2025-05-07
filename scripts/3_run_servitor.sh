#!/usr/bin/env bash

set -e

source ./var_version.sh

read -p "Pairing key: " PAIRING_KEY
read -p "DB Data Source (eg postgres://user:password@host:port/db): " DB_DATA_SOURCE
read -p "TSM DB Data Source (eg postgres://user:password@host:port/db): " TSM_DB_DATA_SOURCE

nohup servitor > /var/log/servitor.log &