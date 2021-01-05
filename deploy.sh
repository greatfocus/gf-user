#!/bin/sh
export PATH=$PATH:/usr/local/go/bin

# Building our app
GOOS=linux GOARCH=amd64 go build

# project config
sudo chmod -R 700 dev.json
sudo chown -R muthurimi dev.json