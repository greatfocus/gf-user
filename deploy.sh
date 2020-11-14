#!/bin/sh
export PATH=$PATH:/usr/local/go/bin

# Building our app
GOOS=linux GOARCH=amd64 go build

# project config
sudo chmod -R 700 dev.json
sudo chown -R muthurimi dev.json

# create service
sudo systemctl stop gf-user 
sudo systemctl disable gf-user  
sudo cp gf-user.service /etc/systemd/system/gf-user.service
systemctl daemon-reload

# start user 
sudo systemctl enable gf-user 
sudo systemctl start gf-user  