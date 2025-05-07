#!/usr/bin/env bash

set -e

sudo cp tacos /usr/bin/tacos
sudo chmod 0500 /usr/bin/tacos
sudo cp ./service/tacos.service /etc/systemd/system/tacos.service
sudo systemctl daemon-reload
sudo systemctl enable tacos
sudo systemctl start tacos
